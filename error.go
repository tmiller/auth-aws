package main

import (
	"fmt"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Printf("auth-aws: fatal: %v\n", err)
		os.Exit(111)
	}
}

func checkOk(ok bool, message string) {
	if !ok {
		fmt.Printf("auth-aws: fatal: %v\n", message)
		os.Exit(111)
	}
}
