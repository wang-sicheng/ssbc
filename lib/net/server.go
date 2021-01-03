package net

import (
	"fmt"
	"github.com/cloudflare/cfssl/log"
	gmux "github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/ssbc/lib/mysql"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"sync"
)

type Server struct {
	// The home directory for the server
	HomeDir string
	// BlockingStart if true makes the Start function blocking;
	// It is non-blocking by default.
	BlockingStart bool
	// The server's configuration
	Config *ServerConfig
	// The server mux
	mux *gmux.Router
	// The current listener for this server
	listener net.Listener
	// An error which occurs when serving
	serveError error

	wait chan bool
	// Server mutex
	mutex sync.Mutex
}

func (s *Server) Init(renew bool) (err error) {
	err = s.init(renew)
	err2 := s.closeDB()
	if err2 != nil {
		log.Errorf("Close DB failed: %s", err2)
	}
	return err
}

func (s *Server) init(renew bool) (err error) {
	//serverVersion := "SSBC v1.0"
	//log.Infof("Server Version: %s", serverVersion)

	// Initialize the config
	err = s.initConfig()
	if err != nil {
		return err
	}
	// Initialize the default CA last

	// Successful initialization
	return nil
}

func (s *Server) closeDB() error {

	return mysql.CloseDB()
}
func (s *Server) initConfig() (err error) {

	if s.HomeDir == "" {

		s.HomeDir, err = os.Getwd()

		if err != nil {
			return errors.Wrap(err, "Failed to get server's home directory")
		}
	}
	// Make home directory absolute, if not already
	absoluteHomeDir, err := filepath.Abs(s.HomeDir)
	if err != nil {
		return fmt.Errorf("Failed to make server's home directory path absolute: %s", err)
	}
	s.HomeDir = absoluteHomeDir
	// Create config if not set
	if s.Config == nil {
		s.Config = new(ServerConfig)
	}
	cfg := s.Config
	// Set log level if debug is true
	if cfg.Debug {
		log.Level = log.LevelDebug
	}
	return nil
}

func (s *Server) Start() (err error) {
	log.Infof("Starting server in home directory: %s", s.HomeDir)

	s.serveError = nil

	if s.listener != nil {
		return errors.New("server is already started")
	}

	// Initialize the server
	err = s.init(false) // 设置主目录和配置文件
	if err != nil {
		err2 := s.closeDB()
		if err2 != nil {
			log.Errorf("Close DB failed: %s", err2)
		}
		return err
	}

	// Register http handlers
	s.registerHandlers()

	// Start listening and serving
	err = s.listenAndServe()
	if err != nil {
		err2 := s.closeDB()
		if err2 != nil {
			log.Errorf("Close DB failed: %s", err2)
		}
		return err
	}
	return nil
}

func (s *Server) listenAndServe() (err error) {

	var listener net.Listener

	c := s.Config

	// Set default listening address and port
	if c.Address == "" {
		c.Address = DefaultServerAddr
	}
	if c.Port == 0 {
		c.Port = DefaultServerPort
	}
	addr := net.JoinHostPort(c.Address, strconv.Itoa(c.Port)) // 将 host 和 post 组装成地址（socket）
	var addrStr string
	addrStr = fmt.Sprintf("http://%s", addr)
	listener, err = net.Listen("tcp", addr) // 监听 socket tcp
	if err != nil {
		return errors.Wrapf(err, "TCP listen failed for %s", addrStr)
	}

	s.listener = listener
	log.Infof("Listening on %s", addrStr)

	// Start serving requests, either blocking or non-blocking
	if s.BlockingStart {
		return s.serve()
	}
	s.wait = make(chan bool)
	go s.serve()

	return nil
}
func (s *Server) serve() error {
	listener := s.listener
	if listener == nil {
		// This can happen as follows:
		// 1) listenAndServe above is called with s.BlockingStart set to false
		//    and returns to the caller
		// 2) the caller immediately calls s.Stop, which sets s.listener to nil
		// 3) the go routine runs and calls this function

		return nil
	}
	s.serveError = http.Serve(listener, s.mux)
	log.Errorf("Server has stopped serving: %s", s.serveError)
	s.closeListener()
	err := s.closeDB()
	if err != nil {
		log.Errorf("Close DB failed: %s", err)
	}
	if s.wait != nil {
		s.wait <- true
	}
	return s.serveError
}
func (s *Server) closeListener() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	port := s.Config.Port
	if s.listener == nil {
		msg := fmt.Sprintf("Stop: listener was already closed on port %d", port)
		log.Debugf(msg)
		return errors.New(msg)
	}
	err := s.listener.Close()
	s.listener = nil
	if err != nil {
		log.Debugf("Stop: failed to close listener on port %d: %s", port, err)
		return err
	}
	log.Debugf("Stop: successfully closed listener on port %d", port)
	return nil
}
func (s *Server) registerHandlers() {
	s.mux = gmux.NewRouter()
	s.registerHandler("testinfo", newTestInfoEndpoint(s))          // 生成测试交易存入redis
	s.registerHandler("recTransHash", receive_trans_bitarry(s))    // 确定公共交易集，建块
	s.registerHandler("recBlock", receiveBlock(s))                 // 接收Block，校验后vote
	s.registerHandler("recBlockVoteRound1", recBlockVoteRound1(s)) // 接收vote1，票数达到要求后投票
	s.registerHandler("recBlockVoteRound2", recBlockVoteRound2(s)) // 接收vote2，票数达到要求后落库（写入数组）
	s.registerHandler("receiveTx", receiveTx(s)) // 接收用户交易
	s.registerHandler("newAccount", newAccount(s)) // 新建账户
	s.registerHandler("newTransaction", sendCoins(s)) // 交易构建
	s.registerHandler("queryTransactions", getTransaction(s)) // 交易查询

}
func (s *Server) registerHandler(path string, se *serverEndpoint) {
	s.mux.Handle("/"+path, se)

}

// Stop the server
// WARNING: This forcefully closes the listening socket and may cause
// requests in transit to fail, and so is only used for testing.
// A graceful shutdown will be supported with golang 1.8.
func (s *Server) Stop() error {
	err := s.closeListener()
	if err != nil {
		return err
	}
	if s.wait == nil {
		return nil
	}
	// Wait for message on wait channel from the http.serve thread. If message
	// is not received in 10 seconds, return
	port := s.Config.Port
	for i := 0; i < 10; i++ {
		select {
		case <-s.wait:
			log.Debugf("Stop: successful stop on port %d", port)
			close(s.wait)
			s.wait = nil
			return nil
		default:
			log.Debugf("Stop: waiting for listener on port %d to stop", port)
			time.Sleep(time.Second)
		}
	}
	log.Debugf("Stop: timed out waiting for stop notification for port %d", port)
	// make sure DB is closed
	err = s.closeDB()
	if err != nil {
		log.Errorf("Close DB failed: %s", err)
	}
	return nil
}
