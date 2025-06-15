package api

import (
	"net/http"
	"sync"
)

type Credentials struct {
	Timezone string
	Username string
	Password string
}

type LoginMethod int

type SessionManager struct {
	session    http.CookieJar
	MaxRetries int
	mu         sync.Mutex
}

type Config struct {
	Host             string
	ExtraHeaders     map[string]string
	Credentials      Credentials
	LoginMethod      LoginMethod
	SecureConnection bool
	Sm               SessionManager
}

const (
	LoginMethodDefault LoginMethod = iota
	LoginMethodUnencrypted
)
