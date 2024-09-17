---
page_title: "Provider: Rad Security"
description: |-
  Rad Security Provider
---

The Rad Security provider is a Terraform provider that allows users to interact with Rad Security through Terraform.

It currently supports cloud registration, and cluster api key management.

## Cloud Registration

The cloud registration resource allows users to register their Azure and AWS accounts to their Rad Security account.

## Cluster API Key Management

The cluster api key management resource allows users to manage the api keys for their cluster.

## Example Usage

An example of how to use the provider can look like the following:

```terraform
provider "rad_security" {
  access_key_id = "rad_security_access_key"
  secret_key    = "rad_security_secret_key"
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
```

Not all resources need to be used together. Depending on the use case, they can be used individually. They are best used within the modules provided by Rad Security.
