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

func TestScrapeSamlResponse(t *testing.T) {
	client := &AdfsClient{"foo", "bar", "adfs.test"}

	f, err := os.Open("testdata/login_success.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	expected := "PD94bWwgdmVyc2lvbj0iMS4wIj8+CjxzYW1scDpSZXNwb25zZ" +
		"SB4bWxuczpzYW1scD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOnByb3RvY29sI" +
		"iBJRD0iXzg3YTEzMjk1LTQ3NWItNGRkZS04ZDdjLTk1YjAyZDEyYmZhOCIgVmVyc2lvb" +
		"j0iMi4wIiBJc3N1ZUluc3RhbnQ9IjIwMTYtMTAtMDhUMDU6NDA6NDEuOTAyWiIgRGVzd" +
		"GluYXRpb249Imh0dHBzOi8vc2lnbmluLmF3cy5hbWF6b24uY29tL3NhbWwiIENvbnNlb" +
		"nQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpjb25zZW50OnVuc3BlY2lmaWVkI" +
		"j4KICA8SXNzdWVyIHhtbG5zPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZ" +
		"XJ0aW9uIj5odHRwOi8vYWRmcy5leGFtcGxlL2FkZnMvc2VydmljZXMvdHJ1c3Q8L0lzc" +
		"3Vlcj4KICA8c2FtbHA6U3RhdHVzPgogICAgPHNhbWxwOlN0YXR1c0NvZGUgVmFsdWU9I" +
		"nVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpzdGF0dXM6U3VjY2VzcyIvPgogIDwvc" +
		"2FtbHA6U3RhdHVzPgogIDxBc3NlcnRpb24geG1sbnM9InVybjpvYXNpczpuYW1lczp0Y" +
		"zpTQU1MOjIuMDphc3NlcnRpb24iIElEPSJfNDgzZjNiNGItNzJjMC00YWRjLTlhZTUtZ" +
		"DkyMWI1YmMxOTQxIiBJc3N1ZUluc3RhbnQ9IjIwMTYtMTAtMDhUMDU6NDA6NDEuOTAyW" +
		"iIgVmVyc2lvbj0iMi4wIj4KICAgIDxJc3N1ZXI+aHR0cDovL2FkZnMuZXhhbXBsZS9hZ" +
		"GZzL3NlcnZpY2VzL3RydXN0PC9Jc3N1ZXI+CiAgICA8ZHM6U2lnbmF0dXJlIHhtbG5zO" +
		"mRzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwLzA5L3htbGRzaWcjIj4KICAgICAgPGRzO" +
		"lNpZ25lZEluZm8+CiAgICAgICAgPGRzOkNhbm9uaWNhbGl6YXRpb25NZXRob2QgQWxnb" +
		"3JpdGhtPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiLz4KI" +
		"CAgICAgICA8ZHM6U2lnbmF0dXJlTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53M" +
		"y5vcmcvMjAwMC8wOS94bWxkc2lnI3JzYS1zaGExIi8+CiAgICAgICAgPGRzOlJlZmVyZ" +
		"W5jZSBVUkk9IiNfNDgzZjNiNGItNzJjMC00YWRjLTlhZTUtZDkyMWI1YmMxOTQxIj4KI" +
		"CAgICAgICAgIDxkczpUcmFuc2Zvcm1zPgogICAgICAgICAgICA8ZHM6VHJhbnNmb3JtI" +
		"EFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnI2VudmVsb" +
		"3BlZC1zaWduYXR1cmUiLz4KICAgICAgICAgICAgPGRzOlRyYW5zZm9ybSBBbGdvcml0a" +
		"G09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMTAveG1sLWV4Yy1jMTRuIyIvPgogICAgI" +
		"CAgICAgPC9kczpUcmFuc2Zvcm1zPgogICAgICAgICAgPGRzOkRpZ2VzdE1ldGhvZCBBb" +
		"Gdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNzaGExIi8+C" +
		"iAgICAgICAgICA8ZHM6RGlnZXN0VmFsdWU+RElHRVNUPC9kczpEaWdlc3RWYWx1ZT4KI" +
		"CAgICAgICA8L2RzOlJlZmVyZW5jZT4KICAgICAgPC9kczpTaWduZWRJbmZvPgogICAgI" +
		"CA8ZHM6U2lnbmF0dXJlVmFsdWU+U0lHTkFUVVJFPC9kczpTaWduYXR1cmVWYWx1ZT4KI" +
		"CAgICAgPEtleUluZm8geG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZ" +
		"HNpZyMiPgogICAgICAgIDxkczpYNTA5RGF0YT4KICAgICAgICAgIDxkczpYNTA5Q2Vyd" +
		"GlmaWNhdGU+Q0VSVElGSUNBVEU8L2RzOlg1MDlDZXJ0aWZpY2F0ZT4KICAgICAgICA8L" +
		"2RzOlg1MDlEYXRhPgogICAgICA8L0tleUluZm8+CiAgICA8L2RzOlNpZ25hdHVyZT4KI" +
		"CAgIDxTdWJqZWN0PgogICAgICA8TmFtZUlEIEZvcm1hdD0idXJuOm9hc2lzOm5hbWVzO" +
		"nRjOlNBTUw6Mi4wOm5hbWVpZC1mb3JtYXQ6cGVyc2lzdGVudCI+RE9NQUlOXGZvbzwvT" +
		"mFtZUlEPgogICAgICA8U3ViamVjdENvbmZpcm1hdGlvbiBNZXRob2Q9InVybjpvYXNpc" +
		"zpuYW1lczp0YzpTQU1MOjIuMDpjbTpiZWFyZXIiPgogICAgICAgIDxTdWJqZWN0Q29uZ" +
		"mlybWF0aW9uRGF0YSBOb3RPbk9yQWZ0ZXI9IjIwMTYtMTAtMDhUMDU6NDU6NDEuOTAyW" +
		"iIgUmVjaXBpZW50PSJodHRwczovL3NpZ25pbi5hd3MuYW1hem9uLmNvbS9zYW1sIi8+C" +
		"iAgICAgIDwvU3ViamVjdENvbmZpcm1hdGlvbj4KICAgIDwvU3ViamVjdD4KICAgIDxDb" +
		"25kaXRpb25zIE5vdEJlZm9yZT0iMjAxNi0xMC0wOFQwNTo0MDo0MS44ODZaIiBOb3RPb" +
		"k9yQWZ0ZXI9IjIwMTYtMTAtMDhUMDY6NDA6NDEuODg2WiI+CiAgICAgIDxBdWRpZW5jZ" +
		"VJlc3RyaWN0aW9uPgogICAgICAgIDxBdWRpZW5jZT51cm46YW1hem9uOndlYnNlcnZpY" +
		"2VzPC9BdWRpZW5jZT4KICAgICAgPC9BdWRpZW5jZVJlc3RyaWN0aW9uPgogICAgPC9Db" +
		"25kaXRpb25zPgogICAgPEF0dHJpYnV0ZVN0YXRlbWVudD4KICAgICAgPEF0dHJpYnV0Z" +
		"SBOYW1lPSJodHRwczovL2F3cy5hbWF6b24uY29tL1NBTUwvQXR0cmlidXRlcy9Sb2xlU" +
		"2Vzc2lvbk5hbWUiPgogICAgICAgIDxBdHRyaWJ1dGVWYWx1ZT5NaWxsZXJUPC9BdHRya" +
		"WJ1dGVWYWx1ZT4KICAgICAgPC9BdHRyaWJ1dGU+CiAgICAgIDxBdHRyaWJ1dGUgTmFtZ" +
		"T0iaHR0cHM6Ly9hd3MuYW1hem9uLmNvbS9TQU1ML0F0dHJpYnV0ZXMvUm9sZSI+CiAgI" +
		"CAgICAgPEF0dHJpYnV0ZVZhbHVlPmFybjphd3M6aWFtOjoxMTExMTExMTExMTE6c2Ftb" +
		"C1wcm92aWRlci9BREZTLGFybjphd3M6aWFtOjoxMTExMTExMTExMTE6cm9sZS9BZG1pb" +
		"jwvQXR0cmlidXRlVmFsdWU+CiAgICAgICAgPEF0dHJpYnV0ZVZhbHVlPmFybjphd3M6a" +
		"WFtOjoyMjIyMjIyMjIyMjI6c2FtbC1wcm92aWRlci9BREZTLGFybjphd3M6aWFtOjoyM" +
		"jIyMjIyMjIyMjI6cm9sZS9Vc2VyPC9BdHRyaWJ1dGVWYWx1ZT4KICAgICAgPC9BdHRya" +
		"WJ1dGU+CiAgICAgIDxBdHRyaWJ1dGUgTmFtZT0iaHR0cHM6Ly9hd3MuYW1hem9uLmNvb" +
		"S9TQU1ML0F0dHJpYnV0ZXMvU2Vzc2lvbkR1cmF0aW9uIj4KICAgICAgICA8QXR0cmlid" +
		"XRlVmFsdWU+Mjg4MDA8L0F0dHJpYnV0ZVZhbHVlPgogICAgICA8L0F0dHJpYnV0ZT4KI" +
		"CAgIDwvQXR0cmlidXRlU3RhdGVtZW50PgogICAgPEF1dGhuU3RhdGVtZW50IEF1dGhuS" +
		"W5zdGFudD0iMjAxNi0xMC0wOFQwNTo0MDo0MS41NTlaIiBTZXNzaW9uSW5kZXg9Il80O" +
		"DNmM2I0Yi03MmMwLTRhZGMtOWFlNS1kOTIxYjViYzE5NDEiPgogICAgICA8QXV0aG5Db" +
		"250ZXh0PgogICAgICAgIDxBdXRobkNvbnRleHRDbGFzc1JlZj51cm46b2FzaXM6bmFtZ" +
		"XM6dGM6U0FNTDoyLjA6YWM6Y2xhc3NlczpQYXNzd29yZFByb3RlY3RlZFRyYW5zcG9yd" +
		"DwvQXV0aG5Db250ZXh0Q2xhc3NSZWY+CiAgICAgIDwvQXV0aG5Db250ZXh0PgogICAgP" +
		"C9BdXRoblN0YXRlbWVudD4KICA8L0Fzc2VydGlvbj4KPC9zYW1scDpSZXNwb25zZT4K"

	actual := client.scrapeSamlResponse(f)
	if expected != actual {
		t.Error("Saml respsonses do not match")
	}
}
