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

# Use ReCaptcha
USE_RECAPTCHA=true
RECAPTCHA_SECRET=<secret key>

```

## Github Actions

The repository also contains a Github actions script to build and push a docker image to a dockerhub repository. 
The action is a manually triggered workflow dispatch that requires the following secrets to be configured in order to run successfully:

| Secret Name        | Value to be configured           |
| -------------- | -------------------- |
| DOCKER_USER     | Username of the account authorized to push docker image to the dockerhub repository |
| DOCKER_PASSWORD     | Password of the account authorized to push docker image to the dockerhub repository |
| DOCKER_ORGANIZATION     | Organization or collection name that hosts the repository being pushed to |
| DOCKER_PARTICIPANT_API_REPO_NAME     | Name of the participant api dockerhub image repository |
| DOCKER_MANAGEMENT_API_REPO_NAME     | Name of the management api dockerhub image repository |

Once this is configured, navigate to the Actions tab on Github > Docker Image CI > Run Workflow

By default the version to be tagged is picked from the latest release version, but it can also be overriden by a user specified tag name.