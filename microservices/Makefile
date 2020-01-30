APPS := $(shell go list ./... |grep cmd)
GIT_VERSION := $(shell git --no-pager describe --tags --always --dirty)
GIT_DATE := $(firstword $(shell git --no-pager show --date=short --format="%ad" --name-only))
CURR_OS := $(shell go env GOOS)
CURR_ARCH := $(shell go env GOARCH)
GO_PLAT := ${CURR_OS}/${CURR_ARCH}
# Use linker flags to provide version/build settings
LDFLAGS=-ldflags '-s -w -X=main.Version=$(GIT_VERSION) -X=main.Date=$(GIT_DATE) -extldflags "-static"'

REPO=carlosedp
VERSION=latest
PLATFORMS=linux/amd64,linux/arm64,linux/riscv64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
noop=
space = $(noop) $(noop)
comma = ,

.PHONY: all $(APPS) $(PLATFORMS) docker $(PLATFORMS)-docker clean

all: help

## all-build: builds all app binaries for current platform
all-build: $(APPS)

$(APPS):
	CGO_ENABLED=0 go build -a -installsuffix cgo $(LDFLAGS) -o $(shell basename $@) $@

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 go build -a -installsuffix cgo $(LDFLAGS) -o $(shell basename $@)-$(os)-$(arch) $<

## linux-amd64: builds all apps for linux/amd64
linux-amd64:
	$(foreach APP,$(APPS),$(shell CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo $(LDFLAGS) -o ${APP} ./cmd/${APP}))

## docker: builds docker images for all apps
docker: linux-amd64
	$(foreach APP,$(APPS),$(shell docker build -t ${REPO}/microservices-demo-${APP}:${VERSION} --platform=linux/amd64 --build-arg=BIN=${APP} .))

## push: pushes all images to REPO
push: docker
	$(foreach APP,$(APPS),$(shell docker push ${REPO}/microservices-demo-${APP}:${VERSION}))

## docker-multi: builds multi-arch Docker images for listed platforms
docker-multi:
	$(shell DOCKER_BUILDKIT=1 docker buildx build -t ${REPO}/microservices-demo-${APP}:${VERSION} --platform=${PLATFORMS} --push .)

clean:
	go clean

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	@echo "\n"
	@echo "Apps: ${APPS}"
	@echo "Platforms: ${PLATFORMS}"