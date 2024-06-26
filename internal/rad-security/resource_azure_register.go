package rad_security

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rad-security/terraform-provider-rad-security/internal/request"
)

func resourceAzureRegister() *schema.Resource {
	return &schema.Resource{
		Description: "Register Azure Subscription and Tenant with Rad Security",

		CreateContext: resourceAzureRegisterCreate,
		ReadContext:   resourceAzureRegisterRead,
		UpdateContext: resourceAzureRegisterUpdate,
		DeleteContext: resourceAzureRegisterDelete,

		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Description: "Subscription ID to use",
				Required:    true,
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "Azure Tenant to use when gathering resources",
				ForceNew:    true,
				Required:    true,
				Sensitive:   true,
			},
			"rad_security_registered": {
				Type:        schema.TypeBool,
				Description: "Target of the API path",
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

func resourceAzureRegisterCreate(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	config := meta.(*Config)
	httpMethod := http.MethodPost
	setValueOnSuccess := config.RadSecurityApiUrl
	diags = resourceAzureRegisterGeneric(ctx, httpMethod, d, setValueOnSuccess, meta)
	return diags
}

func resourceAzureRegisterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl
	targetURI := apiUrlBase + "/cloud/register"
	err := d.Set("api_path", targetURI)
	if err != nil {
		return diag.Errorf("Error setting api_path: %s", err)
	}
	return nil
}

func resourceAzureRegisterUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Update has not yet been implemented
	return nil
}

func resourceAzureRegisterDelete(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	httpMethod := http.MethodDelete
	setValueOnSuccess := ""
	diags = resourceAzureRegisterGeneric(ctx, httpMethod, d, setValueOnSuccess, meta)
	return diags
}

func resourceAzureRegisterGeneric(ctx context.Context, httpMethod string, d *schema.ResourceData, setValueOnSuccess string, meta any) (diags diag.Diagnostics) {
	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl

	targetURI := apiUrlBase + "/cloud/register"
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	tenantID := d.Get("tenant_id").(string)
	subscriptionId := d.Get("subscription_id").(string)

	payload := &RegistrationPayload{
		Type:                "azure",
		AzureTenantID:       tenantID,
		AzureSubscriptionID: subscriptionId,
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
