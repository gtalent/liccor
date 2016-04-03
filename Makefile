all: test build

test:
	@echo "TEST LICCOR SRC"
	@go test -v -cover lib/*.go
	@echo "TEST LICCOR TOOL"
	@go test -v -cover

build:
	@echo "BUILD LICCOR TOOL"
	@go build

clean:
	@echo "CLEANUP LICCOR DIRECTORY"
	@go clean
