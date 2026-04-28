package rad_security

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rad-security/terraform-provider-rad-security/internal/request"
)

func resourceOCIRegister() *schema.Resource {
	return &schema.Resource{
		Description: "Register Oracle Cloud Infrastructure (OCI) tenancy with Rad Security",

		CreateContext: resourceOCIRegisterCreate,
		ReadContext:   resourceOCIRegisterRead,
		UpdateContext: resourceOCIRegisterUpdate,
		DeleteContext: resourceOCIRegisterDelete,

		Schema: map[string]*schema.Schema{
			"tenancy_ocid": {
				Type:        schema.TypeString,
				Description: "OCI Tenancy OCID — uniquely identifies the cloud account.",
				ForceNew:    true,
				Required:    true,
			},
			"user_ocid": {
				Type:        schema.TypeString,
				Description: "OCID of the OCI user that Rad Security authenticates as.",
				ForceNew:    true,
				Required:    true,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Description: "Fingerprint of the API signing key registered against the user.",
				Required:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "OCI home region used to sign API calls (e.g. us-ashburn-1).",
				Required:    true,
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "PEM-encoded RSA private key matching the registered API signing key. Sent over TLS to Rad Security and stored in AWS Secrets Manager.",
				Required:    true,
				Sensitive:   true,
			},
			"rad_security_registered": {
				Type:        schema.TypeBool,
				Description: "Tracks if the account has been successfully registered.",
				Computed:    true,
			},

			// Computed values
			"api_path": {
				Type:        schema.TypeString,
				Description: "Target of the API path",
				Computed:    true,
			},
		},
	}
}

func resourceOCIRegisterCreate(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	config := meta.(*Config)
	httpMethod := http.MethodPost
	setValueOnSuccess := config.RadSecurityApiUrl
	diags = resourceOCIRegisterGeneric(ctx, httpMethod, d, setValueOnSuccess, meta)
	return diags
}

func resourceOCIRegisterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl
	targetURI := apiUrlBase + "/cloud/register"
	err := d.Set("api_path", targetURI)
	if err != nil {
		return diag.Errorf("Error setting api_path: %s", err)
	}
	return nil
}

func resourceOCIRegisterUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Update has not yet been implemented
	return nil
}

func resourceOCIRegisterDelete(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	httpMethod := http.MethodDelete
	setValueOnSuccess := ""
	diags = resourceOCIRegisterGeneric(ctx, httpMethod, d, setValueOnSuccess, meta)
	return diags
}

func resourceOCIRegisterGeneric(ctx context.Context, httpMethod string, d *schema.ResourceData, setValueOnSuccess string, meta any) (diags diag.Diagnostics) {
	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl

	targetURI := apiUrlBase + "/cloud/register"
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	tenancyOCID := d.Get("tenancy_ocid").(string)
	userOCID := d.Get("user_ocid").(string)
	fingerprint := d.Get("fingerprint").(string)
	region := d.Get("region").(string)
	privateKey := d.Get("private_key").(string)

	payload := &RegistrationPayload{
		Type:           "oci",
		OCITenancyOCID: &tenancyOCID,
		OCIUserOCID:    &userOCID,
		OCIFingerprint: &fingerprint,
		OCIRegion:      &region,
		OCIPrivateKey:  &privateKey,
	}

	statusCode, _, diags := request.AuthenticatedRequest(ctx, apiUrlBase, httpMethod, targetURI, accessKey, secretKey, payload)
	if statusCode != http.StatusOK {
		return append(diags, diag.Errorf("Failed to register with Rad Security, received HTTP status: %d", statusCode)...)
	}

	err := d.Set("api_path", targetURI)
	if err != nil {
		return diag.Errorf("Error setting api_path: %s", err)
	}

	if err := d.Set("rad_security_registered", statusCode == http.StatusOK); err != nil {
		return append(diags, diag.Errorf("Error setting rad_security_registered: %s", err)...)
	}

	d.SetId(setValueOnSuccess)

	return nil
}
