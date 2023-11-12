OS := $(shell uname -s)
BINARY_NAME = todogo

ifeq ($(OS),Windows_NT)
	# BINARY_NAME = goapi.exe
	RM = powershell Remove-Item -Path
else
	# BINARY_NAME = goapi
	RM = rm -rf
endif

build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME)
	@echo "Done."

b: build  # Alias for build

run: build
	@echo "Running..."
	@./bin/$(BINARY_NAME)

test:
	@go test -v ./...

clean:
	@echo "Cleaning up..."
	@go clean
	@$(RM) bin
	@echo "Done."