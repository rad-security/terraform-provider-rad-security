package rad_security

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"rad_security_api_url": {
				Type:        schema.TypeString,
				Description: "Rad Security API to target. Defaults to https://api.rad.security",
				Default:     "https://api.rad.security",
				Optional:    true,
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Description: "Rad Security Customer Access ID",
				ForceNew:    true,
				Required:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "Rad Security Customer Secret Key",
				ForceNew:    true,
				Required:    true,
				Sensitive:   true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rad-security_aws_register":          resourceAwsRegister(),
			"rad-security_azure_register":        resourceAzureRegister(),
			"rad-security_cluster_api_key":       resourceClusterAPIKey(),
			"rad-security_google_cloud_register": resourceGoogleCloud(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

type Config struct {
	RadSecurityApiUrl    string
	RadSecurityAccountId string
	AccessKeyId          string
	SecretKey            string
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		RadSecurityApiUrl: d.Get("rad_security_api_url").(string),
		AccessKeyId:       d.Get("access_key_id").(string),
		SecretKey:         d.Get("secret_key").(string),
	}

	return &config, nil
}

type RegistrationPayload struct {
	Type                                        string  `json:"type"`
	ID                                          string  `json:"id"`
	AWSAccountID                                string  `json:"aws_account_id"`
	AWSRoleArn                                  string  `json:"aws_role_arn"`
	AzureSubscriptionID                         string  `json:"azure_subscription_id"`
	AzureTenantID                               string  `json:"azure_tenant_id"`
	AzureServicePrincipalID                     string  `json:"azure_service_principal_id"`
	AzureServicePrincipalSecret                 string  `json:"azure_service_principal_secret"`
	GoogleCloudProjectNumber                    *string `json:"google_cloud_project_number,omitempty"`
	GoogleCloudWorkloadIdentityPoolProviderName *string `json:"google_cloud_workload_identity_pool_provider_name,omitempty"`
	GoogleCloudServiceAccountEmail              *string `json:"google_cloud_service_account_email,omitempty"`
	RadAccountID                                string  `json:"account_id"`
}
