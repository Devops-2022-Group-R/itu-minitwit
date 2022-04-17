# New host complete deployment run through
## Networking
*This section assumes CWD is in .infrastructure/kubernetes/networking*

Deploy the networking namespace
```
kubectl apply -f network-namespace.yaml
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

## Setup secrets
Deploy [sealed secrets](https://github.com/bitnami-labs/sealed-secrets) controller
```
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.17.4/controller.yaml 
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

## Backend and database
*This section assumes CWD is in .infrastructure/kubernetes/backend*

Create the backend namespace
```
kubectl apply -f backend-namespace.yaml
```

### Database
Create database secrets
```bash
kubectl apply -f database-secret.yaml
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
kubectl apply -f monitoring-storage.yaml # Create the nfs volume subfolder "grafana" before
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
kubectl apply -f elasticsearch-storage.yml # Create the nfs volume subfolder "logging" before
kubectl apply -f elasticsearch.yml
```

### Kibana
Setup
```bash
kubectl apply -f kibana.yml
```

### Fluentd
*Note before doing the below, i had to change the parsing method when retrieving logs from /var/log as azure and pure kubernetes apparently have different formats*

Setup.
```
kubectl apply -f fluentd-config.yml
kubectl apply -f fluentd.yml
```