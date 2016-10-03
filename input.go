package main

import "os"

type AuthVars struct {
	Username string
	Password string
	Hostname string
}

func newAuthVars() *AuthVars {
	authVars := &AuthVars{
		Username: os.Getenv("AD_USER"),
		Password: os.Getenv("AD_PASS"),
		Hostname: os.Getenv("AD_HOST"),
	}

	return authVars
}
