package main

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
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

func TestScrapeLoginPage(t *testing.T) {
	client := &AdfsClient{"foo", "bar", "adfs.test"}

	f, err := os.Open("testdata/login_page.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	expUrlValues := url.Values{
		"__VIEWSTATE":          []string{"viewstatedata"},
		"__VIEWSTATEGENERATOR": []string{"viewstategeneratordata"},
		"__EVENTVALIDATION":    []string{"eventvalidationdata"},
		"__db":                 []string{"15"},
		"ctl00$ContentPlaceHolder1$UsernameTextBox": []string{"foo"},
		"ctl00$ContentPlaceHolder1$PasswordTextBox": []string{"bar"},
		"ctl00$ContentPlaceHolder1$SubmitButton":    []string{"Sign In"},
	}

	expFormAction := client.Hostname + "/adfs/ls/?SAMLRequest=REQUEST"
	actFormAction, actUrlValues := client.scrapeLoginPage(f)

	if expFormAction != actFormAction {
		t.Errorf(
			"Form actions do not match \nexp: %s\nact:%s",
			expFormAction,
			actFormAction,
		)
	}

	if !reflect.DeepEqual(expUrlValues, actUrlValues) {
		t.Errorf(
			"Url values do not match \nexp: %s\nact: %s",
			expUrlValues,
			actUrlValues,
		)
	}
}
