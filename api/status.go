package api

import (
	"github.com/PuerkitoBio/goquery"
)

type StatusResponse struct {
	SystemVersion string
	SystemTime    string
	ActivityTime  string
}

// Get status/uptime of router and its firmware version
func (c *Config) Status() (resp *StatusResponse, err error) {
	response, err := c.sessionApiCall("/cgi-bin/luci/admin/system/status", nil, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return &StatusResponse{
		SystemVersion: doc.Find("table th").Eq(1).Text(),
		SystemTime:    doc.Find("table td").Eq(1).Find("p").Eq(0).Text(),
		ActivityTime:  doc.Find("table td").Eq(5).Find("p").Eq(0).Text(),
	}, nil
}
