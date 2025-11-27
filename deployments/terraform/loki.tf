# Loki for log aggregation
resource "helm_release" "loki" {
  name             = "loki"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "loki-stack"
  version          = var.loki_version
  namespace        = "monitoring"
  create_namespace = false
  depends_on       = [helm_release.prometheus]

  values = [
    <<-EOT
      loki:
        enabled: true
        persistence:
          enabled: true
          size: 50Gi
        
        config:
          auth_enabled: false
          
          ingester:
            chunk_idle_period: 3m
            chunk_block_size: 262144
            chunk_retain_period: 1m
            max_transfer_retries: 0
            lifecycler:
              ring:
                kvstore:
                  store: inmemory
                replication_factor: 1

          limits_config:
            enforce_metric_name: false
            reject_old_samples: true
            reject_old_samples_max_age: 168h
            ingestion_rate_mb: 10
            ingestion_burst_size_mb: 20

          schema_config:
            configs:
            - from: 2020-10-24
              store: boltdb-shipper
              object_store: filesystem
              schema: v11
              index:
                prefix: index_
                period: 24h

          server:
            http_listen_port: 3100

          storage_config:
            boltdb_shipper:
              active_index_directory: /data/loki/boltdb-shipper-active
              cache_location: /data/loki/boltdb-shipper-cache
              cache_ttl: 24h
              shared_store: filesystem
            filesystem:
              directory: /data/loki/chunks

          chunk_store_config:
            max_look_back_period: 0s

          table_manager:
            retention_deletes_enabled: true
            retention_period: 720h

      promtail:
        enabled: true
        config:
          clients:
            - url: http://loki:3100/loki/api/v1/push

      fluent-bit:
        enabled: false

      grafana:
        enabled: false
        
      prometheus:
        enabled: false
    EOT
  ]
}

# Grafana datasource for Loki
resource "kubectl_manifest" "loki_datasource" {
  depends_on = [helm_release.loki, helm_release.prometheus]

  yaml_body = <<-YAML
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: loki-grafana-datasource
      namespace: monitoring
      labels:
        grafana_datasource: "1"
    data:
      loki-datasource.yaml: |-
        apiVersion: 1
        datasources:
          - name: Loki
            type: loki
            access: proxy
            url: http://loki:3100
            jsonData:
              maxLines: 1000
  YAML
}
