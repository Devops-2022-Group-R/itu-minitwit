# Session 4
## Used links
- https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/azure-get-started
- https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-linux?pivots=apt
- https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/app_service_custom_hostname_binding

- https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_server
- https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_firewall_rule
- https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_database

## Random notes
If azure says multiple accounts with same username, run `az account clear`.

Subscription ID: `2df7cef3-7027-4cfd-9818-49eab2ef376a`.

To deploy current Terraform:
1. `cd .infrastructure`.
2. `terraform apply -var-file .tfvars`.

### Identity
Identity on backend did not work. We ran:
```sh
RESOURCE_GROUP="itu-minitwit-rg"
APP_NAME="itu-minitwit-backend-as"

az webapp identity assign --resource-group $RESOURCE_GROUP --name $APP_NAME

groupid=$(az ad group create --display-name AzureSqlDbAccess --mail-nickname AzureSqlDbAccess --query objectId --output tsv)

siobjectid=$(az webapp identity show --resource-group $RESOURCE_GROUP --name $APP_NAME --query principalId --output tsv)
```

Query editor in cloud:
```sql
CREATE USER [AzureSqlDbAccess] FROM EXTERNAL PROVIDER;
ALTER ROLE db_datareader ADD MEMBER [AzureSqlDbAccess];
ALTER ROLE db_datawriter ADD MEMBER [AzureSqlDbAccess];
ALTER ROLE db_ddladmin ADD MEMBER [AzureSqlDbAccess];
GO
```
Need to get the above query automated.
