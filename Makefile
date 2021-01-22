FILE_HASH := $(or ${hash},${hash},"empty_hash")
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

build:
	@echo "-- building binary. buildHash ${FILE_HASH}"
	go build -ldflags "-X main.buildHash=${FILE_HASH} -X main.buildTime=${BUILD_TIME}" -o ./bin/collector ./cmd

format:
	@echo "-- format code"
	gofmt -s -w .

lint: format
	@echo "-- linter running"
	golangci-lint run -c .golangci.yaml ./pkg...
	golangci-lint run -c .golangci.yaml ./cmd...