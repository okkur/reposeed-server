SHELL := /bin/bash

TARGET := $(shell echo $${PWD\#\#*/})
.DEFAULT_GOAL: $(TARGET)

VERSION := $(shell cat ./VERSION)
BUILD := `git rev-parse HEAD`

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check run

all: check install

$(TARGET): $(SRC)
	@packr build -o $(TARGET)

build: $(TARGET)
	@true

clean:
	@rm -f $(TARGET)

install:
	@packr install .

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@test -z $(shell gofmt -l ./main.go | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./...); do golint $${d}; done
	@go tool vet ${SRC}

packr:
	@go get -u github.com/gobuffalo/packr/...

run: install
	@$(TARGET)

version:
	@echo $(VERSION)
