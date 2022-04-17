## The plan
- [x] Perform a Security Assessment
    - [x] Risk Identification 
    - [x] Risk Analysis
    - [ ] Pen-Test Your System
    - [x] ZapProxy
    - [ ] wmap - https://www.metasploit.com/
    - [ ] other tools in the [list of OWASP vulnerability scanning tools](https://owasp.org/www-community/Vulnerability_Scanning_Tools))
    - [ ] Fix at least one vulnerability. (e.g. monitoring access control)

- [x] White Hat Attack The Next Team

## Notes

### Pen testing steps
- Zaproxy didn't work as intented with kubernetes [simplyzee](https://github.com/simplyzee/kube-owasp-zap))
- Run ZapProxy via the executable targeting Rhododevdron frontpage https://rhododevdron.swuwu.dk/public, found a few obscure risks
    -  Missing header settings (Anti-clickjacking Header, X-Content-Type-Options Header, Incomplete or No Cache-control Header)
- Run WMAP in a docker container, targeting Rhododevdron frontpage https://rhododevdron.swuwu.dk/public see steps below 

### Zaproxy executable steps
1. Install [Zaproxy](https://www.zaproxy.org/download/)
2. specify url to target Rhododevdron frontpage https://rhododevdron.swuwu.dk/public
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

### Monitoring vulnerability
See scenario discussion [Security](./session09_Security.md)


### Metasploit WMAP

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

Fix issue with postgres connection [Metasploit - Authentication failed for user "msf"](https://github.com/rapid7/metasploit-framework/issues/9696)
```sh
msfdb reinit
```

Ensure postgres server is running
```sh 
pg_lsclusters
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
wmap_sites -a https://rhododevdron.swuwu.dk/public
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

Run scanner on all targets (warning takes a long time) - Seems to halt at brute_dirs module
- [ ] Select or disable specific modules
```msfconsole
wmap_run -e
```

Show vulnerabilites 
```msfconsole
wmap_vulns -l
```