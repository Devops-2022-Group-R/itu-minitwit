# New host complete deployment run through
## Setup storage account for terraform
**Note: Only necessary if using Terraform to manage infrastructure**

```bash
az group create -n core-rg -l northeurope
 
# Create Storage Account
az storage account create -n minitwitterraformstate2 -g core-rg -l northeurope --sku Standard_LRS

# List keys with
az storage account keys list --account-name minitwitterraformstate2

# Create Storage Account Container
az storage container create -n terraformstate --account-name minitwitterraformstate2 --account-key <key-from-created-account> 
```

## Setup secrets
Deploy [sealed secrets](https://github.com/bitnami-labs/sealed-secrets) controller
```
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.18.0/controller.yaml  
```

### To create new secrets, using mssql server as an example
*This example assumes your CWD is in .infrastructure/kubernetes/backend*

*If no secret creation is necessary, just skip to the next step*

Install the kubeseal CLI from https://github.com/bitnami-labs/sealed-secrets/releases. Then (only on new clusters) fetch the certificate from the cluster, and output it to a file
```bash
kubeseal --fetch-cert > clustercert.pem
```
Create a new secret yaml file. DO NOT COMMIT THIS FILE, it's not secure
```bash
kubectl create secret generic itu-minitwit-mssql \
    --namespace itu-minitwit-backend-ns \
    --dry-run=client \
    --from-literal=SA_PASSWORD="mysupersecretpassword" \
    -o yaml > database-secret-unsealed.yaml
```
Create the sealed secret file
```
kubeseal --cert ../clustercert.pem -o yaml < database-secret-unsealed.yaml > database-secret.yaml
```

## Storage
**Note: Only necessary if the host is not on Azure, it is also required to change disk to file in storage-class**

Create Azure Service Principal
```bash
az login
az account set --subscription <subscription-id>
az group list --query "[?name=='itu-minitwit-rg'].id" -o tsv # Get scope
az ad sp create-for-rbac -n "itu-minitwit-cluster-storage-sp" --role "Contributor" --scopes <scope-from-above-command> # Note down the values from this command
```
Create azure.json of the format:
```
{
    "cloud":"AzurePublicCloud",
    "tenantId": "0000000-0000-0000-0000-000000000000", 
    "aadClientId": "0000000-0000-0000-0000-000000000000",
    "aadClientSecret": "0000000-0000-0000-0000-000000000000",
    "subscriptionId": "0000000-0000-0000-0000-000000000000",
    "resourceGroup": "itu-minitwit-rg",
    "location": "northeurope",
    "cloudProviderBackoff": false,
    "useManagedIdentityExtension": false,
    "useInstanceMetadata": true
}
```
`tenantId,subscriptionId` can be retrieved with `az account show`

`aadClientId,aadClientSecret` is from the service principal create before

`resourceGroup,location` can be retrieved with `az group`

Create secret
```bash
cat azure.json | \
kubectl create secret generic azure-cloud-provider \
    --namespace kube-system \
    --dry-run=client \
    --from-file=cloud-config=/dev/stdin \
    -o yaml > azure-cloud-provider-unsealed.yaml

kubeseal --cert ../clustercert.pem -o yaml < azure-cloud-provider-unsealed.yaml > azure-cloud-provider.yaml

kubectl apply -f azure-cloud-provider.yaml
```

Install the driver
```bash
helm repo add azurefile-csi-driver https://raw.githubusercontent.com/kubernetes-sigs/azurefile-csi-driver/master/charts
helm repo update

helm install azurefile-csi-driver azurefile-csi-driver/azurefile-csi-driver \
    --namespace kube-system \
    --set controller.cloudConfigSecretName="azure-cloud-provider" \
    --set controller.cloudConfigSecretNamespace="kube-system" \
    --set node.cloudConfigSecretName="azure-cloud-provider" \
    --set node.cloudConfigSecretNamespace="kube-system"
```

### Create storage-class
```bash
kubectl apply -f storage-class.yaml
```

## Networking
*This section assumes CWD is in .infrastructure/kubernetes/networking*

Deploy the networking namespace
```
kubectl apply -f network-namespace.yaml
```

### Ingress
```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx-ingress ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace -f ingress.yaml
```

### Setup dns zones
#### Azure
```bash
az ad sp create-for-rbac -n ExternalDnsServicePrincipal

spId="9ec83678-21e4-4efd-8496-0e642c2b65b0" # Replace with real id

rg="itu-minitwit-rg"
# Get resource group id
rgId=$(az group show --name $rg --query "id" -o tsv)
# Get dns zone id
dnsId=$(az network dns zone show --name rhododevdron.dk -g $rg --query "id" -o tsv)
dnsId2=$(az network dns zone show --name swuwu.dk -g $rg --query "id" -o tsv)

az role assignment create --role "Reader" --assignee $spId --scope $rgId
az role assignment create --role "DNS Zone Contributor" --assignee $spId --scope $dnsId
az role assignment create --role "DNS Zone Contributor" --assignee $spId --scope $dnsId2

# Get cluster identity id
clusterIdentity=$(az aks show -g $rg -n itu-minitwit-cluster --query "identity.principalId" -o tsv)

# Make the cluster identity a network contributor on the other resource group
az role assignment create --assignee $clusterIdentity --role "Network Contributor" --scope $rgId
```

Create a file, called `external-dns-secret.json` that looks like this, with the data from the service principal:
```json
{
  "tenantId": "01234abc-de56-ff78-abc1-234567890def",
  "subscriptionId": "01234abc-de56-ff78-abc1-234567890def",
  "resourceGroup": "MyDnsResourceGroup",
  "aadClientId": "01234abc-de56-ff78-abc1-234567890def",
  "aadClientSecret": "uKiuXeiwui4jo9quae9o"
}
```
Then run
```
kubectl create secret generic external-dns-config-file --from-file=/path/to/external-dns-secret.json
kubectl apply -f external-dns-deployment.yaml
```

#### Digitalocean
```
TOKEN=""
kubectl create secret generic digitalocean-token-secret --from-literal=token=$TOKEN
kubectl apply -f external-dns-deployment.yaml
```

### Setup and deploy cluster issuer
To provide the sites with SSL certificates we need to create some cluster issuers. First add the helm repo and install the cert-manager helm chart.
```
# Label the cert-manager namespace to disable resource validation
kubectl label namespace ingress-nginx cert-manager.io/disable-validation=true

# Add the Jetstack Helm repository
helm repo add jetstack https://charts.jetstack.io

# Update your local Helm chart repository cache
helm repo update

helm install cert-manager jetstack/cert-manager  \
   --namespace ingress-nginx \
   -f cert-manager-values.yaml
```
Deploy certificate issuer
```
kubectl apply -f cluster-issuer-staging.yaml
kubectl apply -f cluster-issuer-prod.yaml
```

## Backend and database
*This section assumes CWD is in .infrastructure/kubernetes/backend*

Create the backend namespace
```
kubectl apply -f backend-namespace.yaml
```

### Database
Create database secrets
```bash
kubectl apply -f database-secret.yaml # Created in the example above
```
Create database volume and claims (requires an nfs volume)
```bash
kubectl apply -f database-storage.yaml # Create the nfs volume subfolder "database" before this command
```
Create database:
```bash
kubectl apply -f database.yaml
```

### Backend
#### Create a secret
Same thing as with the database, now we create a connection string for the backend to the database, again DO NOT COMMIT THIS FILE
```bash
echo -n "server=itu-minitwit-database-deployment;database=master;user id=sa;password=mysupersecretpassword" | \
kubectl create secret generic itu-minitwit-backend-secrets \
    --namespace itu-minitwit-backend-ns \
    --dry-run=client \
    --from-file=CONNECTION_STRING=/dev/stdin \
    -o yaml > backend-secrets-unsealed.yaml

kubeseal --cert ../clustercert.pem -o yaml < backend-secrets-unsealed.yaml > backend-secrets.yaml
```
#### Deploy
```bash
kubectl apply -f backend-secrets.yaml
kubectl apply -f backend-deployment.yaml
```

## Frontend
*This section assumes CWD is in the .infrastructure/kubernetes of the frontend repo*

Deploy the frontend and it's namespace
```
kubectl apply -f frontend-namespace.yaml
kubectl apply -f frontend-deployment.yaml
```

## Monitoring
*This section assumes CWD is in .infrastructure/kubernetes/monitoring*

Deploy the monitoring and it's namespace
```bash
kubectl apply -f monitoring-namespace.yaml
kubectl apply -f monitoring-storage.yaml
kubectl apply -f monitoring-grafana-deployment.yaml
kubectl apply -f monitoring-prometheus-deployment.yaml
```
*Note: The above will respond to requests, but prometheus won't collect any data unless the DNS rules have been set up*

## Logging
*This section assumes CWD is in .infrastructure/kubernetes/logging*
Create the logging namespace
```
kubectl apply -f logging-namespace.yaml
```

### Elastic search
Setup the operator
```bash
kubectl create -f https://download.elastic.co/downloads/eck/2.1.0/crds.yaml
kubectl apply -f https://download.elastic.co/downloads/eck/2.1.0/operator.yaml
```

Create the pods, persistent volume and namespace
```bash
kubectl apply -f elasticsearch.yml
```

### Kibana
Setup
```bash
kubectl apply -f kibana.yml
```

Get log in password
```bash
kubectl get secret itu-minitwit-elasticsearch-es-elastic-user -n itu-minitwit-logging-ns -o go-template='{{.data.elastic | base64decode}}'
```

### Fluentd
*Note before doing the below, i had to change the parsing method when retrieving logs from /var/log as azure and pure kubernetes apparently have different formats*

Setup.
```
kubectl apply -f fluentd-config.yml
kubectl apply -f fluentd.yml
```