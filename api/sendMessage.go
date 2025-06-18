package api

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (c *Config) SendMessage(phoneNumber string, message string, iface string) (err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Get updated CSRF token from sms form
	homeResponse, err := c.sessionApiCall("/cgi-bin/luci/admin/network/gcom/sms/smslist?smsbox=rec&iface="+iface, nil, nil)
	if err != nil {
		return ErrFetchError
	}
	defer homeResponse.Body.Close()

	doc, err := goquery.NewDocumentFromReader(homeResponse.Body)
	if err != nil {
		return err
	}

	token := doc.Find("input[type='hidden'][name='token']").First().AttrOr("value", "")
	if token == "" {
		return ErrParseError
	}

	_ = writer.WriteField("token", token)
	_ = writer.WriteField("timeclock", fmt.Sprintf("%d", time.Now().Unix()))
	_ = writer.WriteField("cbi.submit", "1")
	_ = writer.WriteField("cbid.smsnew.1.phone", phoneNumber)
	_ = writer.WriteField("cbid.smsnew.1.content", message)
	_ = writer.WriteField("cbid.smsnew.1.send", "Send")

	err = writer.Close()
	if err != nil {
		return ErrParseError
	}

	headers := &Headers{
		"Content-Type": writer.FormDataContentType(),
	}
	response, err := c.sessionApiCall("/cgi-bin/luci/admin/network/gcom/sms/smsnew?nomodal=&iface="+iface, body, headers)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return ErrParseError
	}

	return
}
