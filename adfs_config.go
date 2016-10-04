package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type ADFSConfig struct {
	Username string `ini:"user"`
	Password string `ini:"pass"`
	Hostname string `ini:"host"`
}

func newADFSConfig() *ADFSConfig {

	configPath := fmt.Sprintf("%s/.config/auth-aws/config.ini", os.Getenv("HOME"))
	adfsConfig := new(ADFSConfig)

	cfg, err := ini.Load(configPath)
	if err == nil {
		err = cfg.Section("adfs").MapTo(adfsConfig)
		checkError(err)
	}

	if val, ok := os.LookupEnv("ADFS_USER"); ok {
		adfsConfig.Username = val
	}
	if val, ok := os.LookupEnv("ADFS_PASS"); ok {
		adfsConfig.Password = val
	}
	if val, ok := os.LookupEnv("ADFS_HOST"); ok {
		adfsConfig.Hostname = val
	}

	return adfsConfig
}
