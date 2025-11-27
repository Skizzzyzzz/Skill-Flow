# Architecture Overview

## System Architecture

SkillFlow is a cloud-native intranet social network built with modern DevOps practices and microservices principles.

```
┌─────────────────────────────────────────────────────────────────┐
│                         Internet / Users                         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                     ┌───────▼────────┐
                     │  DNS / CDN     │
                     └───────┬────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                     Load Balancer                                │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                  Traefik Ingress Controller                      │
│                    (TLS Termination)                             │
└─────────┬──────────────────┬──────────────────┬─────────────────┘
          │                  │                  │
    ┌─────▼─────┐      ┌─────▼─────┐     ┌─────▼─────┐
    │ SkillFlow │      │  Keycloak │     │  Grafana  │
    │    API    │      │   (OIDC)  │     │Dashboards │
    └─────┬─────┘      └───────────┘     └───────────┘
          │
┌─────────▼────────────────────────────────────────────────────────┐
│                     SkillFlow Application                         │
│  ┌──────────┐  ┌───────────┐  ┌───────────┐  ┌──────────────┐  │
│  │   HTTP   │  │   Auth    │  │  Business │  │  WebSocket   │  │
│  │ Handlers │  │Middleware │  │   Logic   │  │   Server     │  │
│  └────┬─────┘  └─────┬─────┘  └─────┬─────┘  └──────┬───────┘  │
│       └──────────────┼──────────────┘────────────────┘          │
│                      │                                           │
│              ┌───────▼────────┐                                  │
│              │   Services     │                                  │
│              └───────┬────────┘                                  │
│                      │                                           │
│              ┌───────▼────────┐                                  │
│              │  Repositories  │                                  │
│              └───────┬────────┘                                  │
└──────────────────────┼──────────────────────────────────────────┘
                       │
         ┌─────────────┼─────────────┬──────────────┐
         │             │             │              │
    ┌────▼────┐   ┌────▼────┐   ┌───▼───┐    ┌────▼─────┐
    │PostgreSQL│   │  Redis  │   │ MinIO │    │Elasticsearch│
    │         │   │ (Cache) │   │(Files)│    │  (Search)  │
    └─────────┘   └─────────┘   └───────┘    └───────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Monitoring & Logging                          │
│  ┌────────────┐  ┌──────────┐  ┌─────────┐  ┌─────────────┐   │
│  │ Prometheus │─▶│ Grafana  │  │  Loki   │  │ AlertManager│   │
│  └──────┬─────┘  └──────────┘  └────┬────┘  └─────────────┘   │
│         │                            │                          │
│         └────────────────────────────┘                          │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                      GitOps & CI/CD                              │
│  ┌────────────┐           ┌──────────────┐                      │
│  │   Argo CD  │◀──────────│  Git Repo    │                      │
│  │  (GitOps)  │           │ (Manifests)  │                      │
│  └──────┬─────┘           └──────────────┘                      │
│         │                                                        │
│         └──────────▶ Kubernetes Cluster                         │
└─────────────────────────────────────────────────────────────────┘
```

## Component Details

### Frontend Layer

- **Traefik**: Ingress controller handling all incoming traffic
  - TLS termination
  - Routing to services
  - Rate limiting
  - Request/response middleware

- **Cert-Manager**: Automatic TLS certificate management
  - Let's Encrypt integration
  - Certificate renewal
  - Multiple issuers support

### Application Layer

- **SkillFlow API**: Go-based REST API
  - Gin web framework
  - JWT authentication
  - OIDC integration
  - WebSocket support
  - Structured logging
  - Prometheus metrics

### Authentication & Authorization

- **Keycloak**: Identity and Access Management
  - OIDC provider
  - User management
  - Role-based access control
  - SSO capabilities
  - Multi-realm support

### Data Layer

- **PostgreSQL**: Primary relational database
  - User data
  - Posts, comments, reactions
  - Connections and groups
  - Skills and endorsements

- **Redis**: Caching and session storage
  - Session management
  - Cache layer
  - Real-time features

- **MinIO**: Object storage (S3-compatible)
  - File uploads
  - Media storage
  - Avatars and covers

- **Elasticsearch**: Full-text search
  - User search
  - Content search
  - Skills search

### Monitoring & Observability

- **Prometheus**: Metrics collection and storage
  - Application metrics
  - Infrastructure metrics
  - Custom business metrics
  - Alert rules

- **Grafana**: Visualization and dashboards
  - Pre-built dashboards
  - Custom dashboards
  - Alerting
  - Multiple data sources

- **Loki**: Log aggregation
  - Application logs
  - Container logs
  - Query interface
  - Grafana integration

- **AlertManager**: Alert management
  - Alert routing
  - Grouping and deduplication
  - Notification channels

### GitOps & Deployment

- **Argo CD**: Continuous deployment
  - GitOps workflow
  - Automated sync
  - Rollback capabilities
  - Multi-environment support

- **Terraform**: Infrastructure as Code
  - Cluster provisioning
  - Service deployment
  - Configuration management
  - State management

## Data Flow

### User Registration Flow

```
User → Traefik → API → Auth Service → PostgreSQL
                  ↓
            Profile Creation
                  ↓
            JWT Generation
                  ↓
              Response
```

### OIDC Authentication Flow

```
User → Traefik → API → Keycloak
                         ↓
                   User Login
                         ↓
                  Authorization Code
                         ↓
                    API Exchange
                         ↓
                   JWT Creation
                         ↓
                     Response
```

### Post Creation Flow

```
User → Traefik → API → Auth Middleware
                         ↓
                   Post Service
                         ↓
                   PostgreSQL
                         ↓
                  Search Indexing
                         ↓
                   Elasticsearch
                         ↓
                  Cache Invalidation
                         ↓
                      Redis
                         ↓
                 Notification Service
                         ↓
                   WebSocket Push
```

### File Upload Flow

```
User → Traefik → API → File Service
                         ↓
                   MinIO Upload
                         ↓
                   URL Generation
                         ↓
                   Database Record
                         ↓
                     Response
```

## Security Architecture

### Network Security

- **Network Policies**: Restrict pod-to-pod communication
- **TLS Everywhere**: All external communication encrypted
- **Ingress Filtering**: Rate limiting and WAF rules

### Application Security

- **JWT Authentication**: Stateless authentication
- **OIDC Integration**: Centralized identity management
- **Role-Based Access Control**: Fine-grained permissions
- **Input Validation**: Request validation middleware
- **SQL Injection Prevention**: Parameterized queries (GORM)

### Secrets Management

- **Kubernetes Secrets**: Encrypted at rest
- **External Secrets**: Integration ready (Vault, AWS Secrets Manager)
- **Environment Variables**: Sensitive config from secrets

## Scalability

### Horizontal Scaling

- **API Pods**: Auto-scaled based on CPU/Memory
- **Database Replicas**: Read replicas for scaling reads
- **Cache Layer**: Redis cluster for distributed caching
- **Storage**: S3-compatible distributed storage

### Vertical Scaling

- **Resource Requests/Limits**: Defined for all pods
- **Database Tuning**: Connection pooling, query optimization
- **Cache Optimization**: TTL strategies, eviction policies

## High Availability

- **Multi-replica Deployments**: 3+ replicas for API
- **Database Replication**: Master-slave setup
- **Load Balancing**: Traffic distributed across pods
- **Health Checks**: Liveness and readiness probes
- **Graceful Shutdown**: Clean connection handling

## Disaster Recovery

- **Database Backups**: Automated daily backups
- **Volume Snapshots**: Persistent volume backups
- **GitOps Recovery**: Infrastructure code in Git
- **Multi-region Ready**: Architecture supports multi-region

## Performance Optimization

- **Caching Strategy**: Redis for frequently accessed data
- **Database Indexing**: Optimized indexes
- **Connection Pooling**: Efficient database connections
- **CDN Ready**: Static assets can be served via CDN
- **Compression**: Gzip compression enabled

## Technology Stack Summary

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Language | Go 1.21+ | Application development |
| Framework | Gin | HTTP web framework |
| Database | PostgreSQL 15 | Primary data store |
| Cache | Redis 7 | Caching and sessions |
| Storage | MinIO | Object storage |
| Search | Elasticsearch 8 | Full-text search |
| Ingress | Traefik | Load balancing and routing |
| TLS | Cert-Manager | Certificate management |
| Auth | Keycloak | Identity management |
| Monitoring | Prometheus | Metrics collection |
| Visualization | Grafana | Dashboards |
| Logging | Loki | Log aggregation |
| GitOps | Argo CD | Continuous deployment |
| IaC | Terraform | Infrastructure provisioning |
| Container | Docker | Containerization |
| Orchestration | Kubernetes | Container orchestration |
