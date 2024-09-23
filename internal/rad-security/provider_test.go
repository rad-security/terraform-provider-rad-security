package rad_security

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rad-security/terraform-provider-rad-security/internal/auth"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func setupRadSecurityProvider() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"rad-security": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func testAccCloudRegisterHttpMock(accessKeyID string, secretKey string, registrationResponse *RegistrationPayload) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/authentication/authenticate" {
			var req auth.AuthRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			if req.AccessKeyID == accessKeyID && req.SecretKey == secretKey {
				resp := auth.AuthResponse{Token: "mock_token"}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			} else {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
			}
		} else if r.URL.Path == "/cloud/register" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(registrationResponse)
			return
		} else {
			http.NotFound(w, r)
			return
		}
	}))
}
