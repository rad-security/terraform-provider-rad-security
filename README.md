# terraform-provider-rad-security
This is the official Terraform Provider for Rad Security. Use this provider to interact with the Rad Security api. The provider can be found on the [Terraform Provider Registery](https://registry.terraform.io/providers/rad-security/rad-security/latest).

To configure the provider, you will need a set of cloud api keys. The keys consist of an access and a secret key that can be generated from the Rad Security platform.

To connect your AWS account to your Rad Security account, create a `rad_security_aws_register` resource where you run terraform for your AWS resources.

An example of leveraging this resource can be found in our terraform module examples directory [here](https://github.com/rad-security/terraform-aws-rad-security-connect/blob/main/examples/main.tf)