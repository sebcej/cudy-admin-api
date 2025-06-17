package api

import "github.com/PuerkitoBio/goquery"

type MessageContent string

func (m *Message) Delete() error {
	if m.api == nil || m.iface == "" {
		return ErrUnknownError
	}

	_, err := m.api.sessionApiCall("/cgi-bin/luci/admin/network/gcom/sms/delsms?iface="+m.iface+"&cfg="+m.ID, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) Content() (MessageContent, error) {
	if m.api == nil || m.iface == "" {
		return "", ErrUnknownError
	}

	response, err := m.api.sessionApiCall("/cgi-bin/luci/admin/network/gcom/sms/readsms?iface="+m.iface+"&cfg="+m.ID+"&smsbox="+m.smsbox, nil, nil)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	message := MessageContent(doc.Find("textarea").First().Text())

	return message, nil
}

func (m *Message) Reply(message string) error {
	return m.api.SendMessage(m.PhoneNumber, message, m.iface)
}
