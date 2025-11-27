# Keycloak for OIDC authentication
resource "helm_release" "keycloak" {
  name             = "keycloak"
  repository       = "https://charts.bitnami.com/bitnami"
  chart            = "keycloak"
  version          = "18.0.0"
  namespace        = "keycloak"
  create_namespace = true

  values = [
    <<-EOT
      auth:
        adminUser: admin
        adminPassword: ${var.keycloak_admin_password}

      production: true
      proxy: edge

      postgresql:
        enabled: true
        auth:
          username: keycloak
          password: ${var.keycloak_admin_password}
          database: keycloak
        primary:
          persistence:
            enabled: true
            size: 20Gi

      ingress:
        enabled: true
        ingressClassName: traefik
        annotations:
          cert-manager.io/cluster-issuer: letsencrypt-prod
        hostname: keycloak.${var.domain_name}
        tls: true

      resources:
        requests:
          cpu: 500m
          memory: 1Gi
        limits:
          cpu: 1000m
          memory: 2Gi

      metrics:
        enabled: true
        serviceMonitor:
          enabled: true

      extraEnvVars:
        - name: KEYCLOAK_EXTRA_ARGS
          value: "--spi-login-protocol-openid-connect-legacy-logout-redirect-uri=true"
    EOT
  ]
}

# Keycloak realm configuration
resource "kubectl_manifest" "keycloak_realm" {
  depends_on = [helm_release.keycloak]

  yaml_body = <<-YAML
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: keycloak-realm-config
      namespace: keycloak
    data:
      skillflow-realm.json: |
        {
          "realm": "skillflow",
          "enabled": true,
          "sslRequired": "external",
          "registrationAllowed": true,
          "loginWithEmailAllowed": true,
          "duplicateEmailsAllowed": false,
          "resetPasswordAllowed": true,
          "editUsernameAllowed": false,
          "bruteForceProtected": true,
          "clients": [
            {
              "clientId": "skillflow-api",
              "enabled": true,
              "publicClient": false,
              "protocol": "openid-connect",
              "redirectUris": [
                "https://${var.domain_name}/*",
                "https://api.${var.domain_name}/*"
              ],
              "webOrigins": [
                "https://${var.domain_name}",
                "https://api.${var.domain_name}"
              ],
              "standardFlowEnabled": true,
              "implicitFlowEnabled": false,
              "directAccessGrantsEnabled": true,
              "serviceAccountsEnabled": true,
              "authorizationServicesEnabled": true
            },
            {
              "clientId": "argocd",
              "enabled": true,
              "publicClient": false,
              "protocol": "openid-connect",
              "redirectUris": [
                "https://argocd.${var.domain_name}/*"
              ],
              "webOrigins": [
                "https://argocd.${var.domain_name}"
              ]
            },
            {
              "clientId": "grafana",
              "enabled": true,
              "publicClient": false,
              "protocol": "openid-connect",
              "redirectUris": [
                "https://grafana.${var.domain_name}/*"
              ],
              "webOrigins": [
                "https://grafana.${var.domain_name}"
              ]
            }
          ],
          "roles": {
            "realm": [
              {
                "name": "admin",
                "description": "Administrator role"
              },
              {
                "name": "user",
                "description": "User role"
              },
              {
                "name": "moderator",
                "description": "Moderator role"
              }
            ]
          },
          "groups": [
            {
              "name": "argocd-admins",
              "path": "/argocd-admins"
            }
          ]
        }
  YAML
}
