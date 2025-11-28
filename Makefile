.PHONY: build run test clean docker-build docker-push deploy

APP_NAME=skillflow
VERSION?=latest
DOCKER_REGISTRY?=registry.skillflow.local
GO=go
GOFLAGS=-mod=vendor

build:
	$(GO) build $(GOFLAGS) -o bin/$(APP_NAME) cmd/api/main.go

run:
	$(GO) run cmd/api/main.go

test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

coverage:
	$(GO) tool cover -html=coverage.out

clean:
	rm -rf bin/
	rm -f coverage.out

docker-build:
	docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION) .

docker-push:
	docker push $(DOCKER_REGISTRY)/$(APP_NAME):$(VERSION)

# Kubernetes
k8s-apply:
	kubectl apply -f deployments/kubernetes/

k8s-delete:
	kubectl delete -f deployments/kubernetes/

# Terraform
tf-init:
	cd deployments/terraform && terraform init

tf-plan:
	cd deployments/terraform && terraform plan

tf-apply:
	cd deployments/terraform && terraform apply

# Development
dev-up:
	docker-compose -f deployments/docker-compose.yml up -d

dev-down:
	docker-compose -f deployments/docker-compose.yml down

# Database migrations
migrate-up:
	@if [ ! -f configs/config.local.yaml ]; then \
		echo "Creating local config from config.yaml..."; \
		cp configs/config.yaml configs/config.local.yaml; \
		sed -i '' 's/host: postgres/host: localhost/g' configs/config.local.yaml; \
		sed -i '' 's/host: redis/host: localhost/g' configs/config.local.yaml; \
		sed -i '' 's/endpoint: minio:9000/endpoint: localhost:9000/g' configs/config.local.yaml; \
		sed -i '' 's|http://elasticsearch:9200|http://localhost:9200|g' configs/config.local.yaml; \
	fi
	CONFIG_PATH=configs/config.local.yaml $(GO) run cmd/migrate/main.go up

migrate-down:
	CONFIG_PATH=configs/config.local.yaml $(GO) run cmd/migrate/main.go down

vendor:
	$(GO) mod tidy
	$(GO) mod vendor

lint:
	golangci-lint run ./...
