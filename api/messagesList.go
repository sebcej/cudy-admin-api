package api

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type GetMessagesListResponse struct {
	Messages []Message
}

// Box can be:
// rec: received
// sto: sent
func (c *Config) MessagesList(box string, iface string) (resp *GetMessagesListResponse, err error) {
	response, err := c.sessionApiCall("/cgi-bin/luci/admin/network/gcom/sms/smslist?iface="+iface+"&smsbox="+box, nil, nil)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	var messages []Message

	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		id := s.Find(".btn-primary").AttrOr("onclick", "")

		reg := regexp.MustCompile(`cfg=(.+)(&"?)`)
		sub := reg.FindStringSubmatch(id)
		if len(sub) > 0 {
			id = sub[1]
		}

		messages = append(messages, Message{
			api:    c,
			iface:  iface,
			smsbox: box,

			ID:          id,
			PhoneNumber: s.Find("[id$='-phone'] p").First().Text(),
			Preview:     s.Find("[id$='-content'] p").First().Text(),
			CreatedAt:   s.Find("[id$='-timestamp'] p").First().Text(),
		})
	})

	return &GetMessagesListResponse{
		Messages: messages,
	}, nil
}
