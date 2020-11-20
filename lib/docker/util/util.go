package util

import (
	"encoding/json"
	"fmt"

	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/pkg/errors"
)

func ValidateAndReturnAbsConf(configFilePath, homeDir, cmdName string) (string, string, error) {
	var err error
	var homeDirSet bool
	var configFileSet bool

	defaultConfig := GetDefaultConfigFile(cmdName) // Get the default configuration

	if configFilePath == "" {
		configFilePath = defaultConfig // If no config file path specified, use the default configuration file
	} else {
		configFileSet = true
	}

	if homeDir == "" {
		homeDir = filepath.Dir(defaultConfig) // If no home directory specified, use the default directory
	} else {
		homeDirSet = true
	}

	// Make the home directory absolute
	homeDir, err = filepath.Abs(homeDir)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to get full path of config file")
	}
	homeDir = strings.TrimRight(homeDir, "/")

	if configFileSet && homeDirSet {
		log.Warning("Using both --config and --home CLI flags; --config will take precedence")
	}

	if configFileSet {
		configFilePath, err = filepath.Abs(configFilePath)
		if err != nil {
			return "", "", errors.Wrap(err, "Failed to get full path of configuration file")
		}
		return configFilePath, filepath.Dir(configFilePath), nil
	}

	configFile := filepath.Join(homeDir, filepath.Base(defaultConfig)) // Join specified home directory with default config file name
	return configFile, homeDir, nil
}

func GetDefaultConfigFile(cmdName string) string {

	if cmdName == "SSBC-server" {
		var fname = fmt.Sprintf("%s-config.yaml", cmdName)
		// First check home env variables
		home := "."
		return path.Join(home, fname)
	}

	var fname = fmt.Sprintf("%s-config.yaml", cmdName)
	return path.Join(os.Getenv("HOME"), ".SSBC-client", fname)
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Marshal(from interface{}, what string) ([]byte, error) {
	buf, err := json.Marshal(from)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to marshal %s", what)
	}
	return buf, nil
}

func Unmarshal(from []byte, to interface{}, what string) error {
	err := json.Unmarshal(from, to)
	if err != nil {
		return errors.Wrapf(err, "Failed to unmarshal %s", what)
	}
	return nil
}
