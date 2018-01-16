THIS_FILE := $(lastword $(MAKEFILE_LIST))

build:
	@go build -o bin/gs

.PHONY: build
