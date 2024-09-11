---
page_title: "Provider Configuration"
subcategory: "Guides"
description: |-
  Configuring the Rad Security Provider
---

# Provider Configuration

The Rad Security provider requires an access key ID and secret key. These can be retrieved from the Rad Security UI or through the API. The keys can then be used to authenticate with rad to manage your Rad Account.

```hcl
provider "rad_security" {
  access_key_id = "rad_security_access_key"
  secret_key    = "rad_security_secret_key"
}
```
