#!/bin/sh

# Update the cluster wide storage class
kubectl apply -f ./storage/storage-class.yaml

# Deploy nginx-ingress and cert-manager
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm upgrade nginx-ingress ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace -f ./networking/ingress.yaml
helm upgrade cert-manager jetstack/cert-manager --namespace ingress-nginx -f ./networking/cert-manager-values.yaml

kubectl apply -f ./networking/network-namespace.yaml
kubectl apply -f ./networking/external-dns-deployment.yaml

kubectl apply -f ./networking/cluster-issuer-staging.yaml
kubectl apply -f ./networking/cluster-issuer-prod.yaml


kubectl apply -f ./backend/backend-namespace.yaml

kubectl apply -f ./backend/database-secret.yaml
kubectl apply -f ./backend/database-storage.yaml
kubectl apply -f ./backend/database.yaml

kubectl apply -f ./backend/backend-secrets.yaml
kubectl apply -f ./backend/backend-deployment.yaml


kubectl apply -f ./monitoring/monitoring-namespace.yaml
kubectl apply -f ./monitoring/monitoring-storage.yaml
kubectl apply -f ./monitoring/monitoring-grafana-deployment.yaml
kubectl apply -f ./monitoring/monitoring-prometheus-deployment.yaml


kubectl apply -f ./logging/logging-namespace.yml
kubectl apply -f ./logging/elasticsearch.yml
kubectl apply -f ./logging/kibana.yml

kubectl delete -f ./logging/fluentd-config.yml # The change detection isn't consistent so we delete it before applying it here
kubectl apply -f ./logging/fluentd-config.yml

kubectl apply -f ./logging/fluentd.yml