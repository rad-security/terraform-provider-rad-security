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

func TestResourceGoogleCloudCreate(t *testing.T) {
	radAccessKeyID := "test-access-key-id"
	radSecretKey := "test-secret-key"

	// Test data for Google Cloud registration
	googleServiceAccount := "test-sa@project.iam.gserviceaccount.com"
	googlePoolProvider := "test-pool-provider"
	googleProjectNumber := "123456789"

	resourceName := "rad-security_google_cloud_register.test"

	response := &RegistrationPayload{
		ID:                             "test-registration-id",
		Type:                           "google",
		GoogleCloudServiceAccountEmail: &googleServiceAccount,
		GoogleCloudWorkloadIdentityPoolProviderName: &googlePoolProvider,
		GoogleCloudProjectNumber:                    &googleProjectNumber,
	}

	mockServer := testAccGoogleCloudHttpMock(radAccessKeyID, radSecretKey, response)
	defer mockServer.Close()

	radAuth := auth.New(mockServer.URL)
	providerFactories := setupRadSecurityProvider()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGoogleCloudCreate(
					radAuth.ApiURL,
					radAccessKeyID,
					radSecretKey,
					googleServiceAccount,
					googlePoolProvider,
					googleProjectNumber,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "google_cloud_service_account_email", googleServiceAccount),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_pool_provider_name", googlePoolProvider),
					resource.TestCheckResourceAttr(resourceName, "google_cloud_project_number", googleProjectNumber),
				),
			},
		},
	})
}

func testAccGoogleCloudHttpMock(accessKeyID string, secretKey string, registrationResponse *RegistrationPayload) *httptest.Server {
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
			switch r.Method {
			case http.MethodPost:
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(registrationResponse)
			case http.MethodPut:
				w.WriteHeader(http.StatusNoContent)
				json.NewEncoder(w).Encode(registrationResponse)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		} else if r.URL.Path == "/cloud/"+registrationResponse.ID {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodGet:
				json.NewEncoder(w).Encode(registrationResponse)
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			http.NotFound(w, r)
		}
	}))
}

func testAccResourceGoogleCloudCreate(apiURL, radAccessKeyID, radSecretKey, serviceAccount, poolProvider, projectNumber string) string {
	return fmt.Sprintf(`
provider "rad-security" {
  rad_security_api_url = "%s"
  access_key_id        = "%s"
  secret_key           = "%s"
}

resource "rad-security_google_cloud_register" "test" {
  google_cloud_service_account_email = "%s"
  google_cloud_pool_provider_name    = "%s"
  google_cloud_project_number        = "%s"
}
`, apiURL, radAccessKeyID, radSecretKey, serviceAccount, poolProvider, projectNumber)
}
