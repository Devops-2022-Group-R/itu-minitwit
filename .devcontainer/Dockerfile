# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.202.1/containers/go/.devcontainer/base.Dockerfile

ARG VARIANT="1.17-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

RUN apt-get update \
    && apt-get install -y sqlite3 libsqlite3-dev

# Install Terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/azure-get-started
RUN sudo apt-get install -y gnupg software-properties-common curl apt-transport-https \
    && curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add - \
    && curl https://baltocdn.com/helm/signing.asc | sudo apt-key add - \
    && echo "deb https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list \
    && sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
    && sudo apt-get update && sudo apt-get -y install terraform helm \
    # Install Azure CLI: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli-linux?pivots=apt
    && curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash \
    && az aks install-cli \
    && (cd /; GO111MODULE=on go get github.com/bitnami-labs/sealed-secrets/cmd/kubeseal@main)

RUN git config --global core.autocrlf true