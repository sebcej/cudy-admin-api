package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

func (c *Config) apiCall(path string, body io.Reader, extraHeaders *Headers) (response *http.Response, err error) {
	if c.Api.Sm.session == nil {
		c.Api.Sm.session, err = cookiejar.New(nil)
		if err != nil {
			panic(ErrFetchError)
		}
	}

	client := &http.Client{
		Jar: c.Api.Sm.session,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	pre := ""
	if !strings.HasPrefix(c.Api.Host, "http") {
		pre = "http://"

		if c.Api.SecureConnection {
			pre = "https://"
		}
	}

	method := "GET"
	if body != nil {
		method = "POST"
	}

	req, _ := http.NewRequest(method, pre+c.Api.Host+path, body)

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Host", c.Api.Host)
	req.Header.Add("Referer", pre+c.Api.Host+"/cgi-bin/luci/")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	for key, value := range c.Api.ExtraHeaders {
		req.Header.Add(key, value)
	}
	if extraHeaders != nil {
		for key, value := range *extraHeaders {
			req.Header.Add(key, value)
		}
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	response, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return
}

func (c *Config) sessionApiCall(path string, body io.Reader, extraHeaders *Headers) (response *http.Response, err error) {
	attempts := 0

	for {
		response, err = c.apiCall(path, body, extraHeaders)
		if err != nil {
			return response, err
		}

		// We tried, lets throw an error
		if attempts >= c.Api.Sm.MaxRetries {
			return response, err
		}

		// Session expired, let's login again
		if response.StatusCode == 403 {
			attempts++
			err := c.Login()

			if err != nil {
				return nil, err
			}

			continue
		}

		// Unknown error, let's retry
		if response.StatusCode > 300 || response.StatusCode < 200 {
			attempts++

			time.Sleep(1 * time.Second)

			continue
		}

		return
	}
}

func hashPassword(password, salt string) string {
	// Concatenate password and salt
	input := password + salt

	// Create SHA-256 hash
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)

	// Convert to hexadecimal string
	return hex.EncodeToString(hashBytes)
}

func rateLabel(bytes float64) string {
	uby := "Kbps"
	kby := bytes * 8 / 1024

	if kby >= 1024 {
		uby = "Mbps"
		kby = kby / 1024
	}

	return fmt.Sprintf("%.1f %s", kby, uby)
}
