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
    storage_account_name = "minitwitterraformstate2"
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

resource "azurerm_dns_zone" "cluster_dns_zone_swuwu" {
  name                = "swuwu.dk"
  resource_group_name = azurerm_resource_group.rg.name
}