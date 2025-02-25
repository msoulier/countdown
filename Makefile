.PHONY: help fmt build

help:
	@echo "Targets:"
	@echo "    help fmt count"

fmt:
	goimports -w *.go
	gofmt -w *.go

count: count.go
	go build -o count count.go

build: count
