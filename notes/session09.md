## Logging in Kubernetes

### Storage
Retrieve storage info and create secret
```bash
STORAGE_ACCOUNT_NAME=regnburclusterstorage
SHARE_NAME=cluster-persistent-storage
RG=itu-minitwit-rg

STORAGE_KEY=$(az storage account keys list --resource-group $RG --account-name $STORAGE_ACCOUNT_NAME --query "[0].value" -o tsv)

# Create secret
kubectl create secret generic azure-storage-secret --from-literal=azurestorageaccountname=$STORAGE_ACCOUNT_NAME --from-literal=azurestorageaccountkey=$STORAGE_KEY
```

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
kubectl apply -f logging/elasticsearch-storage.yml
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
Set the password you got earlier in the designated space in fluentd.yml, then
```
kubectl apply -f logging/fluentd-config.yml
kubectl apply -f logging/fluentd.yml
```

- https://medium.com/avmconsulting-blog/how-to-deploy-an-efk-stack-to-kubernetes-ebc1b539d063
- https://blog.kubernauts.io/simple-logging-with-eck-and-fluentd-13824ad65aaf
- https://www.digitalocean.com/community/tutorials/how-to-set-up-an-elasticsearch-fluentd-and-kibana-efk-logging-stack-on-kubernetes


# Security

## The plan
- [x] Perform a Security Assessment
    - [x] Risk Identification 
    - [x] Risk Analysis
    - [x] Pen-Test Your System
    - [x] ZapProxy
    - [x] wmap - https://www.metasploit.com/
    - [x] other tools in the [list of OWASP vulnerability scanning tools](https://owasp.org/www-community/Vulnerability_Scanning_Tools))
    - [x] Fix at least one vulnerability. (e.g. monitoring access control)

- [x] White Hat Attack The Next Team
    - [x] Zaproxy
    - [x] Wmap - Had issues when targeting [GroupA](https://minitwit.thesvindler.net) - OpenSSL tlsv1 alert internal error

## Notes

### Pen testing steps
- Zaproxy didn't produce results as expected with kubernetes [simplyzee](https://github.com/simplyzee/kube-owasp-zap)
- Run ZapProxy via the executable targeting Rhododevdron frontpage https://rhododevdron.swuwu.dk/public, found a few obscure risks
    -  Missing header settings (Anti-clickjacking Header, X-Content-Type-Options Header, Incomplete or No Cache-control Header)
- Run WMAP in a docker container, targeting Rhododevdron frontpage https://rhododevdron.swuwu.dk/public see steps below 

#### Zaproxy executable steps
1. Install [Zaproxy](https://www.zaproxy.org/download/)
2. specify url to target Rhododevdron frontpage https://rhododevdron.swuwu.dk/
    - [x] use traditional spider
    - [x] use ajax spider with FireFox Headless
3. Read alarms in the Zaproxy GUI

#### Zaproxy kubernetes steps

If issue with connecting to the cluster, try setting the subscription Id found in Azure

```sh
az account set --subscription cluser-subscrition-id
```

from .infrastructure/kubernetes/security
```sh
kubectl apply -f zap.yaml
```

```sh
helm repo add simplyzee https://charts.simplyzee.dev
```

Setup environment variables

```sh
export URL_TO_SCAN="url"
```

Run Zap scanner
```sh
helm install "vuln-scan-$(date '+%Y-%m-%d-%H-%M-%S')-job" simplyzee/kube-owasp-zap \
    --namespace owasp-zap \
    --set zapcli.debug.enabled=true \
    --set zapcli.spider.enabled=false \
    --set zapcli.recursive.enabled=false \
    --set zapcli.targetHost=$URL_TO_SCAN
```

Show available logs sorted by most recent date

```sh
kubectl get jobs --namespace owasp-zap | grep -v "COMPLETIONS" | sort
```
Get corresponding pod
```sh
kubectl get pods --namespace owasp-zap
```

Select log
```sh
kubectl logs <podname> --namespace owasp-zap
```

#### Metasploit WMAP

Resources used
- [Setup Metasploit database in Kali Docker Container](https://gist.github.com/pich4ya/e7be40000c4fe7e487460dbebf1832fb)
- [Metasploit WMAP in linux](https://linuxhint.com/metasploit_vurnerability_scanner_linux/)
- [Fix missing 'systemd' in docker container](https://mefmobile.org/fix-systemctl-command-not-found/#:~:text=based%20operating%20systems.-,What%20is%20causing%20the%20%E2%80%9CSystemctl%3A%20command%20not%20found%E2%80%9D%20error,SysV%20init%20instead%20of%20systemd%20.)

Installing Docker image
```ps1
docker pull kalilinux/kali-rolling 
```

```ps1
docker run -it kalilinux/kali-rolling bash
```

Access running container
```ps1
docker exec -it <docker_container_name> bash
```

Commands inside image to setup Metasploit
```sh
apt update && 
apt -y upgrade &&
apt-get install postgresql &&
apt install metasploit-framework &&
dpkg -l | grep systemd &&
apt-get update &&
apt-get install systemd &&
msfdb init && 
service postgresql start
```

Open metasploit console
```ps1
msfconsole
```

To check postgres connection
```msfconsole
db_status
```

Load WMAP plugin
```msfconsole
load wmap
```

Add sites to scan
```msfconsole
wmap_sites -a https://rhododevdron.swuwu.dk/
```

List of sites
```msfconsole
wmap_sites -l
```

Add targets (sites or sub pages, based on wmap_sites) 
```msfconsole
wmap_targets -d [wmap_sites id] 
```
```msfconsole
wmap_targets -t [url]
```

Show enabled modules
```msfconsole
wmap_run -t
```

Run scanner on all targets (warning takes a long time)
```msfconsole
wmap_run -e
```

Show vulnerabilites 
```msfconsole
wmap_vulns -l
```

##### Fix issue with postgres connection [Metasploit - Authentication failed for user "msf"](https://github.com/rapid7/metasploit-framework/issues/9696)
```sh
msfdb reinit
```

Ensure postgres server is running
```sh 
pg_lsclusters
```

Start postgres
```sh
service postgresql start
```

### Monitoring vulnerability
See scenario discussion [Security](./session09_Security.md)

### Steps to test and fix security vulnerabilites locally
1. Change root endpoint in env variable(VITE_API_URL) in the front end dockerfile to localhost 
2. Build new docker image of the front end   
3. Run the front end image map port 8080:80
4. Start itu-minitwit locally
5. Target localhost with WMAP, Zaproxy etc.
6. Add fix,
7. Verify by iterating from step 2. or 4. 
