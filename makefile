BACKEND_DIR := backend
SCRIPTS_DIR := scripts
BINARY_NAME := openyurt-backend

.PHONY: all build run clean

all: build

build:
	@echo "Building backend..."
	@cd $(BACKEND_DIR) && go build -o ../$(BINARY_NAME) main.go

run: build
	@echo "Running backend..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@go clean
