.PHONY: help fmt build

build: count

help:
	@echo "Targets:"
	@echo "    help fmt count build (build everything)"

fmt:
	goimports -w *.go
	gofmt -w *.go

count: count.go
	go build -o count count.go
