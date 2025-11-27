# SkillFlow - Intranet Social Network

SkillFlow is a modern intranet social network built with Go and deployed using a complete DevOps stack.

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP server)
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Storage**: MinIO (S3-compatible)
- **Search**: Elasticsearch

### DevOps Stack
- **Kubernetes**: Container orchestration
- **Terraform**: Infrastructure as Code
- **Argo CD**: GitOps continuous deployment
- **Prometheus**: Metrics and monitoring
- **Grafana**: Observability dashboards
- **Loki**: Log aggregation
- **Cert-Manager**: Automatic TLS certificates
- **Traefik**: Ingress controller
- **OIDC**: Authentication (Keycloak)

## Project Structure

```
├── cmd/                    # Application entrypoints
│   ├── api/               # Main API server
│   └── migrate/           # Database migrations
├── internal/              # Private application code
│   ├── api/              # HTTP handlers
│   ├── auth/             # Authentication & authorization
│   ├── config/           # Configuration
│   ├── domain/           # Business logic & models
│   ├── repository/       # Data access layer
│   └── service/          # Business services
├── pkg/                   # Public libraries
│   ├── logger/           # Logging utilities
│   ├── middleware/       # HTTP middlewares
│   └── utils/            # Helper functions
├── api/                   # API definitions
│   └── openapi/          # OpenAPI/Swagger specs
├── configs/               # Configuration files
├── deployments/          # Deployment configurations
│   ├── kubernetes/       # K8s manifests
│   ├── terraform/        # Terraform IaC
│   ├── argocd/          # Argo CD applications
│   └── docker-compose.yml # Local development
├── scripts/              # Build and utility scripts
└── docs/                 # Documentation

```

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Kubernetes cluster
- Terraform
- kubectl
- helm

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/vern/skillflow.git
cd skillflow
```

2. Start local services:
```bash
make dev-up
```

3. Run migrations:
```bash
make migrate-up
```

4. Start the API server:
```bash
make run
```

### Build

```bash
make build
```

### Testing

```bash
make test
make coverage
```

### Docker

```bash
make docker-build
make docker-push
```

## Deployment

### Kubernetes

```bash
make k8s-apply
```

### Terraform

```bash
make tf-init
make tf-plan
make tf-apply
```

## Features

- User authentication & authorization (OIDC)
- User profiles and social connections
- Posts, comments, reactions
- Real-time notifications (WebSocket)
- File upload & sharing
- Full-text search
- Activity feeds
- Groups and communities
- Direct messaging
- Skills and endorsements

## API Documentation

API documentation is available at `/api/docs` when the server is running.

## Monitoring

- **Prometheus metrics**: `/metrics`
- **Health check**: `/health`
- **Grafana dashboards**: Available via Traefik ingress

## License

MIT
