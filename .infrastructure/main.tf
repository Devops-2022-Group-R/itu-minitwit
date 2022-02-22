# Configure the Azure provider
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 2.65"
    }
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
    size = "B2"
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
