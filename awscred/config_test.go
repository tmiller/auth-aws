package awscred

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var oldCredPath string
var tempPath string

func setup(t *testing.T) {
	oldCredPath = credPath
	path, err := ioutil.TempDir("", "")
	if err != nil {
		credPath = oldCredPath
		t.Fatal(err)
	}
	credPath = fmt.Sprint(path, "/credentials")
}

func teardown(t *testing.T) {
	credPath = oldCredPath
	err := os.RemoveAll(tempPath)
	if err != nil {
		t.Fatal(err)
	}
}

func checkContent(t *testing.T, expected string) {
	c, err := ioutil.ReadFile(credPath)
	if err != nil {
		t.Fatal(err)
	}
	got := string(c)
	if expected != got {
		t.Errorf("\nexpected:\n---\n%s---\ngot:\n---\n%s---", expected, got)
	}
}

func TestNoConfig(t *testing.T) {
	setup(t)
	defer teardown(t)

	expected := `[adfs]
aws_access_key_id     = key
aws_secret_access_key = secret
aws_session_token     = token

`
	c := Credentials{"key", "secret", "token"}
	c.Write()
	checkContent(t, expected)
}
