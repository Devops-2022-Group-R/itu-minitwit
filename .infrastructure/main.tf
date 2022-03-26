# Configure the Azure provider
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 2.65"
    }
  }

  # Terraform state is managed in Azure - create a storage container to use:
  #   az group create -n core-rg -l northeurope
  #   az storage account create -n <account-name> -g core-rg -l northeurope --sku Standard_LRS
  #   az storage container create -n <container-name> --account-name <account-name> --account-key <key-from-created-account> 
  backend "azurerm" {
    resource_group_name  = "core-rg"
    storage_account_name = "minitwitterraformstate"
    container_name       = "terraformstate"
    key                  = "terraform.tfstate"
  }

  required_version = ">= 1.1.0"
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "${var.prefix}-rg"
  location = "northeurope"
}

resource "azurerm_app_service_plan" "asp" {
  name                = "${var.prefix}-asp"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Basic"
    size = "B1"
  }
}

resource "azurerm_app_service" "backend_as" {
  name                = "${var.prefix}-backend-as"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  app_service_plan_id = azurerm_app_service_plan.asp.id
  https_only          = true

  site_config {
    app_command_line = ""
    linux_fx_version = "DOCKER|kongborup/itu-minitwit-server:latest"
  }

  app_settings = {
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
    "DOCKER_REGISTRY_SERVER_URL"          = "https://registry.hub.docker.com"
    "DOCKER_ENABLE_CI"                    = "true"
  }

  connection_string {
    name  = "CONNECTION_STRING"
    type  = "SQLServer"
    value = "server=${var.database_server_name}.database.windows.net;database=${var.database_db_name};fedauth=ActiveDirectoryMSI"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_app_service_custom_hostname_binding" "backend_custom_domain" {
  hostname            = "api.rhododevdron.swuwu.dk"
  app_service_name    = azurerm_app_service.backend_as.name
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_app_service_managed_certificate" "backend_managed_certificate" {
  custom_hostname_binding_id = azurerm_app_service_custom_hostname_binding.backend_custom_domain.id
}

resource "azurerm_app_service_certificate_binding" "backend_certificate_binding" {
  hostname_binding_id = azurerm_app_service_custom_hostname_binding.backend_custom_domain.id
  certificate_id      = azurerm_app_service_managed_certificate.backend_managed_certificate.id
  ssl_state           = "SniEnabled"
}

resource "azurerm_mssql_server" "database_mssql_server" {
  name                = var.database_server_name
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  version             = "12.0"

  administrator_login          = var.database_admin_username
  administrator_login_password = var.database_admin_password

  minimum_tls_version = "1.2"

  azuread_administrator {
    # Security group that should be made in Azure to provide access
    login_username              = "Admins"
    object_id                   = "fa37a2f2-6d36-45e6-8b20-fa037e932ac6"
    azuread_authentication_only = false
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_mssql_firewall_rule" "database_firewall_rule" {
  name             = "${var.prefix}-allow-azure-ips"
  server_id        = azurerm_mssql_server.database_mssql_server.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "0.0.0.0"
}

resource "azurerm_mssql_database" "database_mssql_database" {
  name        = var.database_db_name
  server_id   = azurerm_mssql_server.database_mssql_server.id
  max_size_gb = 2
  sku_name    = "Basic"
}

resource "azurerm_app_service" "monitoring_as" {
  name                = "${var.prefix}-monitoring-as"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  app_service_plan_id = azurerm_app_service_plan.asp.id
  client_cert_enabled = false
  site_config {
    ftps_state       = "Disabled"
    app_command_line = ""
    linux_fx_version = "COMPOSE|${filebase64("../monitoring/docker-compose.yml")}"
  }

  app_settings = {
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
  }
}

resource "azurerm_kubernetes_cluster" "cluster" {
  name                = "${var.prefix}-cluster"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = "regnbur"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "standard_d2as_v4"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_dns_zone" "cluster_dns_zone" {
  name                = "rhododevdron.dk"
  resource_group_name = azurerm_resource_group.rg.name
}