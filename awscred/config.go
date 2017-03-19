package awscred

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

var credPath string = fmt.Sprint(os.Getenv("HOME"), "/.aws/credentials")

func (awsCredentials *Credentials) Write() {

	creds, err := ini.Load(credPath)
	if err != nil {
		creds = ini.Empty()
	}

	creds.NameMapper = ini.TitleUnderscore

	adfs, err := creds.GetSection("adfs")

	if err != nil {
		adfs, err = creds.NewSection("adfs")
		if err != nil {
			fmt.Println(err)
		}
	}

	err = adfs.ReflectFrom(awsCredentials)
	if err != nil {
		fmt.Println(err)
	}

	creds.SaveTo(credPath)
}

type Credentials struct {
	AwsAccessKeyId     string
	AwsSecretAccessKey string
	AwsSessionToken    string
}
