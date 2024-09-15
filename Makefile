binary_name = anydb.exe

## help: Show this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## tidy: Run go mod tidy and go fmt
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: Build the binary
.PHONY: build
build:
	go build -o ./${binary_name} .

## clean: Remove the binary
.PHONY: clean
clean:
	rm ./${binary_name}
