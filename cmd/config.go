package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/util"
)

const (
	longName     = "SS Blockchain"
	shortName    = "SSBC server"
	cmdName      = "SSBC-server"
	envVarPrefix = "SSBC_SERVER"
	homeEnvVar   = "SSBC_SERVER_HOME"

)

const (
	defaultCfgTemplate = `#############################################################################
#   This is a configuration file for the SSBC-server command.
#
#   COMMAND LINE ARGUMENTS AND ENVIRONMENT VARIABLES
#   ------------------------------------------------
#   Each configuration element can be overridden via command line
#   arguments or environment variables.  The precedence for determining
#   the value of each element is as follows:
#   1) command line argument
#      Examples:
#      a) --port 443
#         To set the listening port
#   2) environment variable
#      Examples:
#      a) SSBC_SERVER_PORT=443
#         To set the listening port
#   3) configuration file
#   4) default value (if there is one)
#      All default values are shown beside each element below.
#
#   FILE NAME ELEMENTS
#   ------------------
#   The value of all fields whose name ends with "file" or "files" are
#   name or names of other files.
#   The value of each of these fields can be a simple filename, a
#   relative path, or an absolute path.  If the value is not an
#   absolute path, it is interpretted as being relative to the location
#   of this configuration file.
#
#############################################################################

# Version of config file
version: <<<VERSION>>>

# Server's listening port (default: 8000)
port: 8000

# Enables debug logging (default: false)
debug: false

# Size limit of an acceptable CRL in bytes (default: 512000)
crlsizelimit: 512000

#############################################################################
#  The registry section controls how the SSBC-server does two things:
#  1) authenticates enrollment requests which contain a username and password
#     (also known as an enrollment ID and secret).
#  2) once authenticated, retrieves the identity's attribute names and
#     values which the SSBC-server optionally puts into TCerts
#     which it issues for transacting on the blockchain.
#     These attributes are useful for making access control decisions in
#     chaincode.
#############################################################################
registry:
  # Maximum number of times a password/secret can be reused for enrollment
  # (default: -1, which means there is no limit)
  maxenrollments: -1

  # Contains identity information which is used when LDAP is disabled
  identities:
     - name: <<<ADMIN>>>
       pass: <<<ADMINPW>>>
       type: server
       

#############################################################################
#  Database section
#  Supported types are: "sqlite3", "postgres", and "mysql".
#  The datasource value depends on the type.
#  If the type is "sqlite3", the datasource value is a file name to use
#  as the database store.  Since "sqlite3" is an embedded database, it
#  may not be used if you want to run the fabric-ca-server in a cluster.
#  To run the fabric-ca-server in a cluster, you must choose "postgres"
#  or "mysql".
#############################################################################
db:
  type: sqlite3
  datasource: fabric-ca-server.db
  tls:
      enabled: false
      certfiles:
      client:
        certfile:
        keyfile:

#############################################################################
# BCCSP (BlockChain Crypto Service Provider) section is used to select which
# crypto library implementation to use
#############################################################################
bccsp:
    default: SW
    sw:
        hash: SHA2
        security: 256
        filekeystore:
            # The directory used for the software file-based keystore
            keystore: msp/keystore

#############################################################################
# User
# Password
#############################################################################
boot: root:passwd
`
)

// Initialize config
func (s *ServerCmd) configInit() (err error) {
	if !s.configRequired() {
		return nil
	}

	s.cfgFileName, s.homeDirectory, err = util.ValidateAndReturnAbsConf(s.cfgFileName, s.homeDirectory, cmdName)
	if err != nil {
		return err
	}

	log.Debugf("Home directory: %s", s.homeDirectory)

	// If the config file doesn't exist, create a default one
	if !util.FileExists(s.cfgFileName) {
		err = s.createDefaultConfigFile()
		if err != nil {
			return errors.WithMessage(err, "Failed to create default configuration file")
		}
		log.Infof("Created default configuration file at %s", s.cfgFileName)
	} else {
		log.Infof("Configuration file location: %s", s.cfgFileName)
	}

	// Read the config
	s.myViper.AutomaticEnv() // read in environment variables that match
	err = UnmarshalConfig(s.cfg, s.myViper, s.cfgFileName, true)
	if err != nil {
		return err
	}




	return nil
}

func (s *ServerCmd) createDefaultConfigFile() error {
	var user, pass string = "root", "password"

	//up := s.myViper.GetString("boot")

	//if up == "" {
	//	return errors.New("The '-b user:pass' option is required")
	//}
	//ups := strings.Split(up, ":")
	//if len(ups) < 2 {
	//	return errors.Errorf("The value '%s' on the command line is missing a colon separator", up)
	//}
	//if len(ups) > 2 {
	//	ups = []string{ups[0], strings.Join(ups[1:], ":")}
	//}
	//user = ups[0]
	//pass = ups[1]
	//if len(user) >= 1024 {
	//	return errors.Errorf("The identity name must be less than 1024 characters: '%s'", user)
	//}
	//if len(pass) == 0 {
	//	return errors.New("An empty password in the '-b user:pass' option is not permitted")
	//}
	var myhost string
	var err error
	myhost, err = os.Hostname()
	if err != nil {
		return err
	}

	// Do string subtitution to get the default config
	cfg := strings.Replace(defaultCfgTemplate, "<<<VERSION>>>", "v0.1", 1)
	cfg = strings.Replace(cfg, "<<<ADMIN>>>", user, 1)
	cfg = strings.Replace(cfg, "<<<ADMINPW>>>", pass, 1)
	cfg = strings.Replace(cfg, "<<<MYHOST>>>", myhost, 1)


	// Now write the file
	cfgDir := filepath.Dir(s.cfgFileName)
	err = os.MkdirAll(cfgDir, 0755)
	if err != nil {
		return err
	}

	// Now write the file
	return ioutil.WriteFile(s.cfgFileName, []byte(cfg), 0644)
}

func UnmarshalConfig(config interface{}, vp *viper.Viper, configFile string,
	server bool) error {

	vp.SetConfigFile(configFile)
	err := vp.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "Failed to read config file '%s'", configFile)
	}

	err = vp.Unmarshal(config)
	if err != nil {
		return errors.Wrapf(err, "Incorrect format in file '%s'", configFile)
	}

	return nil
}