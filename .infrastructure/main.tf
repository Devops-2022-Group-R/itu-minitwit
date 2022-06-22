# Configure the Azure provider
terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }

  # Terraform state is managed in Digitalocean - create a space to use
  backend "s3" {
    endpoint                    = "fra1.digitaloceanspaces.com"
    key                         = "terraform.tfstate"
    bucket                      = "minitwit-state"
    region                      = "us-east-1"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
  }

  required_version = ">= 1.1.0"
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_kubernetes_cluster" "cluster" {
  name = "${var.prefix}-cluster"
  region = "fra1"
  version = "1.22.8-do.1"

  node_pool {
    name = "worker-pool"
    size = "s-4vcpu-8gb"
    node_count = 1
  }
}
