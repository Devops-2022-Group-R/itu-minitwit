# Your turn now: Security!

## 1) Perform a Security Assessment 

The following general steps will guide you through a security assessment. Consider using them as steps in a report. The report will become a section in your final project report.

### A. Risk Identification

1. Identifiy assets (e.g. web application)
2. Identify threat sources (e.g. SQL injection)
3. Construct risk scenarios (e.g. Attacker performs SQL injection on web application to download sensitive user data)

### B. Risk Analysis

1. Determine likelihood
2. Determine impact
3. Use a Risk Matrix to prioritize risk of scenarios   
4. Discuss what are you going to do about each of the scenarios

### C. Pen-Test Your System

- Try to test for vulnerabilities in your project by using `wmap`, [`zaproxy`](https://www.zaproxy.org/getting-started/), or any of the tools in the [list of OWASP vulnerability scanning tools](https://owasp.org/www-community/Vulnerability_Scanning_Tools))
- Fix at least one vulnerability that you find; ideally one that is high in your prioritization cf. to your risk analysis 


*To think about*: can you find the traces of the pen test in the logs? Or of your colleagues pen-test?

## 2) White Hat Attack The Next Team

Try to help your fellow colleagues by pen-testing their system. Remember that the goal is to help not to hinder.  Send them a report of what you find. 

For a given group, their "fellow colleagues" are represented by the next group in the [repositories](https://github.com/itu-devops/lecture_notes/blob/master/repositories.py) file. Group R wraps back to 


----
### A. Risk Identification
#### Identify assets
- Data (backend)
- Monitoring (Grafana / Prometheus)
- Flag tool ()
- Azure credentials (Deployment service)
- Front end
- API 

#### Identify threat sources
- SQL Injection
- DDOS
- Broken access control
- Vulnerable and Outdated components

#### Construct risk scenarios
- Attacker performs DDOS making the service unavailable
- Attacker gains access to the monitoring layer, able to change dashboard, alerts and access business information regarding different endpoints  
- Attacker performs SQL Injection to copy, delete or ransom all data from the Database layer
- Attacker abuses security flaw of a vulnerable component and gains access to unwanted parts of the program
- Attacker gains access to Azure credentials, have complete control over the service, change subscribtion, database, resource group, redirect CI/CD pipeline
- Attack gains Access to CircleCi credntials, edit environment variables, stop the pipeline CI/CD

### B. Risk Analysis

#### Risk Matrix
![Risk Accessment Matrix](./RiskAssessmentMatrix.png)

#### Discuss what are you going to do about each of the scenarios
- SQL injection solved by using an ORM as middleware between database and input - gorm 
- DDOS Firewall, Maybe introduce a bandwidth cap on end points, to avoid whole system breaking down.
- Security flaw of vulnerable or outdated components. We use Static analysis tools like Snyk encorporated with the pipeline, to be aware about vulnerable dependencies and update the dependencies.
- Access to CircleCi, Access is granted via github user, on github set a requirement on the organisation repository to require 2FA.   
- Access to Azure crendetials, Require MFA for every user.  
- Access to monitoring, Create a team to limit permissions in grafana, add users to the team. 

### C. Pen-Test Your System

### Pen testing steps
- Zaproxy didn't work as intented with kubernetes [simplyzee](https://github.com/simplyzee/kube-owasp-zap)
- Run ZapProxy via the executable targeting Rhododevdron frontpage https://rhododevdron.swuwu.dk/, found a few obscure risks
    -  Missing header settings on the root endpoint, but not (CSP, Anti-clickjacking Header, X-Content-Type-Options Header, Incomplete or No Cache-control Header)
- Run WMAP in a docker container, targeting Rhododevdron frontpage https://rhododevdron.swuwu.dk/ see steps [Notes](./session09.md) in section Metasploit WMAP 

![Zaproxy results](./ZaproxyAlarms.png)

## 2) White Hat Attack The Next Team group A.
### Pen testing

- We tried SQL injection, but with no luck, seems like you got everything set up nicely in gorm
- You seem to have up to date packages, so finding new exploits is difficult
- We targeted your URL https://minitwit.thesvindler.net with [Zaproxy](https://www.zaproxy.org/download/), see results below

![image](https://user-images.githubusercontent.com/75098556/163984222-ec9a2556-9809-4823-a930-3f1f385dc03d.png)

**Overall niceness and good luck!**

