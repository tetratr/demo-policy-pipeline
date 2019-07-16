# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
COMMIT=$$(git rev-parse --short HEAD)
BINARY_NAME=demo-policy-pipeline

all: test build

build: deps
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o build/$(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -f build/$(BINARY_NAME)
	rm -f $(BINARY_NAME)

deps:
	dep ensure -update

docker-build:
	docker build -t remiphilippe/demo-policy-pipeline .
	docker tag remiphilippe/demo-policy-pipeline remiphilippe/demo-policy-pipeline:$(COMMIT)

docker-publish:
	docker push remiphilippe/demo-policy-pipeline:$(COMMIT)
	docker push remiphilippe/demo-policy-pipeline:latest