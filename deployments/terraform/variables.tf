variable "kubeconfig_path" {
  description = "Path to kubeconfig file"
  type        = string
  default     = "~/.kube/config"
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "skillflow-cluster"
}

variable "domain_name" {
  description = "Domain name for the application"
  type        = string
  default     = "skillflow.local"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "prod"
}

variable "docker_registry" {
  description = "Docker registry URL"
  type        = string
  default     = "registry.skillflow.local"
}

variable "letsencrypt_email" {
  description = "Email for Let's Encrypt certificates"
  type        = string
}

variable "argocd_version" {
  description = "Argo CD version"
  type        = string
  default     = "5.51.0"
}

variable "prometheus_version" {
  description = "Prometheus Operator version"
  type        = string
  default     = "55.0.0"
}

variable "loki_version" {
  description = "Loki stack version"
  type        = string
  default     = "2.9.11"
}

variable "traefik_version" {
  description = "Traefik version"
  type        = string
  default     = "25.0.0"
}

variable "cert_manager_version" {
  description = "Cert Manager version"
  type        = string
  default     = "v1.13.3"
}

variable "keycloak_admin_password" {
  description = "Keycloak admin password"
  type        = string
  sensitive   = true
}
