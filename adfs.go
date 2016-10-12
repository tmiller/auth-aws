package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/howeyc/gopass"
	"github.com/yhat/scrape"

	"gopkg.in/ini.v1"
)

type AdfsClient struct {
	Username string `ini:"user"`
	Password string `ini:"pass"`
	Hostname string `ini:"host"`
}

var (
	settingsPath string = os.Getenv("HOME") + "/.config/auth-aws/config.ini"
)

func (ac *AdfsClient) loadSettingsFile(settingsFile io.Reader) {
	b, err := ioutil.ReadAll(settingsFile)
	checkError(err)

	cfg, err := ini.Load(b)
	if err == nil {
		err = cfg.Section("adfs").MapTo(ac)
		checkError(err)
	}
}

func (ac *AdfsClient) loadEnvVars() {
	if val, ok := os.LookupEnv("ADFS_USER"); ok {
		ac.Username = val
	}
	if val, ok := os.LookupEnv("ADFS_PASS"); ok {
		ac.Password = val
	}
	if val, ok := os.LookupEnv("ADFS_HOST"); ok {
		ac.Hostname = val
	}
}

func (ac *AdfsClient) loadAskVars() {
	reader := bufio.NewReader(os.Stdin)

	if ac.Username == "" {
		fmt.Printf("Username: ")
		user, err := reader.ReadString('\n')
		checkError(err)
		ac.Username = strings.Trim(user, "\n")
	}
	if ac.Password == "" {
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		checkError(err)
		ac.Password = string(pass[:])
	}
	if ac.Hostname == "" {
		fmt.Printf("Hostname: ")
		host, err := reader.ReadString('\n')
		checkError(err)
		ac.Hostname = strings.Trim(host, "\n")
	}
}

func newAdfsClient() *AdfsClient {

	client := new(AdfsClient)

	if settingsFile, err := os.Open(settingsPath); err == nil {
		defer settingsFile.Close()
		client.loadSettingsFile(settingsFile)
	}

	client.loadEnvVars()
	client.loadAskVars()

	if !strings.HasPrefix(client.Hostname, "https://") {
		client.Hostname = "https://" + client.Hostname
	}

	return client
}

func (ac AdfsClient) login() (*http.Response, error) {
	loginUrl := ac.Hostname + "/adfs/ls/IdpInitiatedSignOn.aspx?loginToRp=urn:amazon:webservices"

	cookieJar, err := cookiejar.New(nil)
	checkError(err)

	client := &http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest("GET", loginUrl, nil)
	checkError(err)

	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	checkError(err)

	inputs := scrape.FindAll(root, inputMatcher)
	form, ok := scrape.Find(root, FormMatcher)
	checkOk(ok, "Can't find form")

	formData := url.Values{}

	for _, n := range inputs {
		name := scrape.Attr(n, "name")
		value := scrape.Attr(n, "value")
		switch {
		case strings.Contains(name, "Password"):
			formData.Set(name, ac.Password)
		case strings.Contains(name, "Username"):
			formData.Set(name, ac.Username)
		default:
			formData.Set(name, value)
		}
	}

	action := ac.Hostname + scrape.Attr(form, "action")
	req, err = http.NewRequest("POST", action, strings.NewReader(formData.Encode()))
	checkError(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	return resp, err
}
