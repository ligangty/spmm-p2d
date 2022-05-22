# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=./build
BINARY_NAME=$(BUILD_DIR)/hello
BINARY_UNIX=$(BINARY_NAME)_unix
DOCKERCMD=docker
DOCKERBUILD=$(DOCKERCMD) build
DOCKER_REGISTRY=quay.io/ligangty
DOCKER_REPO=helloservice
DOCKER_TAG=latest
DOCKER_FULL=$(DOCKER_REGISTRY)/$(DOCKER_REPO):$(DOCKER_TAG)

.PHONY: build test clean img build-linux


dev: 
		$(GOBUILD) -trimpath -o $(BINARY_NAME) 

test: 
		$(GOTEST) -v ./...

clean: 
		$(GOCLEAN) ./...
		rm -rf $(BUILD_DIR)

img:
		docker build -t $(DOCKER_FULL) -f ./Dockerfile .

# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -o $(BINARY_NAME) -v ./cmd/indycli
