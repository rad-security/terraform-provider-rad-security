package rad_security

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/rad-security/terraform-provider-rad-security/internal/auth"
)

func TestAccResourceOCIRegister(t *testing.T) {
	radAccessKeyID := "test-access-key-id"
	radSecretKey := "test-secret-key"
	tenancyOCID := "ocid1.tenancy.oc1..test"
	userOCID := "ocid1.user.oc1..test"
	fingerprint := "aa:bb:cc:dd:ee:ff:11:22:33:44:55:66:77:88:99:00"
	region := "us-ashburn-1"
	privateKey := "-----BEGIN RSA PRIVATE KEY-----\ntest-key\n-----END RSA PRIVATE KEY-----\n"

	resourceName := "rad-security_oci_register.test"

	tenancy := tenancyOCID
	user := userOCID
	fp := fingerprint
	reg := region
	pk := privateKey

	response := &RegistrationPayload{
		Type:           "oci",
		OCITenancyOCID: &tenancy,
		OCIUserOCID:    &user,
		OCIFingerprint: &fp,
		OCIRegion:      &reg,
		OCIPrivateKey:  &pk,
	}

	mockServer := testAccCloudRegisterHttpMock(radAccessKeyID, radSecretKey, response)
	defer mockServer.Close()

	radAuth := auth.New(mockServer.URL)
	providerFactories := setupRadSecurityProvider()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOCIRegisterCreate(radAuth.ApiURL, radAccessKeyID, radSecretKey, tenancyOCID, userOCID, fingerprint, region, privateKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tenancy_ocid", tenancyOCID),
					resource.TestCheckResourceAttr(resourceName, "user_ocid", userOCID),
					resource.TestCheckResourceAttr(resourceName, "fingerprint", fingerprint),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "private_key", privateKey),
				),
			},
		},
	})
}

func testAccResourceOCIRegisterCreate(apiURL, radAccessKeyID, radSecretKey, tenancyOCID, userOCID, fingerprint, region, privateKey string) string {
	return fmt.Sprintf(`
provider "rad-security" {
  rad_security_api_url = "%s"
  access_key_id        = "%s"
  secret_key           = "%s"
}

resource "rad-security_oci_register" "test" {
  tenancy_ocid = "%s"
  user_ocid    = "%s"
  fingerprint  = "%s"
  region       = "%s"
  private_key  = <<EOT
%sEOT
}
`, apiURL, radAccessKeyID, radSecretKey, tenancyOCID, userOCID, fingerprint, region, privateKey)
}
