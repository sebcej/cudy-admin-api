package api

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ConnectedDevice struct {
	Name string
	Type string
	IP   string
	Mac  string

	Upload   string
	Download string

	Signal             string
	ConnectionDuration string
}

type ConnectedDevicesResponse struct {
	Devices []ConnectedDevice
	Count   int
}

func (c *Config) ConnectedDevices() (resp *ConnectedDevicesResponse, err error) {
	response, err := c.sessionApiCall("/cgi-bin/luci/admin/network/devices/devlist?detail=1", nil, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	var connectedDevices []ConnectedDevice

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		deviceName := s.Find("td").Eq(1).Find("p").Eq(0).Contents().Not("span").Text()

		if deviceName == "" {
			return
		}

		var ip string
		var mac string
		var upload string
		var download string

		s.Find("td.hidden-xs").Eq(3).Find("p").Eq(0).Contents().Not("br").Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "#text" {
				if i == 0 {
					ip = s.Text()
				} else if i == 1 {
					mac = s.Text()
				}
			}
		})

		s.Find("td.hidden-xs").Eq(4).Find("p").Eq(0).Contents().Not("br").Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "#text" {
				if i == 1 {
					upload = s.Text()
				} else if i == 3 {
					download = s.Text()
				}
			}
		})

		connectedDevices = append(connectedDevices, ConnectedDevice{
			Name: deviceName,
			Type: s.Find("td").Eq(1).Find("p span").Eq(0).Text(),
			IP:   ip,
			Mac:  mac,

			Upload:   strings.TrimSpace(upload),
			Download: strings.TrimSpace(download),

			Signal:             s.Find("td.hidden-xs").Eq(5).Find("p").Eq(0).Text(),
			ConnectionDuration: s.Find("td.hidden-xs").Eq(6).Find("p").Eq(0).Text(),
		})
	})

	return &ConnectedDevicesResponse{
		Devices: connectedDevices,
		Count:   len(connectedDevices),
	}, nil
}
