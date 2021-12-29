PROJECT := gitlab-automation
SHELL := /bin/bash
PKG := github.com/ostaquet/gitlab-automation/modules/
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

help:
	@echo "*** Help..."
	@echo " - lint     : Proceed to code analysis"
	@echo " - fmt      : Format the code according to standards"
	@echo " - test     : Run tests"
	@echo " - snyk     : Run Snyk dependencies check"
	@echo " - coverage : run tests with coverage report (only if all tests passed)"
	@echo " - clean    : Clean the project"
	@echo " - dep      : Update dependencies"
	@echo " - build    : Build the executable"
	
lint:
	@echo "*** Code analysis..."
	@go vet $(shell go list ${PKG}/... | grep -v /vendor/)

fmt:
	@echo "*** Clean code format..."
	@go fmt $(shell go list ${PKG}/... | grep -v /vendor/)

test:
	@echo "*** Run tests..."
	@go test -short $(shell go list ${PKG}/... | grep -v /vendor/)

snyk:
	@echo "*** Run Snyk dependencies checks..."
	@snyk test

clean:
	@echo "*** Clean..."
	@rm -rf vendor/*
	@rm -rf build/*
	@rm -rf $(shell find ./ -name error -type d)
	@rm -rf $(shell find ./ -name generated -type d)
	@rm -f $(shell find ./ -name main -type f)
	@rm -f $(shell find ./ -name cov.out -type f)

coverage:
	@echo "*** Run tests with coverage report..."
	@go test $(shell go list ${PKG}/... | grep -v /vendor/)
	@go test -coverprofile=/tmp/cov.out $(shell go list ${PKG}/... | grep -v /vendor/)
	@go tool cover -func=/tmp/cov.out

dep:
	@echo "*** Update dependencies..."
	@go mod download

build: dep
	@echo "*** Building process..."
	@go build -o build/
	@echo "Executables in build/ directory"
