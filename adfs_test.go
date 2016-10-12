package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func compareADFSConfg(t *testing.T, expected *AdfsClient, actual *AdfsClient) {
	if *expected != *actual {
		t.Errorf(
			"\nexp: %v\nact: %v",
			expected,
			actual,
		)
	}
}

//func captureStdinStdout(t *testing.T, input string, f func()) string {
//	oldStdout := os.Stdout
//	oldStdin := os.Stdin
//
//	rIn, wIn, err := os.Pipe()
//	if err != nil {
//		t.Fatal(err)
//	}
//	os.Stdin = rIn
//
//	rOut, wOut, err := os.Pipe()
//	if err != nil {
//		t.Fatal(err)
//	}
//	os.Stdout = wOut
//
//	wIn.WriteString(input)
//
//	f()
//
//	outC := make(chan string)
//
//	go func() {
//		var buf bytes.Buffer
//		io.Copy(&buf, rOut)
//		outC <- buf.String()
//	}()
//
//	os.Stdout = oldStdout
//	wOut.Close()
//
//	out := <-outC
//
//	os.Stdin = oldStdin
//	wIn.Close()
//
//	return out
//}

func TestLoadSettingsFile(t *testing.T) {
	expected := &AdfsClient{"foo", "bar", "adfs.test"}

	settingsFile := strings.NewReader(
		fmt.Sprintf(
			"[adfs]\n%s\n%s\n%s",
			"user = foo",
			"pass = bar",
			"host = adfs.test",
		),
	)

	actual := new(AdfsClient)

	actual.loadSettingsFile(settingsFile)
	compareADFSConfg(t, expected, actual)
}

func TestLoadEnvVars(t *testing.T) {
	expected := &AdfsClient{"foo", "bar", "adfs.test"}

	os.Setenv("ADFS_USER", "foo")
	os.Setenv("ADFS_PASS", "bar")
	os.Setenv("ADFS_HOST", "adfs.test")

	actual := new(AdfsClient)
	actual.loadEnvVars()

	compareADFSConfg(t, expected, actual)
}

// func TestLoadAskVars(t *testing.T) {
//
// 	expOut := "Username: Password: Hostname: "
// 	expected := &AdfsClient{"foo", "bar", "adfs.test"}
//
// 	actual := new(AdfsClient)
//
// 	input := "foo\nbar\nadfs.test\n"
// 	actOut := captureStdinStdout(t, input, func() {
// 		loadAskVars(actual)
// 	})
//
// 	if expOut != actOut {
// 		t.Errorf("incorrect output\nexpected: %v\nactual: %v", expOut, actOut)
// 	}
// 	compareADFSConfg(t, expected, actual)
// }
