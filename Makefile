# Go parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_MOD=$(GO_CMD) mod
GO_CLEAN=$(GO_CMD) clean
GO_FMT=$(GO_CMD) fmt
GO_WASM=GOOS=js GOARCH=wasm $(GO_BUILD) -o assets/main.wasm  cmd/constructor.go cmd/wasm.go

# Project parameters
MAIN_SRC=cmd/main.go
WASM_SRC=cmd/wasm.go
ASSETS_DIR=assets
OUTPUT_DIR=dist

# Targets
.PHONY: all build wasm assets clean format deps test

all: clean format deps build wasm assets

build:
	@echo "Building the Go project for native..."
	$(GO_BUILD) -o $(OUTPUT_DIR)/main $(MAIN_SRC)

wasm:
	@echo "Building the Go project for WebAssembly..."
	mkdir -p $(OUTPUT_DIR)
	cp $(ASSETS_DIR)/wasm_exec.js $(OUTPUT_DIR)
	GOOS=js GOARCH=wasm $(GO_BUILD) -o $(ASSETS_DIR)/main.wasm cmd/constructor.go cmd/wasm.go

assets:
	@echo "Copying assets..."
	mkdir -p $(OUTPUT_DIR)
	cp -r $(ASSETS_DIR)/* $(OUTPUT_DIR)

clean:
	@echo "Cleaning up..."
	$(GO_CLEAN)
	rm -rf $(OUTPUT_DIR)

format:
	@echo "Formatting Go code..."
	$(GO_FMT) ./...

deps:
	@echo "Downloading dependencies..."
	$(GO_MOD) tidy

test:
	@echo "Running tests..."
	$(GO_CMD) test ./...

