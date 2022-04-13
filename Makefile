# Targets
.PHONY: goose_up
.PHONY: go_get go_build go_test_integration go_test_unit
.PHONY: docker_build_image
.PHONY: service_build service_start_test service_start_demo service_start_dependencies service_stop service_restart
.PHONY: service_run_tests service_test

#
# Project-specific variables
#
# Service name. Used for binary name, docker-compose service name, etc...
SERVICE=currency-exchange-service
# Enable Go Modules.
GO111MODULE=on
# Enable Go Proxy.
GOPROXY=https://proxy.golang.org
# Service go module import path.
GO_SERVICE_IMPORT_PATH=$(shell go list .)

#
# General variables
#
# Path to Docker file
PATH_DOCKER_FILE=$(realpath ./build/Dockerfile)
# Path to docker-compose file
PATH_DOCKER_COMPOSE_FILE=$(realpath ./build/docker-compose.yml)
# Service go module import path.
GO_SERVICE_IMPORT_PATH=$(shell go list .)
# Docker compose starting options.
DOCKER_COMPOSE_OPTIONS= -f $(PATH_DOCKER_COMPOSE_FILE)

#
# Linter targets.
#
lint:
	@docker build -f build/linter/Dockerfile \
	-t $(SERVICE)-linter \
	--build-arg PRIVATE_REPO_KEY="$$(cat ~/.ssh/id_rsa)" \
	. \
	&& docker run $(SERVICE)-linter

#
# Go targets.
#
go_get:
	@echo '>>> Getting go modules.'
	@env GOPROXY=$(GOPROXY) go mod download

go_build:
	@echo '>>> Building go binary.'
	@env GOOS=linux GOPROXY=$(GOPROXY) go build -ldflags="-s -w" -o $(SERVICE)

go_test_unit:
	@echo ">>> Running unit tests."
	@env GIN_MODE=release GOPROXY=$(GOPROXY) go test -v -tags unit

go_test_integration:
	@echo ">>> Running integration tests."
	@env GOPROXY=$(GOPROXY) go test -v -p 1 -tags="integration" ./test/integration/...

#
# Docker targets.
#
docker_build_image:
	@echo ">>> Building docker image with service binary."
	docker build \
		-t $(SERVICE) \
		--build-arg GOPROXY=$(GOPROXY) \
		--build-arg GO111MODULE=$(GO111MODULE) \
		--build-arg GO_SERVICE_IMPORT_PATH=$(GO_SERVICE_IMPORT_PATH) \
		--build-arg PRIVATE_REPO_KEY="$$(cat ~/.ssh/id_rsa)" \
		-f $(PATH_DOCKER_FILE) \
		.

#
# Service targets.
#
service_build:
	@env PRIVATE_REPO_KEY="$$(cat ~/.ssh/id_rsa)" docker-compose $(DOCKER_COMPOSE_OPTIONS) build --no-cache

service_start_dependencies:
	@echo ">>> Start all service dependencies."
	@env PRIVATE_REPO_KEY="$$(cat ~/.ssh/id_rsa)" docker-compose $(DOCKER_COMPOSE_OPTIONS) up \
	-d \
	currency-exchange-service-fixer currency-exchange-service-postgres
	@echo ">>> Sleeping 15 seconds until dependencies start."
	@sleep 15

service_start_test: docker_build_image service_build service_start_dependencies goose_up
	@echo ">>> Starting service."
	@echo ">>> Starting up service container."
	@env PRIVATE_REPO_KEY="$$(cat ~/.ssh/id_rsa)" docker-compose $(DOCKER_COMPOSE_OPTIONS) up -d $(SERVICE)-test

service_start_demo: docker_build_image service_build service_start_dependencies goose_up
	@echo ">>> Starting service."
	@echo ">>> Starting up service container."
	@env PRIVATE_REPO_KEY="$$(cat ~/.ssh/id_rsa)" docker-compose $(DOCKER_COMPOSE_OPTIONS) up -d $(SERVICE)-demo

service_stop:
	@echo ">>> Stopping service."
	@docker-compose  $(DOCKER_COMPOSE_OPTIONS) down -v --remove-orphans

service_run_tests:
	@echo ">>> Running tests over service."
	@env PRIVATE_REPO_KEY="$$(cat ~/.ssh/id_rsa)" docker-compose $(DOCKER_COMPOSE_OPTIONS) \
	run integration-tests

service_test: service_stop service_start_test service_run_tests service_stop

#
# Goose migration targets.
#
goose_up:
	@docker-compose $(DOCKER_COMPOSE_OPTIONS) run \
            --rm \
            -v $$PWD:/app/$(GO_SERVICE_IMPORT_PATH) \
            -w /app/$(GO_SERVICE_IMPORT_PATH) \
            goose-migrate
