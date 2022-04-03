## The plan
- [x] Perform a Security Assessment
- - [x] Risk Identification 
- - [x] Risk Analysis
- - [ ] Pen-Test Your System
- - [x] ZapProxy
- - [ ] wmap - https://www.metasploit.com/
- - [ ] other tools in the [list of OWASP vulnerability scanning tools](https://owasp.org/www-community/Vulnerability_Scanning_Tools))
- - [ ] Fix at least one vulnerability. (e.g. monitoring access control)

- [x] White Hat Attack The Next Team

## Notes

### pen testing steps
- ZapProxy didn't work as intented with kubernetes [simplyzee](https://github.com/simplyzee/kube-owasp-zap))
- run ZapProxy via the executable targeting our  [Rhododevdron frontpage]("https://rhododevdron.swuwu.dk/public"), found a few obscure risks
- - Missing header settings (Anti-clickjacking Header, X-Content-Type-Options Header, Incomplete or No Cache-control Header)

### Zaproxy executable steps
Install [ZapProxy](https://www.zaproxy.org/download/)
specify url to target [Rhododevdron frontpage]("https://rhododevdron.swuwu.dk/public")
- [x] use traditional spider
- [x] use ajax spider with FireFox Headless
Read alarms in the GUI

#### Zaproxy kubernetes steps

from .infrastructure/kubernetes/security
```sh
kubectl apply -f zap.yml
```

```sh
helm repo add simplyzee https://charts.simplyzee.dev
```

```sh
helm install "vuln-scan-$(date '+%Y-%m-%d-%H-%M-%S')-job" simplyzee/kube-owasp-zap \
    --namespace owasp-zap \
    --set zapcli.debug.enabled=true \
    --set zapcli.spider.enabled=false \
    --set zapcli.recursive.enabled=false \
    --set zapcli.scanTypes=$SCAN_TYPE \
    --set zapcli.targetHost=$URL_TO_SCAN
```

```sh
kubectl get jobs --namespace owasp-zap | grep -v "COMPLETIONS" | sort`
```

```sh
kubctl logs vuln-scan-2022-04-03-13-45-07-job-kube-owasp-zap-qmfbl --namespace owasp-zap
```

### monitoring vulnerability
See scenario discussion [security](./session09_Security.md)