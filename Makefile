default: fmt lint install generate

generate:
	cd tools
	go generate ./...

build: generate
	go build -v ./...

build_for_testing: generate
	go build -v -o ./ ./...

install: build
	go install -v ./...

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

.PHONY: fmt  build install generate
