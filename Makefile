NAME       := alerterator
TAG        := navikt/${NAME}
LATEST     := ${TAG}:latest
ROOT_DIR   := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))


.PHONY: build docker docker-push local install test codegen-crd codegen-updater

build:
	cd cmd/alerterator && go build

docker:
    VERSION := $(shell ./version)
	docker image build -t ${TAG}:${VERSION} -t ${TAG} -t ${NAME} -t ${LATEST} -f Dockerfile .
	docker image push ${TAG}:${VERSION}
	docker image push ${LATEST}

local:
	go run cmd/alerterator/main.go --logtostderr --kubeconfig=${KUBECONFIG} --bind-address=127.0.0.1:8080

install:
	cd cmd/alerterator && go install

test:
	go test ./... --coverprofile=cover.out

codegen-crd:
	${ROOT_DIR}/hack/update-codegen.sh

codegen-updater:
	go generate ${ROOT_DIR}/hack/generator/updater.go | goimports > ${ROOT_DIR}/updater/zz_generated.go
