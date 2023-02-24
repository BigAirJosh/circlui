# Variables
export

SHELL := /bin/bash -o errexit -o nounset -o pipefail

MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

VERBOSE ?= false
ifeq (${VERBOSE}, false)
	# --silent drops the need to prepend `@` to suppress command output
	MAKEFLAGS += --silent
endif

USER_ID  := $(shell id -u)
GROUP_ID := $(shell id -g)

## Applications
DOCKER         ?= docker
DOCKER_COMPOSE ?= ${DOCKER} compose

# Helpers
.PHONY: all
all: depend build ## download dependencies and run a build

.PHONY: tinker
tinker: ## drop into built container for debugging purposes
	${DOCKER_COMPOSE} run --rm go /bin/sh

.PHONY: run
run: ## run local using docker container and mocked API
	${DOCKER_COMPOSE} run --rm go go run circlui.go

# Dependencies
.PHONY: depend
depend: ## build base image

# Linting
.PHONY: lint
lint: ## run all linters

# Building
.PHONY: build
build: depend ## build project

# Cleaning
.PHONY: clean
clean: ## clean build artifacts

.PHONY: clean-docker
clean-docker: ## clean docker compose images/containers
	${DOCKER_COMPOSE} down --rmi all -v
	${DOCKER_COMPOSE} rm --stop --force -v

.PHONY: clean-all
clean-all: clean clean-docker ## reset the project to a totally clean state

# Make
print-% :
	echo $* = $($*)

help:
	grep -Eh '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: help
