package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	errors.CheckError(err)

	saml, err := saml.Parse(decodedSamlResponse)

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
	}
	errors.CheckOk(attrRoleIndex >= 0, "Could not find role attribute")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select a role: ")
	userInput, err := reader.ReadString('\n')
	errors.CheckError(err)
	choice, err := strconv.Atoi(strings.Trim(userInput, "\n"))
	errors.CheckError(err)

	chosenValues := strings.Split(saml.Attrs[attrRoleIndex].Values[choice], ",")
	principalARN := chosenValues[0]
	roleARN := chosenValues[1]

	var duration int64 = 3600
	awsSession := session.New(aws.NewConfig().WithRegion("us-east-1"))
	stsClient := sts.New(awsSession)
	assumeRoleInput := sts.AssumeRoleWithSAMLInput{
		DurationSeconds: &duration,
		PrincipalArn:    &principalARN,
		RoleArn:         &roleARN,
		SAMLAssertion:   &samlAssertion,
	}

	creds, err := stsClient.AssumeRoleWithSAML(&assumeRoleInput)
	errors.CheckError(err)

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
