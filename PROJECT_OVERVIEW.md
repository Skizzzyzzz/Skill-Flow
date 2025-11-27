# SkillFlow Project Structure

```
SkillFlow/
│
├── .github/                        # GitHub-specific configurations
│   └── workflows/
│       └── ci-cd.yml              # CI/CD pipeline configuration
│
├── cmd/                           # Application entrypoints
│   ├── api/
│   │   └── main.go               # Main API server
│   └── migrate/
│       └── main.go               # Database migration tool
│
├── internal/                      # Private application code
│   ├── api/
│   │   ├── handlers/             # HTTP request handlers
│   │   │   ├── handlers.go       # Handler initialization
│   │   │   ├── auth.go           # Authentication handlers
│   │   │   ├── user.go           # User management handlers
│   │   │   ├── post.go           # Post handlers
│   │   │   └── other_handlers.go # Additional handlers
│   │   └── router/
│   │       └── router.go         # Route definitions
│   │
│   ├── auth/                     # Authentication logic
│   │
│   ├── config/
│   │   └── config.go             # Configuration management
│   │
│   ├── domain/
│   │   └── models/
│   │       └── models.go         # Data models (User, Post, etc.)
│   │
│   ├── repository/               # Data access layer
│   │   └── repository.go         # Database repositories
│   │
│   └── service/                  # Business logic layer
│       ├── service.go            # Service initialization
│       ├── auth.go               # Auth service
│       └── user.go               # User service
│
├── pkg/                          # Public reusable packages
│   ├── cache/
│   │   └── redis.go              # Redis cache client
│   │
│   ├── database/
│   │   └── database.go           # Database connection
│   │
│   ├── logger/
│   │   └── logger.go             # Structured logging
│   │
│   ├── middleware/
│   │   └── middleware.go         # HTTP middlewares
│   │
│   └── utils/
│       ├── helpers.go            # Utility functions
│       └── helpers_test.go       # Helper tests
│
├── configs/                      # Configuration files
│   ├── config.yaml               # Main configuration
│   └── .env.example              # Environment variables template
│
├── deployments/                  # Deployment configurations
│   │
│   ├── docker-compose.yml        # Local development setup
│   │
│   ├── kubernetes/               # Kubernetes manifests
│   │   ├── namespace.yaml        # Namespace definition
│   │   ├── configmap.yaml        # Configuration map
│   │   ├── secret.yaml           # Secrets
│   │   ├── deployment.yaml       # Application deployment
│   │   ├── ingress.yaml          # Ingress rules (Traefik)
│   │   ├── database.yaml         # Database deployments
│   │   ├── networkpolicy.yaml    # Network policies
│   │   └── monitoring/           # Monitoring configs
│   │       ├── prometheus-alerts.yaml    # Alert rules
│   │       ├── grafana-dashboard.yaml    # Dashboards
│   │       └── servicemonitor.yaml       # Service monitors
│   │
│   ├── terraform/                # Infrastructure as Code
│   │   ├── main.tf               # Main Terraform config
│   │   ├── variables.tf          # Variable definitions
│   │   ├── outputs.tf            # Output definitions
│   │   ├── cert-manager.tf       # Cert-Manager setup
│   │   ├── traefik.tf            # Traefik ingress controller
│   │   ├── prometheus.tf         # Prometheus monitoring
│   │   ├── loki.tf               # Loki log aggregation
│   │   ├── argocd.tf             # Argo CD GitOps
│   │   └── keycloak.tf           # Keycloak identity provider
│   │
│   └── argocd/                   # Argo CD configurations
│       ├── application.yaml      # Application definition
│       ├── project.yaml          # Project definition
│       └── applicationset.yaml   # Multi-environment setup
│
├── scripts/                      # Utility scripts
│   └── setup.sh                  # Automated setup script
│
├── docs/                         # Documentation
│   ├── API.md                    # API documentation
│   ├── DEPLOYMENT.md             # Deployment guide
│   └── ARCHITECTURE.md           # Architecture overview
│
├── .gitignore                    # Git ignore rules
├── .golangci.yml                 # Linter configuration
├── CONTRIBUTING.md               # Contribution guidelines
├── LICENSE                       # MIT License
├── Makefile                      # Build automation
├── Dockerfile                    # Docker image definition
├── go.mod                        # Go module dependencies
└── README.md                     # Project overview
```

## Key Directories

### `/cmd` - Application Entry Points
Contains the main applications. Each subdirectory is a separate binary.
- `api/`: REST API server
- `migrate/`: Database migration tool

### `/internal` - Private Application Code
Application-specific code that shouldn't be imported by other projects.
- `api/`: HTTP layer (handlers, routers)
- `config/`: Configuration management
- `domain/`: Business entities and models
- `repository/`: Data access layer
- `service/`: Business logic layer

### `/pkg` - Public Libraries
Reusable packages that can be imported by other projects.
- `cache/`: Redis caching
- `database/`: Database connections
- `logger/`: Logging utilities
- `middleware/`: HTTP middlewares
- `utils/`: Helper functions

### `/deployments` - Deployment Configurations
Everything needed to deploy the application.
- `kubernetes/`: K8s manifests for all components
- `terraform/`: Infrastructure as Code
- `argocd/`: GitOps configurations
- `docker-compose.yml`: Local development environment

### `/docs` - Documentation
Comprehensive documentation for developers and operators.

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Storage**: MinIO (S3-compatible)
- **Search**: Elasticsearch 8

### DevOps Stack
- **Container Runtime**: Docker
- **Orchestration**: Kubernetes
- **IaC**: Terraform
- **GitOps**: Argo CD
- **Ingress**: Traefik
- **TLS**: Cert-Manager + Let's Encrypt
- **Monitoring**: Prometheus + Grafana
- **Logging**: Loki + Promtail
- **Auth**: Keycloak (OIDC)

## Quick Start

1. **Local Development**
   ```bash
   make dev-up    # Start all services
   make run       # Run API server
   ```

2. **Build Docker Image**
   ```bash
   make docker-build
   ```

3. **Deploy Infrastructure**
   ```bash
   cd deployments/terraform
   terraform init
   terraform apply
   ```

4. **Deploy Application**
   ```bash
   kubectl apply -f deployments/kubernetes/
   ```

5. **GitOps with Argo CD**
   ```bash
   kubectl apply -f deployments/argocd/
   ```

## Features

- ✅ User authentication (JWT + OIDC)
- ✅ User profiles and connections
- ✅ Posts, comments, reactions
- ✅ Real-time notifications (WebSocket)
- ✅ File upload and storage
- ✅ Full-text search
- ✅ Groups and communities
- ✅ Direct messaging
- ✅ Skills and endorsements
- ✅ Full DevOps automation
- ✅ Complete monitoring and logging
- ✅ Horizontal auto-scaling
- ✅ GitOps deployment

## Access URLs (After Deployment)

- **SkillFlow API**: https://api.skillflow.local
- **Argo CD**: https://argocd.skillflow.local
- **Grafana**: https://grafana.skillflow.local
- **Prometheus**: https://prometheus.skillflow.local
- **Keycloak**: https://keycloak.skillflow.local
- **Traefik Dashboard**: https://traefik.skillflow.local

## Next Steps

1. Review documentation in `/docs`
2. Customize configuration in `/configs`
3. Set up CI/CD pipeline
4. Configure Keycloak realms
5. Add custom dashboards
6. Implement additional features

## License

MIT License - See LICENSE file for details
