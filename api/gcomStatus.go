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
	SignalStrength int // 0 when unwkown/dead, 4 when max
	Uploaded       string
	Downloaded     string
	ConnectionTime string
	PublicIp       string
	IP             string
	RawValues      []TableRow
}

func (c *Config) GcomStatus(iface ...string) (resp *GcomStatusResponse, err error) {
	path := "/cgi-bin/luci/admin/network/gcom/status?detail=1"
	if len(iface) == 1 {
		path += "&iface" + iface[0]
	}

	response, err := c.sessionApiCall(path, nil)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	signalStrength := 0
	doc.Find("table td").Eq(2).Find("i.icon").Each(func(i int, s *goquery.Selection) {
		classes := s.AttrOr("class", "")

		if strings.Contains(classes, "icon-4g") {
			reg := regexp.MustCompile(`icon-4g(\d)`)
			res := reg.FindStringSubmatch(classes)[1]
			signalStrength, _ = strconv.Atoi(res)
		}
	})

	var rawValues []TableRow

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Find("td").Eq(1).Find("p").Eq(0).Text())
		value := strings.TrimSpace(s.Find("td").Eq(2).Find("p").Eq(0).Text())

		if label == "" && value == "" {
			return
		}

		rawValues = append(rawValues, TableRow{
			Label: label,
			Value: value,
		})
	})

	usageString := rawValues[1].Value
	splittedUsageString := strings.Split(usageString, "/")

	if len(splittedUsageString) != 2 {
		return nil, ErrParseError
	}

	return &GcomStatusResponse{
		Connected:      doc.Find("table th").Eq(1).HasClass("text-success"),
		NetworkType:    doc.Find("table td").Eq(2).Find("p").Eq(0).Text(),
		SignalStrength: signalStrength,
		Uploaded:       strings.TrimSpace(splittedUsageString[0]),
		Downloaded:     strings.TrimSpace(splittedUsageString[1]),
		PublicIp:       rawValues[3].Value,
		IP:             rawValues[4].Value,
		ConnectionTime: rawValues[5].Value,
		RawValues:      rawValues,
	}, nil
}
