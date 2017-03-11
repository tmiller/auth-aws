package idp

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
	"golang.org/x/net/html/atom"

	"github.com/tmiller/auth-aws/errors"

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

func inputMatcher(n *html.Node) bool {
	return n.DataAtom == atom.Input
}

func formMatcher(n *html.Node) bool {
	return n.DataAtom == atom.Form
}

func NewAdfsClient() *AdfsClient {

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

func (ac *AdfsClient) loadSettingsFile(settingsFile io.Reader) {
	b, err := ioutil.ReadAll(settingsFile)
	errors.Error(err)

	cfg, err := ini.Load(b)
	if err == nil {
		err = cfg.Section("adfs").MapTo(ac)
		errors.Error(err)
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
		errors.Error(err)
		ac.Username = strings.Trim(user, "\n")
	}
	if ac.Password == "" {
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		errors.Error(err)
		ac.Password = string(pass[:])
	}
	if ac.Hostname == "" {
		fmt.Printf("Hostname: ")
		host, err := reader.ReadString('\n')
		errors.Error(err)
		ac.Hostname = strings.Trim(host, "\n")
	}
}

func (ac AdfsClient) scrapeLoginPage(r io.Reader) (string, url.Values) {
	root, err := html.Parse(r)
	errors.Error(err)

	inputs := scrape.FindAll(root, inputMatcher)
	form, ok := scrape.Find(root, formMatcher)
	errors.Ok(ok, "Can't find login form")

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

	return action, formData
}

func (ac AdfsClient) scrapeSamlResponse(r io.Reader) string {
	root, err := html.Parse(r)
	errors.Error(err)

	input, ok := scrape.Find(root, samlResponseMatcher)
	errors.Ok(ok, "Can't find saml response")

	return scrape.Attr(input, "value")
}

func (ac AdfsClient) Login() string {
	loginUrl := ac.Hostname + "/adfs/ls/IdpInitiatedSignOn.aspx?loginToRp=urn:amazon:webservices"

	cookieJar, err := cookiejar.New(nil)
	errors.Error(err)

	client := &http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest("GET", loginUrl, nil)
	errors.Error(err)

	resp, err := client.Do(req)
	errors.Error(err)
	defer resp.Body.Close()

	action, formData := ac.scrapeLoginPage(resp.Body)

	req, err = http.NewRequest("POST", action, strings.NewReader(formData.Encode()))
	errors.Error(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	defer resp.Body.Close()

	return ac.scrapeSamlResponse(resp.Body)
}

func samlResponseMatcher(n *html.Node) bool {
	return n.DataAtom == atom.Input && scrape.Attr(n, "name") == "SAMLResponse"
}
