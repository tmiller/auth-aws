package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/yhat/scrape"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func checkError(err error) {
	if err != nil {
		fmt.Printf("auth: fatal: %v\n", err)
		os.Exit(111)
	}
}

func checkOk(ok bool, message string) {
	if !ok {
		fmt.Printf("auth: fatal: %v\n", message)
		os.Exit(111)
	}
}

func main() {

	user := os.Getenv("AD_USER")
	pass := os.Getenv("AD_PASS")
	host := os.Getenv("AD_HOST")

	baseUrl := fmt.Sprintf("https://%s", host)
	loginUrl := fmt.Sprintf("%s/adfs/ls/IdpInitiatedSignOn.aspx?loginToRp=urn:amazon:webservices", baseUrl)

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
			formData.Set(name, pass)
		case strings.Contains(name, "Username"):
			formData.Set(name, user)
		default:
			formData.Set(name, value)
		}
	}

	action := fmt.Sprint(baseUrl, scrape.Attr(form, "action"))
	req, err = http.NewRequest("POST", action, strings.NewReader(formData.Encode()))
	checkError(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	root, err = html.Parse(resp.Body)
	checkError(err)

	input, ok := scrape.Find(root, samlResponseMatcher)
	checkOk(ok, "Can't find input")
	assertion := scrape.Attr(input, "value")
	samlResponse, err := base64.StdEncoding.DecodeString(assertion)
	checkError(err)

	fmt.Printf("%s\n", samlResponse)
}

func samlResponseMatcher(n *html.Node) bool {
	return n.DataAtom == atom.Input && scrape.Attr(n, "name") == "SAMLResponse"
}

func inputMatcher(n *html.Node) bool {
	return n.DataAtom == atom.Input
}

func FormMatcher(n *html.Node) bool {
	return n.DataAtom == atom.Form
}
