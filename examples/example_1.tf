provider "rad_security" {
  access_key_id = "rad_security_access_key"
  secret_key    = "rad_security_secret_key"
  rad_security_api_url = "https://api.ksoc.com"
}

resource "rad-security_aws_register" "example" {
  rad_security_assumed_role_arn = "arn:aws:iam::<aws_account_number>:role/rad-security-connector"
  aws_account_id                = "aws_account_id"
}

resource "rad-security_azure_register" "example_with_oidc" {
  subscription_id = "123"
  tenant_id       = "456"
}

resource "rad-security_azure_register" "example_with_service_principal_id_and_secret" {
  subscription_id                = "123"
  tenant_id                      = "456"
  service_principal_token_id     = "789"
  service_principal_token_secret = "000"
}

resource "rad-security_cluster_api_key" "example" {}



