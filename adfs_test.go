package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoadSettingsFile(t *testing.T) {
	expected := &ADFSConfig{
		Username: "foo",
		Password: "bar",
		Hostname: "adfs.test",
	}

	settingsFile := strings.NewReader(
		fmt.Sprintf(
			"[adfs]\n%s\n%s\n%s",
			"user = foo",
			"pass = bar",
			"host = adfs.test",
		),
	)

	actual := new(ADFSConfig)

	loadSettingsFile(actual, settingsFile)

	if *expected != *actual {
		t.Errorf(
			"\nexp: %v\nact: %v",
			expected,
			actual,
		)
	}

}
