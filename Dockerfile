FROM golang:alpine as builder
RUN mkdir -p /go/src/github.com/influenzanet/api-gateway
ADD . /go/src/github.com/influenzanet/api-gateway/
WORKDIR /go/src/github.com/influenzanet/api-gateway
RUN apk add --no-cache git && echo "installing go packages.." && while read line; do echo "$line" && go get "$line"; done < packages.txt
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o api-gateway .
FROM scratch
COPY --from=builder /go/src/github.com/influenzanet/api-gateway/api-gateway /app/
COPY ./config.yaml /app/
WORKDIR /app
EXPOSE 3000:3000
CMD ["./api-gateway"]
