---
page_title: "Provider Configuration"
subcategory: "Guides"
description: |-
  Configuring the Rad Security Provider
---

# Provider Configuration

The Rad Security provider requires an access key ID and secret key. Note that they have to be cloud api keys, not plugin or normal api keys. These can be retrieved from the Rad Security UI on the Cloud provisioning keys page or through the API. The api call needed to generate a set of cloud api keys is [here](https://docs.rad.security/reference/post_accounts-access-keys). The keys can then be used to authenticate with Rad to manage your Rad Account.

```hcl
provider "rad_security" {
  access_key_id = "rad_security_access_key"
  secret_key    = "rad_security_secret_key"
  rad_security_api_url = "https://api.ksoc.com"
}
```

The cloud api keys are different from the user access keys. This is to provide a more secure way to seperate the responsibilities of the different types of keys. The user access keys are more permissive. The cloud keys only have the permissions necessary to make changes through the resources that provider provides.
