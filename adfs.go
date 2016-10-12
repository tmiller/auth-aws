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

type ADFSConfig struct {
	Username string `ini:"user"`
	Password string `ini:"pass"`
	Hostname string `ini:"host"`
}

var (
	settingsPath string = os.Getenv("HOME") + "/.config/auth-aws/config.ini"
)

func loadSettingsFile(adfsConfig *ADFSConfig, settingsFile io.Reader) {
	b, err := ioutil.ReadAll(settingsFile)
	checkError(err)

	cfg, err := ini.Load(b)
	if err == nil {
		err = cfg.Section("adfs").MapTo(adfsConfig)
		checkError(err)
	}
}

func loadEnvVars(adfsConfig *ADFSConfig) {
	if val, ok := os.LookupEnv("ADFS_USER"); ok {
		adfsConfig.Username = val
	}
	if val, ok := os.LookupEnv("ADFS_PASS"); ok {
		adfsConfig.Password = val
	}
	if val, ok := os.LookupEnv("ADFS_HOST"); ok {
		adfsConfig.Hostname = val
	}
}

func loadAskVars(adfsConfig *ADFSConfig) {
	reader := bufio.NewReader(os.Stdin)

	if adfsConfig.Username == "" {
		fmt.Printf("Username: ")
		user, err := reader.ReadString('\n')
		checkError(err)
		adfsConfig.Username = strings.Trim(user, "\n")
	}
	if adfsConfig.Password == "" {
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		checkError(err)
		adfsConfig.Password = string(pass[:])
	}
	if adfsConfig.Hostname == "" {
		fmt.Printf("Hostname: ")
		host, err := reader.ReadString('\n')
		checkError(err)
		adfsConfig.Hostname = strings.Trim(host, "\n")
	}
}

func newADFSConfig() *ADFSConfig {

	adfsConfig := new(ADFSConfig)

	if settingsPath != "" {
		if settingsFile, err := os.Open(settingsPath); err == nil {
			defer settingsFile.Close()
			loadSettingsFile(adfsConfig, settingsFile)
		}
	}

	loadEnvVars(adfsConfig)
	loadAskVars(adfsConfig)

	if !strings.HasPrefix(adfsConfig.Hostname, "https://") {
		adfsConfig.Hostname = "https://" + adfsConfig.Hostname
	}

	return adfsConfig
}

func (auth ADFSConfig) login() (*http.Response, error) {
	loginUrl := auth.Hostname + "/adfs/ls/IdpInitiatedSignOn.aspx?loginToRp=urn:amazon:webservices"

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
			formData.Set(name, auth.Password)
		case strings.Contains(name, "Username"):
			formData.Set(name, auth.Username)
		default:
			formData.Set(name, value)
		}
	}

	action := auth.Hostname + scrape.Attr(form, "action")
	req, err = http.NewRequest("POST", action, strings.NewReader(formData.Encode()))
	checkError(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	return resp, err
}
