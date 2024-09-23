package rad_security

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/rad-security/terraform-provider-rad-security/internal/auth"
)

func TestAccResourceAzureRegister(t *testing.T) {
	radAccessKeyID := "test-access-key-id"
	radSecretKey := "test-secret-key"
	azureTenantID := "test-tenant-id"
	azureSubscriptionID := "test-subscription-id"
	azureServicePrincipalTokenID := "test-service-principal-token-id"
	azureServicePrincipalTokenSecret := "test-service-principal-token-secret"

	resourceName := "rad-security_azure_register.test"

	response := &RegistrationPayload{
		Type:                        "azure",
		AzureTenantID:               azureTenantID,
		AzureSubscriptionID:         azureSubscriptionID,
		AzureServicePrincipalID:     "",
		AzureServicePrincipalSecret: "",
	}

	mockServer := testAccCloudRegisterHttpMock(radAccessKeyID, radSecretKey, response)
	defer mockServer.Close()

	radAuth := auth.New(mockServer.URL)
	providerFactories := setupRadSecurityProvider()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAzureRegisterOIDCCreate(radAuth.ApiURL, radAccessKeyID, radSecretKey, azureSubscriptionID, azureTenantID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subscription_id", azureSubscriptionID),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", azureTenantID),
					resource.TestCheckNoResourceAttr(resourceName, "service_principal_token_id"),
					resource.TestCheckNoResourceAttr(resourceName, "service_principal_token_secret"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAzureRegisterServicePrincipalTokenCreate(radAuth.ApiURL, radAccessKeyID, radSecretKey, azureSubscriptionID, azureTenantID, azureServicePrincipalTokenID, azureServicePrincipalTokenSecret),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subscription_id", azureSubscriptionID),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", azureTenantID),
					resource.TestCheckResourceAttr(resourceName, "service_principal_token_id", azureServicePrincipalTokenID),
					resource.TestCheckResourceAttr(resourceName, "service_principal_token_secret", azureServicePrincipalTokenSecret),
				),
			},
		},
	})

	servicePrincipalTokenIdProvidedRegexError := regexp.MustCompilePOSIX(`(Service Principal Token ID cannot be set when Service Principal Token Secret is)`)
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAzureRegisterServicePrincipalTokenCreate(radAuth.ApiURL, radAccessKeyID, radSecretKey, azureSubscriptionID, azureTenantID, azureServicePrincipalTokenID, ""),
				ExpectError: servicePrincipalTokenIdProvidedRegexError,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subscription_id", azureSubscriptionID),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", azureTenantID),
					resource.TestCheckResourceAttr(resourceName, "service_principal_token_id", azureServicePrincipalTokenID),
					resource.TestCheckResourceAttr(resourceName, "service_principal_token_secret", azureServicePrincipalTokenSecret),
				),
			},
		},
	})

	servicePrincipalTokenSecretProvidedRegexError := regexp.MustCompilePOSIX(`(Service Principal Token Secret cannot be set when Service Principal Token ID is)`)
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAzureRegisterServicePrincipalTokenCreate(radAuth.ApiURL, radAccessKeyID, radSecretKey, azureSubscriptionID, azureTenantID, "", azureServicePrincipalTokenSecret),
				ExpectError: servicePrincipalTokenSecretProvidedRegexError,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subscription_id", azureSubscriptionID),
					resource.TestCheckResourceAttr(resourceName, "tenant_id", azureTenantID),
					resource.TestCheckResourceAttr(resourceName, "service_principal_token_id", azureServicePrincipalTokenID),
					resource.TestCheckResourceAttr(resourceName, "service_principal_token_secret", azureServicePrincipalTokenSecret),
				),
			},
		},
	})
}

func testAccResourceAzureRegisterOIDCCreate(apiURL, radAccessKeyID, radSecretKey, azureSubscriptionId, azureTenantId string) string {
	return fmt.Sprintf(`
provider "rad-security" {
  rad_security_api_url = "%s"
  access_key_id        = "%s"
  secret_key           = "%s"
}

resource "rad-security_azure_register" "test" {
  subscription_id = "%s"
  tenant_id       = "%s"
}
`, apiURL, radAccessKeyID, radSecretKey, azureSubscriptionId, azureTenantId)
}

func testAccResourceAzureRegisterServicePrincipalTokenCreate(apiURL, radAccessKeyID, radSecretKey, azureSubscriptionId, azureTenantId, azureServicePrincipalTokenId, azureServicePrincipalTokenSecret string) string {
	return fmt.Sprintf(`
provider "rad-security" {
  rad_security_api_url = "%s"
  access_key_id        = "%s"
  secret_key           = "%s"
}

resource "rad-security_azure_register" "test" {
  subscription_id = "%s"
  tenant_id       = "%s"
  service_principal_token_id = "%s"
  service_principal_token_secret = "%s"
}
`, apiURL, radAccessKeyID, radSecretKey, azureSubscriptionId, azureTenantId, azureServicePrincipalTokenId, azureServicePrincipalTokenSecret)
}
