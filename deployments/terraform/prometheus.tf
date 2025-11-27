# Prometheus Operator (kube-prometheus-stack)
resource "helm_release" "prometheus" {
  name             = "kube-prometheus-stack"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  version          = var.prometheus_version
  namespace        = "monitoring"
  create_namespace = true

  values = [
    <<-EOT
      prometheus:
        prometheusSpec:
          retention: 30d
          storageSpec:
            volumeClaimTemplate:
              spec:
                accessModes: ["ReadWriteOnce"]
                resources:
                  requests:
                    storage: 50Gi
          resources:
            requests:
              cpu: 500m
              memory: 2Gi
            limits:
              cpu: 1000m
              memory: 4Gi
          
          # Service monitors for SkillFlow
          additionalScrapeConfigs:
            - job_name: 'skillflow-api'
              kubernetes_sd_configs:
                - role: pod
                  namespaces:
                    names:
                      - skillflow
              relabel_configs:
                - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
                  action: keep
                  regex: true
                - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
                  action: replace
                  target_label: __metrics_path__
                  regex: (.+)
                - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
                  action: replace
                  regex: ([^:]+)(?::\d+)?;(\d+)
                  replacement: $1:$2
                  target_label: __address__

        ingress:
          enabled: true
          ingressClassName: traefik
          annotations:
            cert-manager.io/cluster-issuer: letsencrypt-prod
          hosts:
            - prometheus.${var.domain_name}
          tls:
            - secretName: prometheus-tls
              hosts:
                - prometheus.${var.domain_name}

      grafana:
        adminPassword: ${var.keycloak_admin_password}
        persistence:
          enabled: true
          size: 10Gi
        
        ingress:
          enabled: true
          ingressClassName: traefik
          annotations:
            cert-manager.io/cluster-issuer: letsencrypt-prod
          hosts:
            - grafana.${var.domain_name}
          tls:
            - secretName: grafana-tls
              hosts:
                - grafana.${var.domain_name}

        dashboardProviders:
          dashboardproviders.yaml:
            apiVersion: 1
            providers:
            - name: 'default'
              orgId: 1
              folder: ''
              type: file
              disableDeletion: false
              editable: true
              options:
                path: /var/lib/grafana/dashboards/default

        dashboards:
          default:
            kubernetes-cluster:
              gnetId: 7249
              revision: 1
              datasource: Prometheus
            kubernetes-pods:
              gnetId: 6417
              revision: 1
              datasource: Prometheus
            traefik:
              gnetId: 11462
              revision: 1
              datasource: Prometheus
            postgres:
              gnetId: 9628
              revision: 7
              datasource: Prometheus
            redis:
              gnetId: 11835
              revision: 1
              datasource: Prometheus

        sidecar:
          dashboards:
            enabled: true
            label: grafana_dashboard
          datasources:
            enabled: true
            label: grafana_datasource

      alertmanager:
        alertmanagerSpec:
          storage:
            volumeClaimTemplate:
              spec:
                accessModes: ["ReadWriteOnce"]
                resources:
                  requests:
                    storage: 10Gi

        ingress:
          enabled: true
          ingressClassName: traefik
          annotations:
            cert-manager.io/cluster-issuer: letsencrypt-prod
          hosts:
            - alertmanager.${var.domain_name}
          tls:
            - secretName: alertmanager-tls
              hosts:
                - alertmanager.${var.domain_name}

      kubeStateMetrics:
        enabled: true

      nodeExporter:
        enabled: true

      prometheusOperator:
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 256Mi
    EOT
  ]
}

# ServiceMonitor for SkillFlow API
resource "kubectl_manifest" "skillflow_service_monitor" {
  depends_on = [helm_release.prometheus]

  yaml_body = <<-YAML
    apiVersion: monitoring.coreos.com/v1
    kind: ServiceMonitor
    metadata:
      name: skillflow-api
      namespace: monitoring
      labels:
        app: skillflow-api
    spec:
      selector:
        matchLabels:
          app: skillflow-api
      namespaceSelector:
        matchNames:
          - skillflow
      endpoints:
      - port: http
        path: /metrics
        interval: 30s
  YAML
}
