package redis

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/gomodule/redigo/redis"
	"time"
)

//keys that already used
//purpose, keys, type
//transaction pool， transPool， List
//common transaction cache， CommonTxCache， String       map[string][]byte
//common transaction cache for verify, CommonTxCache4verify, Set
//common Tx from nodes， CommonTx+currentBlock.Hash+sender, Set
//received senders, CommonTx+currentBlock.Hash, Set
//verify Block Tx Cache, verifyBlockTxCache+Block.hash, Set
//vote round 1 vote cache, v.Hash+"round1", Set
//vote round 2 vote cache, v.Hash+"round2", Set

type RedisConn struct {
	conn redis.Conn
}

//连接
func redisConn(netWork, address string) (redis.Conn, error) {
	conn, err := redis.Dial(netWork, address)
	if err != nil {
		log.Info(err.Error())
		return nil, err
	}
	return conn, nil
}

//set string
func (conn *RedisConn) SetString(key, value string) error {
	_, err := conn.conn.Do("SET", key, value)
	if err != nil {
		log.Info(err.Error())
		return err
	}
	return nil
}

//get string
func (conn *RedisConn) GetString(key string) (string, error) {
	value, err := redis.String(conn.conn.Do("GET", key))
	if err != nil {
		log.Info(err.Error())
		return "", err
	}
	return value, nil
}

//set out time
func (conn *RedisConn) SetExpire(key string, time int) error {
	_, err := conn.conn.Do("expire", key, time)
	if err != nil {
		log.Info(err.Error())
		return err
	}
	return nil
}

var Pool redis.Pool
var Redisconn redis.Conn

func init() {
	Pool = redis.Pool{
		//最大的激活连接数，同时最多有N个连接
		MaxActive: 0,
		//最大的空闲连接数，即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
		MaxIdle: 10,
		//空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout: 300 * time.Second,
		//
		Wait: true,
		//连接方法
		Dial: func() (redis.Conn, error) {
			return redisConn("tcp", "127.0.0.1:6379")
		},
	}
	conn, err := redisConn("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Info("redis fail init fail:", err.Error())
		//panic(err.Error())

	}
	Redisconn = conn
}

func SADD(keyValue ...string) error {
	_, err := Redisconn.Do("SADD", keyValue[0], keyValue[1])
	if err != nil {
		log.Info(err.Error())
		return err
	}
	return nil
}

func SCARD(key string) (int, error) {
	res, err := redis.Int(Redisconn.Do("SCARD", key))
	if err != nil {
		log.Info(err.Error())
		return 0, err
	}
	return res, nil
}

func ToInt(value interface{}, err error) (int, error) {
	res, err := redis.Int(value, nil)
	if err != nil {
		log.Info("to int err:", err)
		return 0, err
	}
	return res, nil
}
