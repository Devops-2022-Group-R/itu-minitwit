
## Logging in Kubernetes
### Elasticsearch
See:
- https://www.elastic.co/blog/how-to-run-elastic-cloud-on-kubernetes-from-azure-kubernetes-service
- https://www.elastic.co/guide/en/cloud-on-k8s/current/index.html

Install ECK operator:
```bash
kubectl create -f https://download.elastic.co/downloads/eck/2.1.0/crds.yaml

kubectl apply -f https://download.elastic.co/downloads/eck/2.1.0/operator.yaml
```

Start pod:
```bash
kubectl apply -f logging/elasticsearch.yml
```

Check status:
```bash
kubectl get pods --selector='elasticsearch.k8s.elastic.co/cluster-name=itu-minitwit-elasticsearch' -n itu-minitwit-logging-ns
# or
kubectl get elasticsearch -n itu-minitwit-logging-ns
```

Get password:
```bash
PASSWORD=$(kubectl get secret itu-minitwit-elasticsearch-es-elastic-user -n itu-minitwit-logging-ns -o go-template='{{.data.elastic | base64decode}}')
```

In a separate terminal, set up port forwarding from local machine to ES:
```bash
kubectl port-forward service/itu-minitwit-elasticsearch-es-http 9200 -n itu-minitwit-logging-ns
```
See if above port forwarding works:
```bash
curl -u "elastic:$PASSWORD" -k "https://localhost:9200"
```

## Kibana
Start pod:
```bash
kubectl apply -f logging/kibana.yml
```

See status:
```bash
kubectl get pod --selector='kibana.k8s.elastic.co/name=itu-minitwit-kibana' -n itu-minitwit-logging-ns
# or
kubectl get kibana -n itu-minitwit-logging-ns
```

In a separate terminal, set up port forwarding from the local machine to Kibana:
```bash
kubectl port-forward service/itu-minitwit-kibana-kb-http 5601 -n itu-minitwit-logging-ns
```

See if it's responding:
```http
curl -k -i https://localhost:5601

# should respond with

HTTP/1.1 302 Found
location: /login?next=%2F
...
```

## Fluentd
- https://medium.com/avmconsulting-blog/how-to-deploy-an-efk-stack-to-kubernetes-ebc1b539d063
- https://blog.kubernauts.io/simple-logging-with-eck-and-fluentd-13824ad65aaf
- https://www.digitalocean.com/community/tutorials/how-to-set-up-an-elasticsearch-fluentd-and-kibana-efk-logging-stack-on-kubernetes
