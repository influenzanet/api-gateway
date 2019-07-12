.PHONY: build test docker

DOCKER_OPTS ?= --rm

help:
	@echo "Service building targets"
	@echo "	 build : build service command"
	@echo "  test  : run test suites"
	@echo "  docker: build docker image"
	@echo "Env:"
	@echo "  DOCKER_OPTS : default docker build options (default : $(DOCKER_OPTS))"
	@echo "  TEST_ARGS : Arguments to pass to go test call"

build:
	go build .

test:
	go test $(TEST_ARGS)

docker:
	docker build $(DOCKER_OPTS) .