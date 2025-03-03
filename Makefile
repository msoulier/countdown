.PHONY: help fmt build

build: countdown

help:
	@echo "Targets:"
	@echo "    help fmt countdown build (build everything)"

fmt:
	goimports -w *.go
	gofmt -w *.go

countdown: countdown.go
	go build -o countdown countdown.go
