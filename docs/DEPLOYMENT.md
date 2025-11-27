# SkillFlow Deployment Guide

## Prerequisites

- Kubernetes cluster (1.25+)
- kubectl configured
- Helm 3.x
- Terraform 1.5+
- Docker
- Git

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/vern/skillflow.git
cd skillflow
```

### 2. Local Development

Start all services locally with Docker Compose:

```bash
make dev-up
```

Run database migrations:

```bash
make migrate-up
```

Start the API server:

```bash
make run
```

API will be available at `http://localhost:8080`

### 3. Build Docker Image

```bash
make docker-build
```

### 4. Deploy Infrastructure with Terraform

Navigate to terraform directory:

```bash
cd deployments/terraform
```

Create `terraform.tfvars`:

```hcl
kubeconfig_path          = "~/.kube/config"
cluster_name             = "skillflow-cluster"
domain_name              = "skillflow.local"
environment              = "prod"
docker_registry          = "registry.skillflow.local"
letsencrypt_email        = "admin@skillflow.local"
keycloak_admin_password  = "your-secure-password"
```

Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

This will deploy:
- **Traefik** - Ingress controller
- **Cert-Manager** - TLS certificate management
- **Prometheus** - Metrics and monitoring
- **Grafana** - Observability dashboards
- **Loki** - Log aggregation
- **Argo CD** - GitOps continuous deployment
- **Keycloak** - Identity and access management

### 5. Deploy Application with Kubernetes

Create namespace and apply manifests:

```bash
kubectl apply -f deployments/kubernetes/namespace.yaml
kubectl apply -f deployments/kubernetes/secret.yaml
kubectl apply -f deployments/kubernetes/configmap.yaml
kubectl apply -f deployments/kubernetes/database.yaml
kubectl apply -f deployments/kubernetes/deployment.yaml
kubectl apply -f deployments/kubernetes/ingress.yaml
```

Or use make:

```bash
make k8s-apply
```

### 6. Deploy with Argo CD (GitOps)

Apply Argo CD application:

```bash
kubectl apply -f deployments/argocd/project.yaml
kubectl apply -f deployments/argocd/application.yaml
```

Get Argo CD admin password:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

Access Argo CD UI at `https://argocd.skillflow.local`

### 7. Configure Keycloak

Access Keycloak admin console at `https://keycloak.skillflow.local`

1. Login with admin credentials
2. Create `skillflow` realm
3. Create clients for:
   - skillflow-api
   - argocd
   - grafana
4. Create users and roles

### 8. Access Services

Once deployed, services will be available at:

- **SkillFlow API**: `https://api.skillflow.local`
- **Argo CD**: `https://argocd.skillflow.local`
- **Grafana**: `https://grafana.skillflow.local`
- **Prometheus**: `https://prometheus.skillflow.local`
- **Traefik Dashboard**: `https://traefik.skillflow.local`
- **Keycloak**: `https://keycloak.skillflow.local`

## Monitoring

### Prometheus Metrics

Prometheus scrapes metrics from:
- SkillFlow API (`/metrics` endpoint)
- Kubernetes cluster components
- Traefik ingress controller
- PostgreSQL
- Redis

### Grafana Dashboards

Pre-configured dashboards available in Grafana:
- Kubernetes Cluster Overview
- Kubernetes Pods
- Traefik Dashboard
- PostgreSQL Database
- Redis Cache
- SkillFlow API Metrics

### Logs with Loki

View logs in Grafana using Loki datasource:
- Application logs
- Kubernetes pod logs
- System logs

Query example:
```
{namespace="skillflow", app="skillflow-api"}
```

## Scaling

### Horizontal Pod Autoscaler

HPA is configured to scale based on:
- CPU utilization (70%)
- Memory utilization (80%)

Scale range: 3-10 replicas

Manually scale:
```bash
kubectl scale deployment skillflow-api --replicas=5 -n skillflow
```

### Database Scaling

For production, consider:
- PostgreSQL replication
- Connection pooling (PgBouncer)
- Read replicas

## Backup and Recovery

### Database Backup

```bash
kubectl exec -n skillflow postgres-0 -- pg_dump -U skillflow skillflow > backup.sql
```

### Restore Database

```bash
kubectl exec -i -n skillflow postgres-0 -- psql -U skillflow skillflow < backup.sql
```

### Volume Snapshots

Use your cloud provider's volume snapshot feature for persistent volumes.

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -n skillflow
kubectl describe pod <pod-name> -n skillflow
kubectl logs <pod-name> -n skillflow
```

### Check Services

```bash
kubectl get svc -n skillflow
kubectl get ingress -n skillflow
```

### Check Argo CD Sync Status

```bash
kubectl get applications -n argocd
argocd app get skillflow-app
argocd app sync skillflow-app
```

### Check Prometheus Targets

Visit `https://prometheus.skillflow.local/targets`

### Database Connection Issues

```bash
kubectl exec -it postgres-0 -n skillflow -- psql -U skillflow
```

## Security Best Practices

1. **Change default passwords** in secrets
2. **Enable RBAC** for Kubernetes
3. **Use network policies** to restrict pod communication
4. **Enable Pod Security Standards**
5. **Regular security updates** for base images
6. **Scan images** for vulnerabilities
7. **Use secrets management** (HashiCorp Vault, Sealed Secrets)
8. **Enable audit logging**
9. **Implement rate limiting** in Traefik
10. **Use TLS** for all communications

## CI/CD Pipeline

### GitHub Actions Example

```yaml
name: Deploy SkillFlow

on:
  push:
    branches: [main]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker image
        run: make docker-build
      
      - name: Push to registry
        run: make docker-push
      
      - name: Update Argo CD
        run: |
          argocd app sync skillflow-app
```

## Maintenance

### Update Dependencies

```bash
go get -u ./...
go mod tidy
```

### Update Helm Charts

```bash
helm repo update
terraform plan
terraform apply
```

### Database Migrations

```bash
make migrate-up
```

## Support

For issues and questions:
- GitHub Issues: https://github.com/vern/skillflow/issues
- Documentation: https://docs.skillflow.local

## License

MIT License - see LICENSE file for details
