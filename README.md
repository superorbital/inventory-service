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

$ curl -X GET 127.0.0.1:8080/items/1000
{"id":1000,"name":"car","tag":"mustang"}

$ curl -X GET 127.0.0.1:8080/items
[{"id":1000,"name":"car","tag":"mustang"}]

$ curl -H 'Content-Type: application/json' -X PUT -d '{"name":"car", "tag":"shelby"}' 127.0.0.1:8080/items/1000
[{"id":1000,"name":"car","tag":"shelby"}]

$ curl -H 'Content-Type: application/json' -X DELETE 127.0.0.1:8080/items/1000

$ curl -X GET 127.0.0.1:8080/items
[]

$ docker compose down
```

## Releasing

```sh
$ go install github.com/mitchellh/gox@latest
$ gox -osarch='darwin/amd64 darwin/arm64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm linux/arm64 linux/mips linux/mips64 linux/mips64le linux/mipsle linux/s390x netbsd/386 netbsd/amd64 netbsd/arm openbsd/386 openbsd/amd64 windows/386 windows/amd64' -output './bin/builds/inventory-service_{{.OS}}_{{.Arch}}'
```

- Create a release in Github with the resulting binaries.
- The container image will be built, tagged latest, and pushed to Docker Hub with eadch merge to main.

## Pre-Commit Hooks

- See: [pre-commit](https://pre-commit.com/)
  - [pre-commit/pre-commit-hooks](https://github.com/pre-commit/pre-commit-hooks)
  - [antonbabenko/pre-commit-terraform](https://github.com/antonbabenko/pre-commit-terraform)

### Install

#### Local Install (macOS)

- **IMPORTANT**: All developers committing any code to this repo, should have these pre-commit hooks installed locally. Github actions may also run these at some point, but it is generally faster and easier to run them locally, in most cases.

```sh
brew install pre-commit jq shellcheck shfmt git-secrets go-critic golangci-lint
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.15.0
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
go install github.com/BurntSushi/toml/cmd/tomlv@master
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install github.com/orijtech/structslop/cmd/structslop
go install github.com/sqs/goreturns@latest
go install golang.org/x/lint/golint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install mvdan.cc/gofumpt@latest

mkdir -p ${HOME}/.git-template/hooks
git config --global init.templateDir ${HOME}/.git-template
```

- Close and reopen your terminal
- Make sure that you run these commands from the root of this git repo!

```sh
cd inventory-service
pre-commit init-templatedir -t pre-commit ${HOME}/.git-template
pre-commit install
```

- Test it

```sh
pre-commit run -a
git diff
```

### Checks

See:

- [.pre-commit-config.yaml](./.pre-commit-config.yaml)

#### Configuring Hooks

- [.pre-commit-config.yaml](./.pre-commit-config.yaml)
