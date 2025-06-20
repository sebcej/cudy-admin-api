package api

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (c *Config) Login() (err error) {
	form := url.Values{}
	password := c.Credentials.Password

	if c.LoginMethod == LoginMethodDefault {
		homeResponse, err := c.apiCall("/cgi-bin/luci/", nil, nil)
		if err != nil {
			return ErrFetchError
		}
		defer homeResponse.Body.Close()

		doc, err := goquery.NewDocumentFromReader(homeResponse.Body)
		if err != nil {
			return err
		}

		csrfToken := doc.Find("input[type='hidden'][name='_csrf']").AttrOr("value", "")
		token := doc.Find("input[type='hidden'][name='token']").AttrOr("value", "")
		salt := doc.Find("input[type='hidden'][name='salt']").AttrOr("value", "")

		form.Add("_csrf", csrfToken)
		form.Add("token", token)
		form.Add("salt", salt)

		if salt != "" {
			password = hashPassword(password, salt)

			if token != "" {
				password = hashPassword(password, token)
			}
		}
	}

	form.Add("luci_language", "autp")
	form.Add("luci_username", c.Credentials.Username)
	form.Add("luci_password", password)

	form.Add("timeclock", fmt.Sprintf("%d", time.Now().Unix()))
	form.Add("zonename", c.TimeZone)

	formBody := bytes.NewReader([]byte(form.Encode()))
	headers := &Headers{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	response, err := c.apiCall("/cgi-bin/luci/", formBody, headers)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 302:
		return nil
	case 403:
		return ErrWrongCredentials
	default:
		return ErrUnknownError
	}

}
