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

### Zaproxy executable steps
1. Install [Zaproxy](https://www.zaproxy.org/download/)
2. specify url to target Rhododevdron frontpage https://rhododevdron.swuwu.dk/public
    - [x] use traditional spider
    - [x] use ajax spider with FireFox Headless
3. Read alarms in the Zaproxy GUI

#### Zaproxy kubernetes steps

from .infrastructure/kubernetes/security
```sh
kubectl apply -f zap.yml
```

```sh
helm repo add simplyzee https://charts.simplyzee.dev
```

Setup environment variables

```sh
export URL_TO_SCAN="url"
```

```sh
export SCAN_TYPE="all"
```


Run Zap scanner
```sh
helm install "vuln-scan-$(date '+%Y-%m-%d-%H-%M-%S')-job" simplyzee/kube-owasp-zap \
    --namespace owasp-zap \
    --set zapcli.debug.enabled=true \
    --set zapcli.spider.enabled=false \
    --set zapcli.recursive.enabled=false \
    --set zapcli.scanTypes=$SCAN_TYPE \
    --set zapcli.targetHost=$URL_TO_SCAN
```

Show available logs sorted by most recent date

```sh
kubectl get jobs --namespace owasp-zap | grep -v "COMPLETIONS" | sort`
```

Select log
```sh
kubctl logs vuln-scan-2022-01-01-01-01-01-job-kube-owasp-zap-qmfbl --namespace owasp-zap
```

### Monitoring vulnerability
See scenario discussion [security](./session09_Security.md)
