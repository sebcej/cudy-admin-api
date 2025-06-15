package api

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GcomStatusResponse struct {
	Connected      bool
	NetworkType    string
	SignalStrength int // 0 when unwkown, 4 when max
	Uploaded       string
	Downloaded     string
}

func (c *Config) GcomStatus() (resp *GcomStatusResponse, err error) {
	response, err := c.sessionApiCall("/cgi-bin/luci/admin/network/gcom/status", nil)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	signalStrength := 0
	doc.Find("table td").Eq(1).Find("i.icon").Each(func(i int, s *goquery.Selection) {
		classes := s.AttrOr("class", "")

		if strings.Contains(classes, "icon-4g") {
			reg := regexp.MustCompile(`icon-4g(\d)`)
			res := reg.FindStringSubmatch(classes)[1]
			signalStrength, _ = strconv.Atoi(res)
		}
	})

	usageString := doc.Find("table tr").Eq(3).Find("td").Eq(1).Find("p").Eq(0).Text()
	splittedUsageString := strings.Split(usageString, "/")

	if len(splittedUsageString) != 2 {
		return nil, ErrParseError
	}

	return &GcomStatusResponse{
		Connected:      doc.Find("table th").Eq(1).HasClass("text-success"),
		NetworkType:    doc.Find("table td").Eq(1).Find("p").Eq(0).Text(),
		SignalStrength: signalStrength,
		Uploaded:       strings.TrimSpace(splittedUsageString[0]),
		Downloaded:     strings.TrimSpace(splittedUsageString[1]),
	}, nil
}
