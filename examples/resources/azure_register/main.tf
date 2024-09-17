resource "rad-security_azure_register" "with_oidc" {
  subscription_id = "123"
  tenant_id       = "456"
}

resource "rad-security_azure_register" "with_service_principal_id_and_secret" {
  subscription_id                = "123"
  tenant_id                      = "456"
  service_principal_token_id     = "789"
  service_principal_token_secret = "000"
}
