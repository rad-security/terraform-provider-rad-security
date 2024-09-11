---
page_title: "Connecting with Azure"
subcategory: "Guides"
description: |-
  A guide to connect with Azure
---

# Connecting with Azure

There are two methods to connect with Azure. The first method is using a Service Principal with OIDC. The second method is using ID and Secret. This resource is used to send the necessary details to the Rad Security API to be able to syncronize resources in your Azure subscription.

## Connecting with OIDC

This method uses OIDC for authentication. This is the recommended method.

In your Terraform, use the `rad-security_azure_register` resource:

```hcl
resource "rad-security_azure_register" "example" {
  subscription_id = "your-azure-subscription-id"
  tenant_id       = "your-azure-tenant-id"
}
```

No ID or Secret is needed, and Rad Security will use OIDC to syncronize resources from your Azure subscription.

## Connecting with ID and Secret

This method uses a Service Principal ID and Secret for authentication. With an existing Azure Service Principal. This is not the recommended method. OIDC is recommended as it is more secure and easier to manage.

In your Terraform configuration, use the rad-security_azure_register resource:

```hcl
resource "rad-security_azure_register" "azure_connection" {
  subscription_id           = "your-azure-subscription-id"
  tenant_id                 = "your-azure-tenant-id"
  service_principal_id      = "your-service-principal-id"
  service_principal_secret  = "your-service-principal-secret"
}
```

Replace the placeholder values with your actual Azure details:

- `subscription_id`: Your Azure Subscription ID
- `tenant_id`: Your Azure Tenant ID
- `service_principal_id`: The ID of the Service Principal you created
- `service_principal_secret`: The secret of the Service Principal
