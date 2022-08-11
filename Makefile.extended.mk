CURRENT_DIRECTORY := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
INFRA_DIRECTORY   := $(realpath $(CURRENT_DIRECTORY)/../infra)
SERVICE_NAME      := $(notdir $(CURRENT_DIRECTORY))
CONFIG_SWARM      := docker-compose.swarm.yml
REGISTRY_HOST     := registry.services.phlx.ru
CLUSTER           := swarm
DOTENV            := .env
SERVICES          := server worker

.PHONY: path
# Show command how add Go binaries to PATH making it accessible after `go install ...`
path:
	@echo 'export PATH="$$PATH:$$(go env GOPATH)/bin"'

.PHONY: run-server
# Run cmd/server
run-server:
	@go run cmd/server/wire_gen.go cmd/server/main.go -conf=./configs -dotenv=.env.local

.PHONY: run-worker
# Run cmd/worker
run-worker:
	@go run cmd/worker/wire_gen.go cmd/worker/main.go -conf=./configs -dotenv=.env.local

.PHONY: vendor
# Make ./vendor folder with dependencies
vendor:
	@go mod tidy && go mod vendor && go mod verify

.PHONY: gen
# Makes go generate ./...
gen:
	@go generate ./...

.PHONY: test
# Makes go test ./...
test:
	@go test -race -parallel 10 ./...

.PHNOY: wire
# Wire dependencies with google/wire
wire:
	@go run -mod=mod github.com/google/wire/cmd/wire ./...

.PHONY: ent
# Run ent for generate schema
ent:
	@go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema

.PHONY: lint
# Run linter fo Golang files
lint:
	@docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run

.PHONY: update
# Update service in Docker Swarm without downtime
update:
	@set -e; for service in ${SERVICES}; \
		do docker pull ${REGISTRY_HOST}/${SERVICE_NAME}_$${service}:latest \
			&& docker service update \
			--with-registry-auth \
			--image ${REGISTRY_HOST}/${SERVICE_NAME}_$${service}:latest \
			${CLUSTER}_${SERVICE_NAME}_$${service} ; \
	done

.PHONY: deploy
# Deploy to Docker Swarm
deploy:
	@env \
		$$(cat ${INFRA_DIRECTORY}/${DOTENV} | sed '/^[[:blank:]]*#/d;s/#.*//' | xargs) \
		docker stack deploy \
		--orchestrator swarm \
		--with-registry-auth \
		-c "${CURRENT_DIRECTORY}"/${CONFIG_SWARM} \
		${CLUSTER}

.PHONY: undeploy
# Remove service from Docker Swarm
undeploy:
	@set -e; for service in ${SERVICES}; \
		do docker service rm ${CLUSTER}_${SERVICE_NAME}_$${service} ; \
	done

.PHONY: push
# Build and push image to registry
push:
	@set -e; for service in ${SERVICES}; \
		do docker build -t ${REGISTRY_HOST}/${SERVICE_NAME}_$${service}:latest \
			-f ${CURRENT_DIRECTORY}/Dockerfile-$${service} ${CURRENT_DIRECTORY}/. \
			&& docker push ${REGISTRY_HOST}/${SERVICE_NAME}_$${service}:latest ; \
	done

.PHONY: push
# Build and push image to registry
pull:
	@set -e; for service in ${SERVICES}; \
		do docker pull ${REGISTRY_HOST}/${SERVICE_NAME}_$${service}:latest ; \
	done

.PHONY: env
# Display environment variables from infra .env
env:
	@echo $$(cat ${INFRA_DIRECTORY}/${DOTENV} | sed '/^[[:blank:]]*#/d;s/#.*//' | xargs)

.PHONY: logs
# Display Docker Swarm container logger
logs:
	@docker logs "$$(docker ps -q -f name=${CLUSTER}_${SERVICE_NAME})"
