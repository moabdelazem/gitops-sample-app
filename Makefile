APP_NAME    := gitops-app
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME  ?= $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS     := -s -w \
               -X github.com/moabdelazem/gitops-sample-app/pkg/version.Version=$(VERSION) \
               -X github.com/moabdelazem/gitops-sample-app/pkg/version.GitCommit=$(GIT_COMMIT) \
               -X github.com/moabdelazem/gitops-sample-app/pkg/version.BuildTime=$(BUILD_TIME)

DOCKER_USERNAME := "moabdelazem"
DOCKER_IMG  := $(APP_NAME)
DOCKER_TAG  := $(VERSION)

.PHONY: build run clean docker-build docker-run help

build:
	@echo "Building $(APP_NAME) $(VERSION) ($(GIT_COMMIT))..."
	go build -ldflags="$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/main.go
	@echo "Binary: bin/$(APP_NAME)"

run: build
	@echo "Starting $(APP_NAME)..."
	./bin/$(APP_NAME)

clean:
	@rm -rf bin/
	@echo "Cleaned."

docker-build:
	@echo "Building Docker image $(DOCKER_IMG):$(DOCKER_TAG)..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(DOCKER_USERNAME)/$(DOCKER_IMG):$(DOCKER_TAG) \
		-t $(DOCKER_USERNAME)/$(DOCKER_IMG):latest .

docker-run:
	@echo "Running $(DOCKER_IMG):$(DOCKER_TAG)..."
	docker run --rm -p 8080:8080 \
		-e APP_ENV=development \
		$(DOCKER_USERNAME)/$(DOCKER_IMG):$(DOCKER_TAG)
