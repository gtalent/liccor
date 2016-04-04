all: test build

test: test-src test-cli

test-src:
	@echo "TEST LICCOR SRC"
	@go test -v -cover lib/*.go
	@echo "===\n"

test-cli:
	@echo "TEST LICCOR TOOL"
	@go test -v -cover
	@echo "===\n"

build:
	@echo "BUILD LICCOR TOOL"
	@go build

lint:
	@golint ./...

clean:
	@echo "CLEANUP LICCOR DIRECTORY"
	@go clean
