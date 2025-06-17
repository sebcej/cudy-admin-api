package api

import (
	"net/http"
	"time"
)

type Credentials struct {
	Username string
	Password string
}

type LoginMethod int

type SessionManager struct {
	session http.CookieJar

	MaxRetries int
	RetryWait  time.Duration
}

type Headers map[string]string

type Api struct {
	Host             string
	ExtraHeaders     Headers
	SecureConnection bool
	Sm               SessionManager
}

type Config struct {
	Api         Api
	Credentials Credentials
	LoginMethod LoginMethod
	TimeZone    string
}

type TableRow struct {
	Label     string
	Value     string
	ExtraData map[string]string
}

type Message struct {
	api    *Config
	iface  string
	smsbox string

	ID          string
	PhoneNumber string
	Preview     string
	CreatedAt   string
}

const (
	LoginMethodDefault LoginMethod = iota
	LoginMethodUnencrypted
)
