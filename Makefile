.PHONY: build test docker-participant-api docker-management-api management-api participant-api

DOCKER_OPTS ?= --rm
VERSION := $(shell git describe --tags --abbrev=0)
TARGET_DIR ?= ./

help:
	@echo "Service building targets"
	@echo "	 build : build service command"
	@echo "  test  : run test suites"
	@echo "  docker-participant-api: build docker image for participant-api"
	@echo "  docker-management-api: build docker image for management-api"
	@echo "Env:"
	@echo "  DOCKER_OPTS : default docker build options (default : $(DOCKER_OPTS))"
	@echo "  TEST_ARGS : Arguments to pass to go test call"

    
management-api:
	go build -o $(TARGET_DIR) ./cmd/management-api

participant-api:
	go build -o $(TARGET_DIR) ./cmd/participant-api

build: management-api participant-api

test:
	go test $(TEST_ARGS)

docker-participant-api:
	docker build -t  github.com/influenzanet/participant-api:$(VERSION)  -f build/docker/participant-api/Dockerfile $(DOCKER_OPTS) .

docker-management-api:
	docker build -t  github.com/influenzanet/management-api:$(VERSION)  -f build/docker/management-api/Dockerfile $(DOCKER_OPTS) .