output "argocd_server_url" {
  description = "Argo CD server URL"
  value       = "https://argocd.${var.domain_name}"
}

output "grafana_url" {
  description = "Grafana dashboard URL"
  value       = "https://grafana.${var.domain_name}"
}

output "prometheus_url" {
  description = "Prometheus server URL"
  value       = "https://prometheus.${var.domain_name}"
}

output "traefik_dashboard_url" {
  description = "Traefik dashboard URL"
  value       = "https://traefik.${var.domain_name}"
}

output "keycloak_url" {
  description = "Keycloak URL"
  value       = "https://keycloak.${var.domain_name}"
}

output "skillflow_api_url" {
  description = "SkillFlow API URL"
  value       = "https://api.${var.domain_name}"
}
