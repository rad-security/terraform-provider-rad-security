package rad_security

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/rad-security/terraform-provider-rad-security/internal/request"
)

type CreateAccessKeyReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ClusterAPIAccesskey struct {
	ID        string     `json:"id"`
	SecretKey string     `json:"secret_key"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	RevokedBy *string    `json:"revoked_by,omitempty"`
}

func resourceClusterAPIKey() *schema.Resource {
	return &schema.Resource{
		Description: "Create new cluster access keys to use",

		CreateContext: resourceClusterAPIKeyCreate,
		ReadContext:   resourceClusterAPIKeyRead,
		DeleteContext: resourceClusterAPIKeyDelete,

		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Description: "Rad Security Cluster Access Key",
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "Rad Security Cluster Secret Key",
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
			},
			"prefix": {
				Type:        schema.TypeString,
				Description: "Prefix for the access key",
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceClusterAPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	var clusterAPIKeys ClusterAPIAccesskey

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	prefix := d.Get("prefix").(string)

	if prefix != "" {
		prefix = fmt.Sprintf("%s-%s", prefix, string(randomString))
	} else {
		prefix = string(randomString)
	}

	currentTime := time.Now()
	formattedTime := currentTime.Format(time.RFC3339)
	formattedName := fmt.Sprintf("%s Terraform %s", prefix, formattedTime)

	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl

	targetURI := apiUrlBase + "/accounts/access_keys"
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	payload := &CreateAccessKeyReq{
		Type: "plugin",
		Name: formattedName,
	}

	statusCode, body, diags := request.AuthenticatedRequest(ctx, apiUrlBase, http.MethodPost, targetURI, accessKey, secretKey, payload)
	if statusCode != http.StatusCreated {
		return append(diags, diag.Errorf("Failed to register with Rad Security, received HTTP status: %d", statusCode)...)
	}

	err := json.Unmarshal(body, &clusterAPIKeys)
	if err != nil {
		return diag.Errorf("Error decoding JSON: %s", err)
	}

	err = d.Set("access_key", clusterAPIKeys.ID)
	if err != nil {
		return diag.Errorf("Error setting access_key: %s", err)
	}

	err = d.Set("secret_key", clusterAPIKeys.SecretKey)
	if err != nil {
		return diag.Errorf("Error setting secret_key: %s", err)
	}

	d.SetId(clusterAPIKeys.ID)

	return diags
}

func resourceClusterAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var clusterAPIKeys ClusterAPIAccesskey

	cloudAccessKey, ok := d.GetOk("access_key")
	if !ok {
		return diag.Errorf("Missing access_key")
	}

	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	targetURI := apiUrlBase + "/accounts/access_keys/" + cloudAccessKey.(string)

	payload := []byte{}

	statusCode, body, diags := request.AuthenticatedRequest(ctx, apiUrlBase, http.MethodGet, targetURI, accessKey, secretKey, payload)
	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Errorf("Failed to read cluster api keys. received HTTP status: %d", statusCode)...)
	}

	err := json.Unmarshal(body, &clusterAPIKeys)
	if err != nil {
		return diag.Errorf("Decoding JSON: %s", err)
	}

	if clusterAPIKeys.RevokedAt != nil {
		d.SetId("")
	}

	return nil
}

func resourceClusterAPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	cloudAccessKey, ok := d.GetOk("access_key")
	if !ok {
		return diag.Errorf("Missing access_key")
	}

	config := meta.(*Config)
	apiUrlBase := config.RadSecurityApiUrl
	accessKey := config.AccessKeyId
	secretKey := config.SecretKey

	targetURI := apiUrlBase + "/accounts/access_keys/" + cloudAccessKey.(string) + "/revoke"

	payload := []byte{}

	statusCode, _, diags := request.AuthenticatedRequest(ctx, apiUrlBase, http.MethodPut, targetURI, accessKey, secretKey, payload)
	if statusCode != http.StatusNoContent {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return append(diags, diag.Errorf("Failed to delete cluster api keys. received HTTP status: %d", statusCode)...)
	}

	d.SetId("")
	return nil
}
