# Session 4
## Used links
- https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/azure-get-started
- https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-linux?pivots=apt
- https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/app_service_custom_hostname_binding

## Random notes
If azure says multiple accounts with same username, run `az account clear`.

Subscription ID: `2df7cef3-7027-4cfd-9818-49eab2ef376a`.

To deploy current Terraform:
1. `cd .infrastructure`.
2. `terraform apply -var-file .tfvars`.
