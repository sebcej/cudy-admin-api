package cudy

import (
	"github.com/sebcej/cudy-admin-api/api"
)

func Init(host string, username string, password string) (config *api.Config, err error) {
	config = &api.Config{
		Credentials: api.Credentials{
			Timezone: "Europe/Rome",
			Username: username,
			Password: password,
		},
		Host:             host,
		SecureConnection: false,
		Sm: api.SessionManager{
			MaxRetries: 1,
		},
	}

	config.LoginMethod = api.LoginMethodDefault

	return config, nil
}
