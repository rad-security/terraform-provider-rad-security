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
			"service_principal_token_id": {
				Type:        schema.TypeString,
				Description: "Optional: Service Principal Token ID to use when authenticating  with token id and secret. OIDC based auth is the preferred option as it is more secure.",
				Optional:    true,
			},
			"service_principal_token_secret": {
				Type:        schema.TypeString,
				Description: "Optional: Service Principal Token Secret to use when authenticating with token id and secret. OIDC based auth is the preferred option as it is more secure.",
				Optional:    true,
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
	servicePrincipalTokenID := d.Get("service_principal_token_id").(string)
	servicePrincipalTokenSecret := d.Get("service_principal_token_secret").(string)

	if servicePrincipalTokenID == "" && servicePrincipalTokenSecret != "" {
		return append(diags, diag.Errorf("Service Principal Token Secret cannot be set when Service Principal Token ID is")...)
	}

	if servicePrincipalTokenSecret == "" && servicePrincipalTokenID != "" {
		return append(diags, diag.Errorf("Service Principal Token ID cannot be set when Service Principal Token Secret is")...)
	}

	payload := &RegistrationPayload{
		Type:                        "azure",
		AzureTenantID:               tenantID,
		AzureSubscriptionID:         subscriptionId,
		AzureServicePrincipalID:     servicePrincipalTokenID,
		AzureServicePrincipalSecret: servicePrincipalTokenSecret,
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
