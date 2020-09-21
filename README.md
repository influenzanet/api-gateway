# API Gateway

## Role

API Gateway provides the public interface for InfluenzaNet system. It exposes operations useable by clients to interact with an Influenzanet instance.

This repository hosts the participant and management API gateways (service mapping HTTP logic to gRPC services).

## Build
### Docker
Dockerfile(s) are located in `build/docker/participant-api` and `build/docker/management-api`. The default Dockerfile is using a multistage build and create a minimal image base on `scratch`.
To trigger the build process using the default docker file call:

``` sh
make docker-participant-api
make docker-managment-api
```

This will use the most recent git tag to tag the docker image.

#### Contribute

Feel free to create your own Dockerfile (e.g. compiling and deploying to specific target images), eventually others may need the same.
You can create a pull request with adding the Dockerfile into `build/docker/*` with a good name that it can be identified well, and add a short description to `build/docker/readme.md` about the purpose and speciality of it.

An example to run your created docker image - with the set environment variables - can be found e.g. [here](build/docker/participant-api/example).

## Settings

Turn endpoints on/off:

``` sh
# Allow deleting participant data
USE_DELETE_PARTICIPANT_DATA_ENDPOINT=true

# to disable signup
DISABLE_SIGNUP_WITH_EMAIL_ENDPOINT=true
```
