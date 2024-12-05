package rad_security

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rad-security/terraform-provider-rad-security/internal/request"
)

func resourceGoogleCloud() *schema.Resource {
	return &schema.Resource{
		Description: "Register Google Cloud Project with Workload Federation",

		CreateContext: resourceGoogleCloudWorkloadFederationCreate,
		ReadContext:   resourceGoogleCloudWorkloadFederationRead,
		UpdateContext: resourceGoogleCloudWorkloadFederationUpdate,
		DeleteContext: resourceGoogleCloudWorkloadFederationDelete,

		Schema: map[string]*schema.Schema{
			"google_cloud_service_account_email": {
				Type:        schema.TypeString,
				Description: "Google Cloud service account to impersonate",
				Required:    true,
			},
			"google_cloud_pool_provider_name": {
				Type:        schema.TypeString,
				Description: "Google Cloud pool provider name",
				Required:    true,
			},
			"google_cloud_project_number": {
				Type:        schema.TypeString,
				Description: "Google Cloud project number to sync resources with",
				Required:    true,
			},
		},
	}
}

func resourceGoogleCloudWorkloadFederationCreate(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	var cloudAccount RegistrationPayload
	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl

	targetURI := apiUrlBase + "/cloud/register"
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	googleCloudServiceAccountEmail := d.Get("google_cloud_service_account_email").(string)
	googleCloudPoolProviderName := d.Get("google_cloud_pool_provider_name").(string)
	googleCloudProjectNumber := d.Get("google_cloud_project_number").(string)

	payload := &RegistrationPayload{
		Type:                     "google",
		GoogleCloudProjectNumber: &googleCloudProjectNumber,
		GoogleCloudWorkloadIdentityPoolProviderName: &googleCloudPoolProviderName,
		GoogleCloudServiceAccountEmail:              &googleCloudServiceAccountEmail,
	}

	statusCode, body, diags := request.AuthenticatedRequest(ctx, apiUrlBase, http.MethodPost, targetURI, accessKey, secretKey, payload)
	if statusCode != http.StatusOK {
		return append(diags, diag.Errorf("Failed to register with Rad Security, received HTTP status: %d", statusCode)...)
	}
	err := json.Unmarshal(body, &cloudAccount)
	if err != nil {
		return diag.Errorf("Error decoding JSON: %s", err)
	}

	compositeID := cloudAccount.ID + ":" + cloudAccount.RadAccountID
	d.SetId(compositeID)

	return diags
}

func resourceGoogleCloudWorkloadFederationRead(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	var cloudAccount RegistrationPayload

	idParts := strings.Split(d.Id(), ":")
	if len(idParts) != 2 {
		return diag.Errorf("Invalid ID format, expected ResourceID:RadAccountID")
	}
	cloudAccountID := idParts[0]

	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl

	targetURI := apiUrlBase + "/cloud/" + cloudAccountID
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	statusCode, body, diags := request.AuthenticatedRequest(ctx, apiUrlBase, http.MethodGet, targetURI, accessKey, secretKey, nil)
	if statusCode != http.StatusOK {
		return append(diags, diag.Errorf("Failed to register with Rad Security, received HTTP status: %d", statusCode)...)
	}

	err := json.Unmarshal(body, &cloudAccount)
	if err != nil {
		return diag.Errorf("Error decoding JSON: %s", err)
	}

	return diags
}

func resourceGoogleCloudWorkloadFederationUpdate(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	var cloudAccount RegistrationPayload
	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl

	idParts := strings.Split(d.Id(), ":")
	if len(idParts) != 2 {
		return diag.Errorf("Invalid ID format, expected ResourceID:RadAccountID")
	}
	cloudAccountID := idParts[0]
	radAccountID := idParts[1]

	targetURI := apiUrlBase + "/cloud/" + cloudAccountID
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	googleCloudServiceAccountEmail := d.Get("google_cloud_service_account_email").(string)
	googleCloudPoolProviderName := d.Get("google_cloud_pool_provider_name").(string)
	googleCloudProjectNumber := d.Get("google_cloud_project_number").(string)

	payload := &RegistrationPayload{
		Type:                     "google",
		GoogleCloudProjectNumber: &googleCloudProjectNumber,
		GoogleCloudWorkloadIdentityPoolProviderName: &googleCloudPoolProviderName,
		GoogleCloudServiceAccountEmail:              &googleCloudServiceAccountEmail,
		ID:                                          cloudAccountID,
		RadAccountID:                                radAccountID,
	}

	statusCode, body, diags := request.AuthenticatedRequest(ctx, apiUrlBase, http.MethodPut, targetURI, accessKey, secretKey, payload)
	if statusCode != http.StatusOK {
		return append(diags, diag.Errorf("Failed to register with Rad Security, received HTTP status: %d", statusCode)...)
	}
	err := json.Unmarshal(body, &cloudAccount)
	if err != nil {
		return diag.Errorf("Error decoding JSON: %s", err)
	}

	return diags
}

func resourceGoogleCloudWorkloadFederationDelete(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	idParts := strings.Split(d.Id(), ":")
	if len(idParts) != 2 {
		return diag.Errorf("Invalid ID format, expected ResourceID:RadAccountID")
	}
	cloudAccountID := idParts[0]

	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl

	targetURI := apiUrlBase + "/cloud/" + cloudAccountID
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	statusCode, _, diags := request.AuthenticatedRequest(ctx, apiUrlBase, http.MethodGet, targetURI, accessKey, secretKey, nil)
	if statusCode != http.StatusOK {
		return append(diags, diag.Errorf("Failed to delete cloud account registration with Rad Security, received HTTP status: %d", statusCode)...)
	}
	d.SetId("")

	return diags
}
