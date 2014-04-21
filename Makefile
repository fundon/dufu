TEST=.
BENCH=.

get:
	@go get -d ./...

fmt:
	@go fmt ./...

build:	get
	@mkdir -p bin
	@go build -a -o bin/space

.PHONY: bench fmt get build
