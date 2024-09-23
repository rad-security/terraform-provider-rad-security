package rad_security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/rad-security/terraform-provider-rad-security/internal/auth"
)

func TestAccResourceClusterApiKeyCreate(t *testing.T) {
	radAccessKeyID := "test-access-key-id"
	radSecretKey := "test-secret-key"
	radClusterAccessKeyId := "test-cluster-access-key-id"
	radClusterSecretKey := "test-cluster-secret-key"

	resourceName := "rad-security_cluster_api_key.test"

	response := &ClusterAPIAccesskey{
		ID:        radClusterAccessKeyId,
		SecretKey: radClusterSecretKey,
	}

	mockServer := testAccClusterApiKeyHttpMock(radAccessKeyID, radSecretKey, response)
	defer mockServer.Close()

	radAuth := auth.New(mockServer.URL)
	providerFactories := setupRadSecurityProvider()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceClusterApiKeyCreate(radAuth.ApiURL, radAccessKeyID, radSecretKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_key", radClusterAccessKeyId),
					resource.TestCheckResourceAttr(resourceName, "secret_key", radClusterSecretKey),
				),
			},
		},
	})
}

func testAccClusterApiKeyHttpMock(accessKeyID string, secretKey string, clusterApiKeyResponse *ClusterAPIAccesskey) *httptest.Server {
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
				err := json.NewEncoder(w).Encode(resp)
				if err != nil {
					panic(err)
				}
			} else {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
			}
		} else if r.URL.Path == "/accounts/access_keys" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			err := json.NewEncoder(w).Encode(clusterApiKeyResponse)
			if err != nil {
				panic(err)
			}
			return
		} else if r.URL.Path == "/accounts/access_keys/"+clusterApiKeyResponse.ID {
			err := json.NewEncoder(w).Encode(clusterApiKeyResponse)
			if err != nil {
				panic(err)
			}
			return
		} else if r.URL.Path == "/accounts/access_keys/"+clusterApiKeyResponse.ID+"/revoke" {
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			http.NotFound(w, r)
			return
		}
	}))
}

func testAccResourceClusterApiKeyCreate(apiURL, radAccessKeyID, radSecretKey string) string {
	return fmt.Sprintf(`
provider "rad-security" {
  rad_security_api_url = "%s"
  access_key_id        = "%s"
  secret_key           = "%s"
}

resource "rad-security_cluster_api_key" "test" {}
`, apiURL, radAccessKeyID, radSecretKey)
}
