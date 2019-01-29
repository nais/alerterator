NAME       := alerterator
TAG        := navikt/${NAME}
LATEST     := ${TAG}:latest
ROOT_DIR   := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))


.PHONY: build docker docker-push local install test codegen-crd codegen-updater

build:
	cd cmd/alerterator && go build

docker:
	docker image build -t ${TAG}:$(shell /bin/cat ./version) -t ${TAG} -t ${NAME} -t ${LATEST} -f Dockerfile .

docker-push:
	docker image push ${TAG}:$(shell /bin/cat ./version)
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