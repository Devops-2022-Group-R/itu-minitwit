## Running the services
Running the Prometheus server (powershell, you might want to replace `${}` with `$()`):
```ps1
docker run -it --rm -p 9090:9090 -v ${pwd}/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

Running the Grafana server:
```ps1
docker run -it --rm -p 3000:3000 --name=grafana -v ${pwd}/monitoring/grafana/provisioning:/etc/grafana/provisioning grafana/grafana
```

## Making changes
To update Prometheus rules, edit `monitoring/prometheus.yml`.

To update Grafana settings, use Grafana's provisioning ([see docs](https://grafana.com/docs/grafana/latest/administration/provisioning/)). In short, Prometheus is added as a data source in `monitoring/grafana/provisioning/datasources`, and the general dashboard settings are defined in `monitoring/grafana/provisioning/dashboards`.

To update the Grafana dashboard (or create a new one), open it in the Grafana UI, make your changes, see if they work, and save the dashboard JSON model in `monitoring/grafana/provisioning/dashboards`. Overwrite the existing file if it's an update or save it as a new file to create a new dashboard.


## Extra notes
The hosts are a bit annoying because we are using Docker. Instead `localhost` etc., we must use `host.docker.internal`
