package cudy

import (
	"time"

	"github.com/sebcej/cudy-admin-api/api"
)

func Init(host string, username string, password string) (config *api.Config, err error) {
	config = &api.Config{
		Api: api.Api{
			Host:             host,
			SecureConnection: false,
			Sm: api.SessionManager{
				MaxRetries: 2,
				RetryWait:  1 * time.Second,
			},
		},
		Credentials: api.Credentials{
			Username: username,
			Password: password,
		},
		LoginMethod: api.LoginMethodDefault,
		TimeZone:    "Europe/Rome",
	}

	return config, nil
}
