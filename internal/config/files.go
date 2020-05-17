package config

import (
	"os"
	"path/filepath"
)

var (
	CAFile         = configFile("ca.pem")
	ServerCertFile = configFile("server.pem")
	ServerKeyFile  = configFile("server-key.pem")
)

func configFile(filename string) string {
	if os.Getenv("CONFIG_DIR") != "" {
		return filepath.Join(os.Getenv("CONFIG_DIR"), filename)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".proglog", filename)
}
