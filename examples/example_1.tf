provider "rad_security" {
  access_key_id = "rad_security_access_key"
  secret_key    = "rad_security_secret_key"
}

resource "rad-security_aws_register" "example" {
  rad_security_assumed_role_arn = "arn:aws:iam::<aws_account_number>:role/rad-security-connector"
  aws_account_id                = "aws_account_id"
}

resource "rad-security_azure_register" "example" {
  subscription_id = "123"
  tenant_id       = "456"
}
resource "rad-security_cluster_api_key" "example" {}



