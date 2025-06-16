package cudy_test

import (
	"testing"

	"github.com/sebcej/cudy-admin-api"
	"github.com/sebcej/cudy-admin-api/api"
	testdata "github.com/sebcej/cudy-admin-api/tests"
)

type LoginTest struct {
	name           string
	modifyConfig   func(*api.Config)
	expectedErrMsg string
}

func TestLogin(t *testing.T) {
	hiddenInputs := map[string]string{
		"_csrf": "db71a06b67327ad21fdcd7fed4803102",
		"token": "113139b3c8f685db338c5a132054ae67",
		"salt":  "192d8c8c1d50f05d50f91546272aad69",
	}

	// Create test server
	serverCfg := testdata.Config{
		Username:       "admin",
		Password:       "password",
		HashedPassword: "fb021c647416fe8c73f43d78d45e25bb0d35567fee5a30cb2f2a666571b9de30",
		Salt:           hiddenInputs["salt"],
		Token:          hiddenInputs["token"],
		HiddenInputs:   hiddenInputs,
	}
	testServer := testdata.NewLoginTestServer(t, serverCfg)
	defer testServer.Close()

	// Test cases
	tests := []LoginTest{
		{
			name: "Valid credentials",
			modifyConfig: func(c *api.Config) {
				// Use correct credentials (already set)
			},
		},
		{
			name: "Invalid password",
			modifyConfig: func(c *api.Config) {
				c.Credentials.Password = "wrongPassword"
			},
			expectedErrMsg: api.ErrWrongCredentials.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cudyApi, _ := cudy.Init(testServer.URL, serverCfg.Username, serverCfg.Password)
			tt.modifyConfig(cudyApi)
			err := cudyApi.Login()

			if tt.expectedErrMsg != "" {
				if err == nil || err.Error() != tt.expectedErrMsg {
					t.Errorf("DefaultLogin() error = %v, expectedErrMsg %q", err, tt.expectedErrMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DefaultLogin() error = %v", err)
				}
			}
		})
	}
}
