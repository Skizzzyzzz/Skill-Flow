# Argo CD Installation
resource "helm_release" "argocd" {
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  version          = var.argocd_version
  namespace        = "argocd"
  create_namespace = true

  values = [
    <<-EOT
      global:
        domain: argocd.${var.domain_name}

      server:
        extraArgs:
          - --insecure

        ingress:
          enabled: true
          ingressClassName: traefik
          annotations:
            cert-manager.io/cluster-issuer: letsencrypt-prod
            traefik.ingress.kubernetes.io/router.tls: "true"
          hosts:
            - argocd.${var.domain_name}
          tls:
            - secretName: argocd-tls
              hosts:
                - argocd.${var.domain_name}

        metrics:
          enabled: true
          serviceMonitor:
            enabled: true

      controller:
        metrics:
          enabled: true
          serviceMonitor:
            enabled: true

      repoServer:
        metrics:
          enabled: true
          serviceMonitor:
            enabled: true

      redis:
        metrics:
          enabled: true
          serviceMonitor:
            enabled: true

      configs:
        params:
          server.insecure: true

        cm:
          url: https://argocd.${var.domain_name}
          
          # Enable auto-sync
          application.instanceLabelKey: argocd.argoproj.io/instance

          # OIDC configuration (Keycloak)
          oidc.config: |
            name: Keycloak
            issuer: https://keycloak.${var.domain_name}/realms/skillflow
            clientID: argocd
            clientSecret: $argocd-oidc:clientSecret
            requestedScopes: ["openid", "profile", "email", "groups"]

        rbac:
          policy.default: role:readonly
          policy.csv: |
            p, role:admin, applications, *, */*, allow
            p, role:admin, clusters, *, *, allow
            p, role:admin, repositories, *, *, allow
            p, role:admin, projects, *, *, allow
            
            g, argocd-admins, role:admin
    EOT
  ]
}

# Wait for Argo CD to be ready
resource "null_resource" "wait_for_argocd" {
  depends_on = [helm_release.argocd]

  provisioner "local-exec" {
    command = "kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd"
  }
}
