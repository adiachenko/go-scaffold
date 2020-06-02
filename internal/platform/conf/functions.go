package conf

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// GetEnv returns the value of the given env variable "key"
// and when its value is empty, it returns the specified "defaultVal".
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.TrimSuffix(value, "\n")
	}

	return defaultVal
}

// GetEnvAsInt returns the value of the given env variable "key" as integer
func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// GetEnvAsBool returns the value of the given env variable "key" as boolean
func GetEnvAsBool(name string, defaultVal bool) bool {
	valueStr := GetEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Load reads the configuration file and assigns values from .env
func Load() bool {
	_, realpath, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(realpath), "../../../") + "/.env"

	if PathExists(path) {
		// read config as ini file
		cfg, err := ini.InsensitiveLoad(path)
		if err != nil {
			panic(err)
		}

		// read the default section (not using sections in config)
		s, err := cfg.GetSection("")
		if err != nil {
			panic(err)
		}

		// set the key/value pairs on the env
		for k, v := range s.KeysHash() {
			if os.Getenv(strings.ToUpper(k)) == "" {
				if err := os.Setenv(strings.ToUpper(k), v); err != nil {
					logrus.WithError(err).Fatal(err.Error())
				}
			}
		}

		return true
	}

	return false
}

// PathExists returns whether the given file or directory exists or not
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}
