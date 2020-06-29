PROJECT_NAME := kubepox
VERSION_FILE := ./version/version.go
VERSION := 0.1
REVISION=$(shell git log -1 --pretty=format:"%H")
BUILD_NUMBER := latest
DOCKER_REGISTRY?=aporeto
DOCKER_IMAGE_NAME?=$(PROJECT_NAME)
DOCKER_IMAGE_TAG?=$(BUILD_NUMBER)

deps:
	go get -v ./...

codegen:
	echo 'package version' > $(VERSION_FILE)
	echo '' >> $(VERSION_FILE)
	echo '// VERSION is the version of kubepox' >> $(VERSION_FILE)
	echo 'const VERSION = "$(VERSION)"' >> $(VERSION_FILE)
	echo '' >> $(VERSION_FILE)
	echo '// REVISION is the revision of kubepox' >> $(VERSION_FILE)
	echo 'const REVISION = "$(REVISION)"' >> $(VERSION_FILE)

build: codegen
	CGO_ENABLED=1 go build -a -installsuffix cgo cmd/kubepox/main.go

package: build
	mv main docker/kubepox

docker_build: package
	docker \
		build \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) docker

docker_push: docker_build
	docker \
		push \
		$(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

binary: deps build
	mv main kubepox
