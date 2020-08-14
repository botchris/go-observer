.PHONY: all deps gometalinter test cover

all: linters test cover

deps:
	GOBIN=$(CURDIR)/bin  go install github.com/golangci/golangci-lint/cmd/golangci-lint

linters:
	bin/golangci-lint run -v

test:
	go test -v -race -cpu=1,2,4 -coverprofile=coverage.txt -covermode=atomic

cover:
	go tool cover -html=coverage.txt -o coverage.html
