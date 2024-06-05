# Define variables
DOCKER_COMPOSE = docker-compose
ENV_FILE = .env

# Default target
.PHONY: all
all: build up

# Build the Docker images
.PHONY: build
build: build_app build_api

.PHONY: build_app
build_app:
	$(DOCKER_COMPOSE) -f deployment/docker-compose.yml build

.PHONY: build_api
build_api:
	$(DOCKER_COMPOSE) -f deployment/docker-compose.api.yml build

# Start the services
.PHONY: up
up: $(ENV_FILE) up_app up_api

.PHONY: up_app
up_app:
	$(DOCKER_COMPOSE) -f deployment/docker-compose.yml up -d

.PHONY: up_api
up_api:
	$(DOCKER_COMPOSE) -f deployment/docker-compose.api.yml up -d

# Stop the services
.PHONY: down
down: down_app down_api

.PHONY: down_app
down_app:
	$(DOCKER_COMPOSE) -f deployment/docker-compose.yml down

.PHONY: down_api
down_api:
	$(DOCKER_COMPOSE) -f deployment/docker-compose.api.yml down

# Restart the services
.PHONY: restart
restart: down up

# Generate Swagger documentation
.PHONY: swagger
swagger:
	swag init -g cmd/api -o ../../docs

# Check environment variables
$(ENV_FILE):
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Error: .env file not found!"; \
		exit 1; \
	fi
