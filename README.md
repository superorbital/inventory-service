# inventory-service

A demo API that implements a very simple inventory service.

This directory contains an example server using the OpenAPI code generator which implements the OpenAPI [inventory](./inventory-openapi.yaml) example. This is hard forked and heavily inspired by the example pet store code from [deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen/tree/master/examples/petstore-expanded).

Among other things, we use this to provide an example backend for a simple custom Terraform provider.

## Usage

```sh
$ docker compose up -d

$ curl -X GET 127.0.0.1:8080/items
[]

$ curl -H 'Content-Type: application/json' -X POST -d '{"name":"car", "tag":"mustang"}' 127.0.0.1:8080/items
{"id":1000,"name":"car","tag":"mustang"}

curl -X GET 127.0.0.1:8080/items/1000
{"id":1000,"name":"car","tag":"mustang"}

$ curl -X GET 127.0.0.1:8080/items
[{"id":1000,"name":"car","tag":"mustang"}]

$ curl -H 'Content-Type: application/json' -X PUT -d '{"name":"car", "tag":"shelby"}' 127.0.0.1:8080/items/1000
[{"id":1000,"name":"car","tag":"shelby"}]

$ curl -H 'Content-Type: application/json' -X DELETE 127.0.0.1:8080/items/1002

$ curl -X GET 127.0.0.1:8080/items
[]

$ docker compose down
```
>>>>>>> f602e8e (Initial codebase)
