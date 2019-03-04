# Middlewares

This package contains middleware that is used by the API-Gateway of the Influenzanet Project to minimize code duplication inside of endpoints and increse maintainability and changeability of middleware code. All middleware is written in [Go](https://golang.org/).

The middleware is written for the [Gin Web Framework](https://github.com/gin-gonic/gin).

## Installation

The package is included in the API Gateway repository and does not require any separate installation procedure.

## Functionality

Here you can find the different middleware functions provided and what functionality they provide.

### Require Payload

Checks whether the request contains a body/payload and blocks the request if there is no payload using the response: `403 Bad Request`.

This middleware should only be used for requests that usually provide payloads such as `POST`.

### Extract Token

Checks for the existence of the Authorization header and extracts the encoded authorization/token string from the header and attaches it to the `Gin Context` as `encodedToken`. Requests that have no Authorization header will be blocked: `403 Bad Request`.

This middleware should only be used for endpoints that require the user to be authenticated.

### Extract URL Token

Checks for the existence of a query parameter called token that is then attached to the `Gin Context` as `urlToken` or the request is blocked with: `403 Bad Request`.

This middleware should only be used for endpoints that require an url token to be present.

### Validate Token

Uses the encoded token that the `Extract Token` middleware added to the `Gin Context` and validates it via a call to the authentication service. If valid the decoded token is also added to the `Gin Context` as `parsedToken`.

This middleware requires the `encodedToken` parameter to be present on the `Gin Context` and should not be used without the `Extract Token` middleware.

## TODO

--
