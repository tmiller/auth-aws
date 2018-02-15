package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tmiller/auth-aws/awscred"
	"github.com/tmiller/auth-aws/errors"
	"github.com/tmiller/auth-aws/idp"
	"github.com/tmiller/auth-aws/saml"

	"github.com/yhat/scrape"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {

	adfsClient := idp.NewAdfsClient()

	samlAssertion := adfsClient.Login()

	decodedSamlResponse, err := base64.StdEncoding.DecodeString(samlAssertion)
	errors.Error(err)

	saml, err := saml.Parse(decodedSamlResponse)

	var defaultDuration int64 = 3600
	var sessionDuration int64 = 3600
	var sessionNotOnOrAfter time.Time
	var sessionNotOnOrAfterDuration int64 = 3600
	attrRoleIndex := -1
	for ai, attrs := range saml.Attrs {
		if attrs.Name == "https://aws.amazon.com/SAML/Attributes/Role" {
			attrRoleIndex = ai
			for vi, val := range attrs.Values {
				splitVal := strings.Split(val, "/")
				role := splitVal[len(splitVal)-1]
				fmt.Printf("[%d] %v\n", vi, role)
			}
		}
		if attrs.Name == "https://aws.amazon.com/SAML/Attributes/SessionDuration" {
			sessionDuration, err = strconv.ParseInt(saml.Attrs[ai].Values[0], 10, 64)
			if err != nil {
				fmt.Errorf("can't parse the SessionDuration SAML attribute")
			}
		}
		if attrs.Name == "https://aws.amazon.com/SAML/Attributes/SessionNotOnOrAfter" {
			sessionNotOnOrAfter, err = time.Parse(time.RFC3339, saml.Attrs[ai].Values[0])
			if err != nil {
				fmt.Errorf("can't parse the SessionNotOnOrAfter SAML attribute")
			}

			// We could do something like:
			//
			//     sessionNotOnOrAfterDuration = int(time.Until(sessionNotOnOrAfter).Seconds())
			//
			// but Seconds() returns a float64 for some reason and if we cast
			// it to int64 there's some kind of overflow risk.
			sessionNotOnOrAfterDuration = sessionNotOnOrAfter.Unix() - time.Now().Unix()
		}
	}
	errors.Ok(attrRoleIndex >= 0, "Could not find role attribute")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select a role: ")
	userInput, err := reader.ReadString('\n')
	errors.Error(err)
	choice, err := strconv.Atoi(strings.Trim(userInput, "\n"))
	errors.Error(err)

	chosenValues := strings.Split(saml.Attrs[attrRoleIndex].Values[choice], ",")
	principalARN := chosenValues[0]
	roleARN := chosenValues[1]

	// Choose the maximum session expiration value out of the possibilities.  We
	// can only ask for up to 1 hour, but the SAML assertion can provide either
	// of those, and the smallest of the two SAML attributes wins.
	var durationSeconds int64
	durationValues := []int64{defaultDuration, sessionDuration, sessionNotOnOrAfterDuration}
	for _, val := range durationValues {
		durationSeconds = val
		for _, v := range durationValues {
			if v < durationSeconds {
				durationSeconds = v
			}
		}
	}

	awsSession := session.New(aws.NewConfig().WithRegion("us-east-1"))
	stsClient := sts.New(awsSession)

	assumeRoleInput := sts.AssumeRoleWithSAMLInput{
		DurationSeconds: &durationSeconds,
		PrincipalArn:    &principalARN,
		RoleArn:         &roleARN,
		SAMLAssertion:   &samlAssertion,
	}

	creds, err := stsClient.AssumeRoleWithSAML(&assumeRoleInput)
	errors.Error(err)

	awsCredentials := &awscred.Credentials{
		AwsAccessKeyId:     *creds.Credentials.AccessKeyId,
		AwsSecretAccessKey: *creds.Credentials.SecretAccessKey,
		AwsSessionToken:    *creds.Credentials.SessionToken,
	}

	awsCredentials.Write()

}

func samlResponseMatcher(n *html.Node) bool {
	return n.DataAtom == atom.Input && scrape.Attr(n, "name") == "SAMLResponse"
}
