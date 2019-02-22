# Middlewares

This repository contains middleware that is used by multiple services in the Influenzanet project to minimize the code duplication and maintain changeability. All middleware is written in [Go](https://golang.org/).

The middleware is written for the [Gin Web Framework](https://github.com/gin-gonic/gin).

Dependancy management uses [Dep](https://github.com/golang/dep).

## Installation

Download the package using

```sh
go get -u https://github.com/Influenzanet/middlewares
```

To use the middleware functions import the middleware package into your go project:

```go
import (
  middlewares "github.com/Influenzanet/middlewares"
)
```

To apply a middleware to a route using Gin:

```go
router := gin.Default()
v1 := router.Group("/v1")
v1.Use(middlewares.RequirePayload())
```

## TODO

Describe functionality of each middleware.