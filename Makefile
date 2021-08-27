GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=app
VERSION_FILE=VERSION
C_VERSION=`cat $(VERSION_FILE)`
SERVICE_PORT?=8001

.PHONY: all test build 

all: help

## Build:
build: ## Build your project and put the output binary in out/bin/
	mkdir -p out/bin
	GOOS=linux GOARCH=amd64  $(GOCMD) build -ldflags "-s" -a -installsuffix cgo -o $(BINARY_NAME) main.go

clean: ## Remove build related file
	rm -fr ./bin
	rm -fr ./out
	rm -f app

## Docker:
docker-build: ## Use the dockerfile to build the container
	docker build --rm --tag $(BINARY_NAME):$(C_VERSION) .


docker-lint: ## Use Hadolint to check for common errors
	docker run --rm -i  hadolint/hadolint hadolint  - < ./Dockerfile 
kind-create: 
	kind create cluster

kind-deploy-app:
	kubectl kustomize k8s/overlay/demo 
	
k8s-kind-deploy-work:
	kind-create
	kind-deploylfconfig.me.app
