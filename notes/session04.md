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

### Database
Changes necessary for the azure deployment in terms of the database:
- When opening a database connection we run the migrations and check if the simulator user exists, if it doesn't create it
- Switch to [microsoft sql server driver](https://github.com/go-gorm/sqlserver)
- To use azure managed identity it was necessary to use azuread package in the microsoft sql server driver, [documentation here](https://github.com/denisenkom/go-mssqldb#azure-active-directory-authentication)

### Apply terraform with circle ci
Create a service principal in Azure
```
az ad sp create-for-rbac --name "terraform-contributor" --role Contributor --scopes /subscriptions/<subscription-id>
```

Create a storage account
```sh
# Create Resource Group
az group create -n <resource-group-name> -l northeurope
 
# Create Storage Account
az storage account create -n <account-name> -g <resource-group-name> -l northeurope --sku Standard_LRS
 
# Create Storage Account Container
az storage container create -n <container-name> --account-name <account-name> --account-key <key-from-created-account> 
```

### Hostname binding
The terraform script creates a hostname binding to a hostname that Azure does not have control over. To gain this control access the [portal](portal.azure.com) find the App Service and click "Custom domains" the hostname that you have added will have a red x next to it. Click it and you will be guided to add some records to the domain host of your choosing. In our case we had to add 

- CNAME record
  - Key: `api.rhododevdron.swuwu.dk `
  - Value: `itu-minitwit-backend-as.azurewebsites.net`
- TXT record
  - Key: `swuwu.dk`
  - Value: `MS=<key shown in azure>`
- TXT record
  - Key: `asuid.api.rhododevdron.swuwu.dk`
  - Value: `<key shown in azure>`

This has to also be done for the frontend. Completing these steps allows Azure to show content on that domain as well as create an SSL certificate for the domain. It's a one time setup which remains even in case of a resource teardown.

### Create deployment webhook
First get the URL from [here](https://github.com/sajayantony/appservicedemo), the url is:
```
https://<publishingusername>:<publishingpwd>@<publishurl>/docker/hook
```
Secondly retrieve the publish profile by going to your app service, pressing "Deployment Center" -> "Manage publish profile" -> "Download publish profile"

The downloaded XML contains alot of info but the necessary keys are:
- `publishUrl` which maps to `publishurl` in the url above
- `userName` which maps to `publishingusername`
- `userPWD` which maps to `publishingpwd`

Place the keys in the url and create a webhook on docker hub with it. This process should be repeated for every app service.