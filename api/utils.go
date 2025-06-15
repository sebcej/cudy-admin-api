package api

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func (c *Config) apiCall(path string, body io.Reader) (response *http.Response, err error) {
	if c.Sm.session == nil {
		c.Sm.session, err = cookiejar.New(nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	client := &http.Client{
		Jar: c.Sm.session,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	pre := ""
	if !strings.HasPrefix(c.Host, "http") {
		pre = "http://"

		if c.SecureConnection {
			pre = "https://"
		}
	}

	method := "GET"
	if body != nil {
		method = "POST"
	}

	req, _ := http.NewRequest(method, pre+c.Host+path, body)

	req.Header.Add("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Host", c.Host)
	req.Header.Add("Origin", pre+c.Host)
	req.Header.Add("Referer", pre+c.Host+"/cgi-bin/luci")

	for key, value := range c.ExtraHeaders {
		req.Header.Add(key, value)
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

func (c *Config) sessionApiCall(path string, body io.Reader) (response *http.Response, err error) {
	attempts := 0

	for {
		c.Sm.mu.Lock()
		response, err = c.apiCall(path, body)
		c.Sm.mu.Unlock()
		if err != nil {
			return response, err
		}

		// We tried, lets throw an error
		if attempts >= c.Sm.MaxRetries {
			return response, err
		}

		// Session expired, let's login again
		if response.StatusCode == 403 {
			c.Sm.mu.Lock()

			attempts++
			err := c.Login()
			c.Sm.mu.Unlock()

			if err != nil {
				return nil, err
			}

			continue
		}

		// Unknown error, let's retry just in case
		if response.StatusCode > 300 || response.StatusCode < 200 {
			attempts++
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
