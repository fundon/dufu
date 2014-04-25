TEST=.
BENCH=.

get:
	@go get -d ./...

fmt:
	@go fmt ./...

build:	get
	@mkdir -p bin
	@go build -a -o bin/dufu

gox-build: get
	@mkdir -p bin
	@gox  -output bin/"dufu_{{.OS}}_{{.Arch}}"

.PHONY: bench fmt get build gox-build
