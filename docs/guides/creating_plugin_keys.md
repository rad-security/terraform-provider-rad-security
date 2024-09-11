---
page_title: "Creating Plugin Keys"
subcategory: "Guides"
description: |-
  Configuring plugin keys
---

# Creating Plugin Keys

Plugin keys are used when connecting a Kubernetes cluster to Rad Security. Normally, plugin keys have to be generated manually within th e UI or through authenticating with the api, and making a request to create a plugin key. The `rad-security_cluster_api_key` resource can be used to create a plugin key through Terraform.

```hcl
resource "rad-security_cluster_api_key" "example" {}
```

The `access_key` and `secret_key` outputs can be used directly in the plugins helm chart or in a secret for the helm chart to use.

If the plugin key gets deleted, the resource will be recreated.
