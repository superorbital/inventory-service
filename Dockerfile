# docker buildx build --platform  linux/amd64,linux/arm64 --tag superorbital/inventory-service --push .
# https://hub.docker.com/repository/docker/superorbital/inventory-service

FROM golang:1.19-alpine AS build

RUN apk --no-cache add \
    bash \
    gcc \
    musl-dev \
    openssl
RUN mkdir -p /go/src/github.com/superorbital/inventory-service
WORKDIR /go/src/github.com/superorbital/inventory-service
ADD . /go/src/github.com/superorbital/inventory-service
RUN go build --ldflags '-linkmode external -extldflags "-static"' .

FROM alpine:3.14 AS deploy

WORKDIR /

RUN apk --no-cache add curl
COPY --from=build /go/src/github.com/superorbital/inventory-service /

HEALTHCHECK --interval=15s --timeout=3s \
  CMD curl -f 127.0.0.1:8080/items?limit=1 || exit 1

ENTRYPOINT ["/inventory-service"]
CMD ["-port", "8080"]
