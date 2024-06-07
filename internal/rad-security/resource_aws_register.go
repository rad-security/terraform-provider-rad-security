package rad_security

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rad-security/terraform-provider-rad-security/internal/request"
)

func resourceAwsRegister() *schema.Resource {
	return &schema.Resource{
		Description: "Register AWS account with Rad Security",

		CreateContext: resourceAwsRegisterCreate,
		ReadContext:   resourceAwsRegisterRead,
		UpdateContext: resourceAwsRegisterUpdate,
		DeleteContext: resourceAwsRegisterDelete,

		Schema: map[string]*schema.Schema{
			"rad_security_assumed_role_arn": {
				Type:        schema.TypeString,
				Description: "Rad Security Role to Trust",
				Required:    true,
			},
			"aws_account_id": {
				Type:        schema.TypeString,
				Description: "Rad Security Customer AWS account ID",
				ForceNew:    true,
				Required:    true,
				Sensitive:   true,
			},
			"rad_security_registered": {
				Type:        schema.TypeBool,
				Description: "Tracks if the account has been successfully registered",
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

func resourceAwsRegisterCreate(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	config := meta.(*Config)
	httpMethod := http.MethodPost
	setValueOnSuccess := config.RadSecurityApiUrl
	diags = resourceAwsRegisterGeneric(ctx, httpMethod, d, setValueOnSuccess, meta)
	return diags
}

func resourceAwsRegisterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl
	targetURI := apiUrlBase + "/cloud/register"
	err := d.Set("api_path", targetURI)
	if err != nil {
		return diag.Errorf("Error setting api_path: %s", err)
	}
	return nil
}

func resourceAwsRegisterUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Update has not yet been implemented
	return nil
}

func resourceAwsRegisterDelete(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	httpMethod := http.MethodDelete
	setValueOnSuccess := ""
	diags = resourceAwsRegisterGeneric(ctx, httpMethod, d, setValueOnSuccess, meta)
	return diags
}

func resourceAwsRegisterGeneric(ctx context.Context, httpMethod string, d *schema.ResourceData, setValueOnSuccess string, meta any) (diags diag.Diagnostics) {
	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl

	targetURI := apiUrlBase + "/cloud/register"
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey
	awsAccountID := d.Get("aws_account_id").(string)

	payload := &RegistrationPayload{
		Type:         "aws",
		AWSAccountID: awsAccountID,
		AWSRoleArn:   "arn:aws:iam::" + awsAccountID + ":role/rad-security-connect",
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
