#!/bin/bash

# SkillFlow Setup Script
# This script automates the initial setup of SkillFlow

set -e

echo "=========================================="
echo "SkillFlow Setup Script"
echo "=========================================="
echo ""

# Check prerequisites
echo "Checking prerequisites..."

command -v kubectl >/dev/null 2>&1 || { echo "kubectl is required but not installed. Aborting." >&2; exit 1; }
command -v helm >/dev/null 2>&1 || { echo "helm is required but not installed. Aborting." >&2; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo "terraform is required but not installed. Aborting." >&2; exit 1; }

echo "✓ All prerequisites met"
echo ""

# Get configuration
read -p "Enter domain name (default: skillflow.local): " DOMAIN
DOMAIN=${DOMAIN:-skillflow.local}

read -p "Enter environment (dev/staging/prod, default: prod): " ENV
ENV=${ENV:-prod}

read -p "Enter Let's Encrypt email: " LE_EMAIL
if [ -z "$LE_EMAIL" ]; then
    echo "Let's Encrypt email is required. Aborting."
    exit 1
fi

read -sp "Enter Keycloak admin password: " KC_PASSWORD
echo ""
if [ -z "$KC_PASSWORD" ]; then
    echo "Keycloak admin password is required. Aborting."
    exit 1
fi

echo ""
echo "Configuration:"
echo "  Domain: $DOMAIN"
echo "  Environment: $ENV"
echo "  Let's Encrypt Email: $LE_EMAIL"
echo ""

read -p "Proceed with installation? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

# Create terraform.tfvars
echo ""
echo "Creating Terraform configuration..."
cat > deployments/terraform/terraform.tfvars <<EOF
kubeconfig_path          = "$HOME/.kube/config"
cluster_name             = "skillflow-cluster"
domain_name              = "$DOMAIN"
environment              = "$ENV"
docker_registry          = "registry.$DOMAIN"
letsencrypt_email        = "$LE_EMAIL"
keycloak_admin_password  = "$KC_PASSWORD"
EOF

echo "✓ Terraform configuration created"

# Deploy infrastructure
echo ""
echo "Deploying infrastructure with Terraform..."
cd deployments/terraform
terraform init
terraform plan
terraform apply -auto-approve

echo "✓ Infrastructure deployed"
cd ../..

# Wait for services to be ready
echo ""
echo "Waiting for services to be ready..."
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=traefik -n traefik --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=cert-manager -n cert-manager --timeout=300s
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=prometheus -n monitoring --timeout=300s

echo "✓ Core services are ready"

# Generate secrets
echo ""
echo "Generating secrets..."

JWT_SECRET=$(openssl rand -base64 32)
DB_PASSWORD=$(openssl rand -base64 24)
OIDC_CLIENT_SECRET=$(openssl rand -base64 32)
MINIO_ACCESS_KEY="minioadmin"
MINIO_SECRET_KEY=$(openssl rand -base64 24)

# Update secrets
kubectl delete secret skillflow-secrets -n skillflow 2>/dev/null || true
kubectl create secret generic skillflow-secrets -n skillflow \
  --from-literal=jwt-secret="$JWT_SECRET" \
  --from-literal=db-password="$DB_PASSWORD" \
  --from-literal=oidc-client-secret="$OIDC_CLIENT_SECRET" \
  --from-literal=minio-access-key="$MINIO_ACCESS_KEY" \
  --from-literal=minio-secret-key="$MINIO_SECRET_KEY"

echo "✓ Secrets generated and stored"

# Deploy application
echo ""
echo "Deploying SkillFlow application..."
kubectl apply -f deployments/kubernetes/

echo "✓ Application deployed"

# Deploy Argo CD application
echo ""
echo "Deploying Argo CD GitOps configuration..."
kubectl apply -f deployments/argocd/

echo "✓ Argo CD configured"

# Get access information
echo ""
echo "=========================================="
echo "Installation Complete!"
echo "=========================================="
echo ""
echo "Access URLs:"
echo "  SkillFlow API: https://api.$DOMAIN"
echo "  Argo CD: https://argocd.$DOMAIN"
echo "  Grafana: https://grafana.$DOMAIN"
echo "  Prometheus: https://prometheus.$DOMAIN"
echo "  Keycloak: https://keycloak.$DOMAIN"
echo ""
echo "Get Argo CD admin password:"
echo "  kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath=\"{.data.password}\" | base64 -d"
echo ""
echo "Keycloak admin credentials:"
echo "  Username: admin"
echo "  Password: $KC_PASSWORD"
echo ""
echo "Important: Save the following secrets securely:"
echo "  JWT Secret: $JWT_SECRET"
echo "  Database Password: $DB_PASSWORD"
echo "  OIDC Client Secret: $OIDC_CLIENT_SECRET"
echo "  MinIO Access Key: $MINIO_ACCESS_KEY"
echo "  MinIO Secret Key: $MINIO_SECRET_KEY"
echo ""
echo "Next steps:"
echo "1. Configure DNS to point to your cluster's load balancer"
echo "2. Configure Keycloak realm and clients"
echo "3. Run database migrations: make migrate-up"
echo "4. Build and push Docker image: make docker-build docker-push"
echo ""
echo "For detailed instructions, see docs/DEPLOYMENT.md"
echo ""
