package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func compareADFSConfg(t *testing.T, expected *ADFSConfig, actual *ADFSConfig) {
	if *expected != *actual {
		t.Errorf(
			"\nexp: %v\nact: %v",
			expected,
			actual,
		)
	}
}

func TestLoadSettingsFile(t *testing.T) {
	expected := &ADFSConfig{"foo", "bar", "adfs.test"}

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
	compareADFSConfg(t, expected, actual)
}

func TestLoadEnvVars(t *testing.T) {
	expected := &ADFSConfig{"foo", "bar", "adfs.test"}

	os.Setenv("ADFS_USER", "foo")
	os.Setenv("ADFS_PASS", "bar")
	os.Setenv("ADFS_HOST", "adfs.test")

	actual := new(ADFSConfig)
	loadEnvVars(actual)

	compareADFSConfg(t, expected, actual)
}
