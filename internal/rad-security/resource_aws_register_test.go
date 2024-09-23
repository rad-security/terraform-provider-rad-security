package rad_security

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/rad-security/terraform-provider-rad-security/internal/auth"
)

func TestAccResourceAWSRegister(t *testing.T) {
	radAccessKeyID := "test-access-key-id"
	radSecretKey := "test-secret-key"
	radSecurityAssumedRoleArn := "arn:aws:iam::123456789012:role/test-role"
	awsAccountID := "123456789012"

	resourceName := "rad-security_aws_register.test"

	response := &RegistrationPayload{
		Type:         "aws",
		AWSAccountID: awsAccountID,
		AWSRoleArn:   radSecurityAssumedRoleArn,
	}

	mockServer := testAccCloudRegisterHttpMock(radAccessKeyID, radSecretKey, response)
	defer mockServer.Close()

	radAuth := auth.New(mockServer.URL)
	providerFactories := setupRadSecurityProvider()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAWSRegisterCreate(radAuth.ApiURL, radAccessKeyID, radSecretKey, radSecurityAssumedRoleArn, awsAccountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rad_security_registered", "true"),
				),
			},
		},
	})
}

func testAccResourceAWSRegisterCreate(apiURL, radAccessKeyID, radSecretKey, radSecurityAssumedRoleArn, awsAccountID string) string {
	return fmt.Sprintf(`
provider "rad-security" {
  rad_security_api_url = "%s"
  access_key_id        = "%s"
  secret_key           = "%s"
}

resource "rad-security_aws_register" "test" {
  rad_security_assumed_role_arn = "%s"
  aws_account_id       = "%s"
}
`, apiURL, radAccessKeyID, radSecretKey, radSecurityAssumedRoleArn, awsAccountID)
}
