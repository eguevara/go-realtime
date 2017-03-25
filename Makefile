# The binary to build (just the basename).
BIN := realtime

# This repo's root import path (under GOPATH).
PKG := github.com/eguevara/go-realtime

# Where to push the docker image.
REGISTRY ?= erick

# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)
BUILDTAGS=

.PHONY: all fmt vet lint build test install
.DEFAULT: default

all: build fmt lint test vet install

build:
	@echo "+ $@"
	@go build -tags "$(BUILDTAGS) cgo" .

fmt:
	@echo "+ $@"
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

lint:
	@echo "+ $@"
	@go list ./... | grep -v /vendor/ | xargs -L1 golint

test: fmt lint vet
	@echo "+ $@"
	@go test -v -tags "$(BUILDTAGS) cgo" $(shell go list ./... | grep -v vendor)

vet:
	@echo "+ $@"
	@go vet $(shell go list ./... | grep -v vendor)

install:
	@echo "+ $@"
	@go install .
