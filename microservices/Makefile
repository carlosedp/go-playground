APPS := $(go list ./... |grep cmd)
# APPS := echo-app publisher formatter client
GIT_VERSION := $(shell git --no-pager describe --tags --always --dirty)
GIT_DATE := $(firstword $(shell git --no-pager show --date=short --format="%ad" --name-only))

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags '-X=main.Version=$(GIT_VERSION) -X=main.Build=$(BUILD) -extldflags "-static"'

REPO=carlosedp
VERSION=latest
PLATFORMS=linux/amd64 linux/arm64 linux/riscv64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
noop=
space = $(noop) $(noop)
comma = ,

.PHONY: all $(APPS) $(PLATFORMS) docker $(PLATFORMS)-docker clean

all: $(APPS)

$(APPS):
	CGO_ENABLED=0 go build -a -installsuffix cgo $(LDFLAGS) -o $(shell basename $@) ./cmd/$@

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-X main.version $(GIT_VERSION) -extldflags "-static"' -o ${APP}-$(os)-$(arch) $<

docker: $(PLATFORMS)
	docker buildx build -t ${REPO}/${APP} --platform=$(subst $(space),$(comma),$(PLATFORMS)) --push .

docker-legacy: $(PLATFORMS)
	DOCKER_BUILDKIT=1 docker build -t ${REPO}/${APP} .

docker-multi: $(PLATFORMS)
	$(foreach PLAT,$(PLATFORMS),$(shell DOCKER_BUILDKIT=1 docker build -t ${REPO}/${APP}:${}{VERSION} --platform=$(PLAT) .))

release: docker

clean:
	rm -rf ${APP}-*

