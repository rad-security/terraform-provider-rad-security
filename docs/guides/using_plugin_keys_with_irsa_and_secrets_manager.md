---
page_title: "Using Plugin Keys With IRSA and AWS Secrets Manager"
subcategory: "Guides"
description: |-
  How to use plugin keys stored in AWS Secrets Manager and access themwith IRSA
---

# Using Plugin Keys With IRSA and AWS Secrets Manager

You can use the `rad-security_cluster_api_key` resource to create plugin keys for your cluster. The resource outputs can be used within an AWS Secrets Manager secret. You can pass the secret to the Rad Security Helm Chart to use plugin keys by referencing the AWS Secret Manager secret ARN in the helm chart. 

The AWS Secrets Manager secret has to be in the following format:

```json
{
  "access-key-id": "YOUR_ACCESS_KEY_ID",
  "secret-access-key": "YOUR_SECRET_ACCESS_KEY"
}
```


To avoid passing access and secret IAM keys in the helm chart, you can IRSA instead. This example shows how to configure IRSA for an EKS cluster with the minimum required IAM permissions. Depending on how your terraform code is structured, you will need to adjust the IAM policy or IRSA OIDC providers as needed. 

```terraform
terraform {
  required_providers {
    rad-security = {
      source = "rad-security/rad-security"
    }
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "rad-security" {
  access_key_id        = "YOUR_ACESS_KEY_ID"
  secret_key           = "YOUR_SECRET_KEY"
  rad_security_api_url = "https://api.ksoc.com"
}


locals {
  provider_arn = "YOUR_OIDC_PROVIDER_ARN"
  rad_namespace = "ksoc"
  rad_service_accounts = ["ksoc-sbom", "ksoc-guard", "agent-ksoc-k9", "ksoc-node-agent", "ksoc-sync", "ksoc-watch"]
  namespace_service_accounts = [for sa in local.rad_service_accounts : "${local.rad_namespace}:${sa}"]
}

resource "rad-security_cluster_api_key" "this" {}

resource "aws_secretsmanager_secret" "rad_cluster_secret" {
  name        = "rad-cluster-secret-example"
  description = "RAD cluster secret to store an Cluster API Keys"

}

resource "aws_secretsmanager_secret_version" "rad_cluster_secret" {
  secret_id     = aws_secretsmanager_secret.rad_cluster_secret.id
  secret_string = jsonencode({
    access-key-id = "${rad-security_cluster_api_key.this.access_key}",
    secret-key = "${rad-security_cluster_api_key.this.secret_key}"
  })
}

resource "aws_iam_policy" "rad_secret_read_access" {
  name        = "rad-secret-read-access"
  path        = "/"
  description = "Rad API Cluster Key secret policy to allow reading from Secrets Manager"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ]
        Resource = "${aws_secretsmanager_secret.rad_cluster_secret.arn}"
      }
    ]
  })
}

module "iam_eks_role" {
  source    = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  role_name = "rad-plugins-irsa-role"
  create_role = true

  role_policy_arns = {
    policy              = "${aws_iam_policy.rad_secret_read_access.arn}"
  }

  oidc_providers = {
    one = {
      provider_arn               = local.provider_arn
      namespace_service_accounts = local.namespace_service_accounts
    }
  }
}
```

The following values need to be set in the helm chart:

```yaml
ksocGuard:
    serviceAccountAnnotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::YOUR_AWS_ACCOUNT_ID:role/rad-plugins-irsa-role
ksocSbom:
    serviceAccountAnnotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::YOUR_AWS_ACCOUNT_ID:role/rad-plugins-irsa-role
ksocSync:
    serviceAccountAnnotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::YOUR_AWS_ACCOUNT_ID:role/rad-plugins-irsa-role
ksocWatch: 
    serviceAccountAnnotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::YOUR_AWS_ACCOUNT_ID:role/rad-plugins-irsa-role
ksocNodeAgent:
    serviceAccountAnnotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::YOUR_AWS_ACCOUNT_ID:role/rad-plugins-irsa-role
k9:
    serviceAccountAnnotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::YOUR_AWS_ACCOUNT_ID:role/rad-plugins-irsa-role
```

These values set the service account annotations for IRSA. The Rad Plugins will use the default configuration chain to retrieve the necessary AWS credentials to authenticate with AWS and use those credentials to pull the plugin keys from AWS Secrets Manager. 
