##########################
# STAGE 1
##########################
FROM golang:1.18-alpine as builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN mkdir -p /go/src/github.com/influenzanet/api-gateway
ENV GO111MODULE=on
ADD . /go/src/github.com/influenzanet/api-gateway/
WORKDIR /go/src/github.com/influenzanet/api-gateway
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR /go/src/github.com/influenzanet/api-gateway/cmd/management-api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app .

##########################
# STAGE 2
##########################
FROM scratch
# we need those on scratch for connecting to mongoDB:
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/influenzanet/api-gateway/cmd/management-api/app /app/
WORKDIR /app
EXPOSE 3232:3232
CMD ["./app"]
