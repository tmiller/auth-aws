package main

import "os"

type ADFSConfig struct {
	Username string
	Password string
	Hostname string
}

func newADFSConfig() *ADFSConfig {
	authVars := &ADFSConfig{
		Username: os.Getenv("AD_USER"),
		Password: os.Getenv("AD_PASS"),
		Hostname: os.Getenv("AD_HOST"),
	}

	return authVars
}
