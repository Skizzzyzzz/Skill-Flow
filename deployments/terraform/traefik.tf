# Traefik Ingress Controller
resource "helm_release" "traefik" {
  name             = "traefik"
  repository       = "https://traefik.github.io/charts"
  chart            = "traefik"
  version          = var.traefik_version
  namespace        = "traefik"
  create_namespace = true

  values = [
    <<-EOT
      additionalArguments:
        - "--log.level=INFO"
        - "--accesslog=true"
        - "--metrics.prometheus=true"
        - "--entryPoints.websecure.forwardedHeaders.trustedIPs=10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

      ports:
        web:
          port: 80
          expose: true
          redirectTo: websecure
        websecure:
          port: 443
          expose: true
          tls:
            enabled: true

      service:
        type: LoadBalancer

      ingressRoute:
        dashboard:
          enabled: true
          matchRule: Host(`traefik.${var.domain_name}`)
          entryPoints: ["websecure"]

      persistence:
        enabled: true
        size: 1Gi

      resources:
        requests:
          cpu: 100m
          memory: 128Mi
        limits:
          cpu: 500m
          memory: 256Mi

      metrics:
        prometheus:
          enabled: true
          addEntryPointsLabels: true
          addServicesLabels: true
    EOT
  ]
}
