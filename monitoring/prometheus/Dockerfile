FROM prom/prometheus:v2.33.5

ARG PROMETHEUS_YML='prometheus-dev.yml'

COPY $PROMETHEUS_YML /etc/prometheus/prometheus.yml

EXPOSE 9090
