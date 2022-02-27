# ITU MiniTwit
*Updated by: Group R*

To build the MiniTwit Docker image:
```sh
docker build -t minitwit .
```

To run the MiniTwit Docker container:
```sh
docker run -p 8080:8080 -it minitwit
```

## Deploy
The project is automatically deployed when a new `:latest` tag is pushed to [Docker Hub](https://hub.docker.com/repository/docker/kongborup/itu-minitwit-server). This happens automatically once a commit is made to the `master` branch.

This assumes that Azure has been set up correctly.

### Azure setup
#### Prerequisites
- [Install Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/azure-get-started)
- [Install the Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) and follow the initial steps of this [Terraform guide](https://learn.hashicorp.com/tutorials/terraform/azure-build)

#### Steps
1. Make sure the working directory is `.infrastructure` (`cd .infrastructure` from the project root).
2. Ensure that you have the `.infrastructure/.tfvars` file and replace the `<add_value>` fields with your secrets and config values.
3. Run `terraform apply -var-file .tfvars`
