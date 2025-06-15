package testdata

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Config holds test server configuration
type Config struct {
	Username       string
	Password       string
	HashedPassword string
	Salt           string
	Token          string
	HiddenInputs   map[string]string
}

// NewTestServer creates an HTTP test server with /cgi-bin/luci routes
func NewLoginTestServer(t *testing.T, cfg Config) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/luci/":
			if r.Method == http.MethodGet {
				// Return HTML with hidden inputs
				html := `<html><body><form>`
				for name, value := range cfg.HiddenInputs {
					html += `<input type="hidden" name="` + name + `" value="` + value + `">`
				}
				html += `</form></body></html>`
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(html))
			} else if r.Method == http.MethodPost {
				// Parse form data
				if err := r.ParseForm(); err != nil {
					t.Errorf("Server: failed to parse form: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// Validate luci_username and luci_password
				username := r.FormValue("luci_username")
				password := r.FormValue("luci_password")

				if username == cfg.Username && password == cfg.HashedPassword {
					w.Header().Set("Location", "/success")
					w.WriteHeader(http.StatusFound) // 302
				} else {
					w.WriteHeader(http.StatusForbidden) // 403
				}
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}
