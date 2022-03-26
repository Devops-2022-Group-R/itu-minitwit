We deviate from the official schedule this week. We will transfer our services to a Kubernetes cluster to make adding logging later easier.

## The Plan
- [x] Fix Terraform
- [x] Create Kubernetes cluster in Azure
- [x] Create automated Kubernetes deployment
- [x] Create Kubernetes pod definitions
   - [x] Main API server
   - [x] Monitoring setup
   - [x] Frontend
- [x] Remove services moved to Kubernetes from Terraform
- [x] Networking
   - [x] Create public ip
   - [x] Create ingress controller with public ip as the load balancer ip
   - [x] Create dns zone
   - [x] Create service principal and permissions for external-dns
   - [x] Deploy external-dns
   - [x] Alter service to use ingress controller
   - [x] Deploy cert-manager
   - [x] Create certificates
   - [x] hope for the best 
- [ ] Set up persistent storage

## Terms that everyone should know
- Cluster
- Namespace
- Pod
- "Expose"
- Load balancer
- Ingres controller
- API gateway

We will likely not make use of some of the above (especially the bottom ones), but it's great to know about.

## Notes
Connect to the cluster:
```
az aks get-credentials --resource-group itu-minitwit-rg --name itu-minitwit-cluster
```

### Setup Networking
```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx-ingress ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace -f ingress.yaml
```

```bash
az ad sp create-for-rbac -n ExternalDnsServicePrincipal

# Get resource group id
rgId=$(az group show --name $rg --query "id" -o tsv)
# Get dns zone id
dnsId=$(az network dns zone show --name rhododevdron.dk -g $rg --query "id" -o tsv)

az role assignment create --role "Reader" --assignee <service principal appId> --scope $rgId
az role assignment create --role "DNS Zone Contributor" --assignee <service principal appId> --scope $dnsId

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

### Setup SSL
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

### Setup sealed secrets
```
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.17.3/controller.yaml
```

### Deployment
From .infrastructure/kubernetes/backend run the following to deploy the backend
```
kubectl apply -f backend-namespace.yaml
kubectl apply -f backend-secrets.yaml
kubectl apply -f backend-deployment.yaml
```

From .infrastructure/kubernetes/monitoring run the following to deploy the monitoring
```
kubectl apply -f monitoring-namespace.yaml
kubectl apply -f monitoring-grafana-deployment.yaml
kubectl apply -f monitoring-prometheus-deployment.yaml
```

From .infrastructure/kubernetes/frontend run the following to deploy the monitoring
```
kubectl apply -f frontend-namespace.yaml
kubectl apply -f frontend-deployment.yaml
```