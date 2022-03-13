Running the Prometheus server (powershell, you might want to replace `${}` with `$()`):
```ps1
docker run -it --rm -p 9090:9090 -v ${pwd}/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

Running the Grafana server:
```ps1
docker run -it --rm -p 3000:3000 --name=grafana -v ${pwd}/monitoring/grafana/provisioning:/etc/grafana/provisioning grafana/grafana
```

Creating the data source in Grafana, set the URL to `http://host.docker.internal:9090`.
