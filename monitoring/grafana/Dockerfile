FROM grafana/grafana:8.4.3

ENV PROMETHEUS_HOST='host.docker.internal:9090'

COPY provisioning /etc/grafana/provisioning

EXPOSE 3000
