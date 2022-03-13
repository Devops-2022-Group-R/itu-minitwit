FROM prom/prometheus:v2.33.5

ARG APP_HOST='host.docker.internal:8080'
ARG PROMETHEUS_YML='prometheus-dev.yml'

COPY $PROMETHEUS_YML /etc/prometheus/prometheus.yml

EXPOSE 9090
