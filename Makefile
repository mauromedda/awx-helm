# This is a Makefile useful to wrap Help Chart test in CICD Pipelines
# It executes 3 tests:
#   - linting check
#   - syntax
#   - semantic (integration)
.PHONY: all
all: dprep lint syntax

SHELL=/usr/bin/env bash


CHART_NAME ?= $(shell basename `git rev-parse --show-toplevel`)
COMMIT_ID ?= $(shell git rev-parse --short HEAD)

HELM_VERSION ?= 2.14.1
DOCKER_IMAGE_HELM := helm-runner

dprep:
	@echo '===> Build the test container'
	docker build --build-arg VERSION=${HELM_VERSION} -t helm-runner -f Dockerfile.terratest .
	@echo '===> Go and helm container ready'

lint:
	@echo '===> [HELM] Start the linting test'
	@docker run -ti --rm --volume $(PWD):/apps/$(CHART_NAME):rw ${DOCKER_IMAGE_HELM} helm lint ${CHART_NAME}
	@echo '===> Lint test finished'

syntax:
	@echo '[TERRATEST|HELM] Start the syntax (unit) testing'
	@docker run -ti --rm --volume $(PWD):/apps/$(CHART_NAME):rw -w /apps/$(CHART_NAME)/tests ${DOCKER_IMAGE_HELM} ;\
		go test -mod=vendor -v -tags helm -run=TestHelmBasicTemplateRenderedDeployment .
	@echo '[TERRATEST|HELM] Unit tests completed with success!'

integration:
	@echo '[TERRATEST|HELM] Start the integration tests - TBD'