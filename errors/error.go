package errors

import (
	"fmt"
	"os"
)

func Error(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "auth-aws: fatal: %v\n", err)
		os.Exit(111)
	}
}

func Ok(ok bool, message string) {
	if !ok {
		fmt.Fprintf(os.Stderr, "auth-aws: fatal: %v\n", message)
		os.Exit(111)
	}
}
