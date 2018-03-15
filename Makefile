GOCMD = go
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%S%z')
VERSION := $(shell git describe --always --dirty --match 'v[0-9]*')
COMMIT := $(shell git rev-parse HEAD)

.PHONY: all
all: test dist

test:
	$(GOCMD) test

dist:
	$(GOCMD) build -a -installsuffix cgo -ldflags=" \
		-X main.buildDate=$(BUILD_DATE) \
		-X main.buildVersion=$(VERSION) \
		-X main.buildCommit=$(COMMIT) \
		-s -w \
		"
