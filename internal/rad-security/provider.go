package rad_security

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
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
					Sensitive:   true,
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
				"rad-security_aws_register":    resourceAwsRegister(),
				"rad-security_azure_register":  resourceAzureRegister(),
				"rad-security_cluster_api_key": resourceClusterAPIKey(),
			},
			ConfigureContextFunc: configureProvider,
		}
		return p
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
	Type                string `json:"type"`
	AWSAccountID        string `db:"aws_account_id" json:"aws_account_id"`
	AWSRoleArn          string `db:"aws_role_arn" json:"aws_role_arn"`
	AzureSubscriptionID string `db:"azure_subscription_id" json:"azure_subscription_id"`
	AzureTenantID       string `db:"azure_tenant_id" json:"azure_tenant_id"`
}
