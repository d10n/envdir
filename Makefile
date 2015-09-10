GOCMD = go
BUILD_DATE := $(shell date --utc --iso-8601=seconds)
VERSION := $(shell git describe --always --dirty --match 'v[0-9]*')
COMMIT := $(shell git rev-parse HEAD)

.PHONY: dist

dist:
	$(GOCMD) build -ldflags=" \
		-X main.buildDate=$(BUILD_DATE) \
		-X main.buildVersion=$(VERSION) \
		-X main.buildCommit=$(COMMIT) \
		-s \
		"
