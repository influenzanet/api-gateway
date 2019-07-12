# API Gateway

## Role

API Gateway provides the public interface for InfluenzaNet system. It exposes operations useable by clients to interact with an Influenzanet instance

This repository is a [Go implementation of influenzanet service](https://github.com/influenzanet/influenzanet/wiki/Go-based-service-organisation)

## Configuration

To run the service expect a configuration file (by default config.yml) defining the endpoints for internal microservices

For example

```yaml
service_urls:
  authentication: "localhost:3201"
  user_management: "localhost:3200"
```

## Configuration Environment
 
 - CONFIG_FILE (optional) location of the configuration file

## Run

Several make targets are available

 - build  : build the service locally
 - test   : run tests (env `TEST_ARGS` to add more arguments )
 - docker : build the docker image (`DOCKER_OPTS` to add more arguments)
 - help : show all available targets