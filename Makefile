# Makefile for Go Breakout Game - Desktop and WASM builds

# Variables
BINARY_NAME=breakout
WASM_NAME=breakout.wasm
DOCS_DIR=docs
GO_VERSION=$(shell go version | cut -d' ' -f3)

# Default target
.PHONY: all
all: build

# Build desktop version
.PHONY: build
build:
	go build -o $(BINARY_NAME) main.go

# Build WASM version
.PHONY: wasm
wasm: clean-wasm
	mkdir -p $(DOCS_DIR)
	GOOS=js GOARCH=wasm go build -o $(DOCS_DIR)/$(WASM_NAME) main.go
	@if [ -f "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" $(DOCS_DIR)/; \
	elif [ -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then \
		cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(DOCS_DIR)/; \
	else \
		echo "Error: wasm_exec.js not found in Go installation"; exit 1; \
	fi

# Create HTML file for WASM deployment
.PHONY: wasm-html
wasm-html: wasm
	@echo "Creating HTML file for WASM deployment..."
	@./scripts/create-html.sh $(DOCS_DIR)

# Full WASM build with HTML
.PHONY: wasm-full
wasm-full: wasm-html

# Clean desktop build
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Clean WASM build
.PHONY: clean-wasm
clean-wasm:
	rm -f $(DOCS_DIR)/$(WASM_NAME) $(DOCS_DIR)/wasm_exec.js $(DOCS_DIR)/index.html

# Clean all builds
.PHONY: clean-all
clean-all: clean clean-wasm
	rm -rf $(DOCS_DIR)

# Run desktop version
.PHONY: run
run: build
	./$(BINARY_NAME)

# Serve WASM version locally for testing
.PHONY: serve
serve: wasm-full
	@echo "Starting local server for WASM testing..."
	@echo "Open http://localhost:8080 in your browser"
	@which python3 > /dev/null && cd $(DOCS_DIR) && python3 -m http.server 8080 || \
	 which python > /dev/null && cd $(DOCS_DIR) && python -m SimpleHTTPServer 8080 || \
	 echo "Python not found. Please serve the docs/ directory with any HTTP server."

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build desktop version"
	@echo "  wasm       - Build WASM version only"
	@echo "  wasm-html  - Build WASM version with HTML"
	@echo "  wasm-full  - Complete WASM build (same as wasm-html)"
	@echo "  run        - Build and run desktop version"
	@echo "  serve      - Build WASM and serve locally for testing"
	@echo "  clean      - Clean desktop build"
	@echo "  clean-wasm - Clean WASM build"
	@echo "  clean-all  - Clean all builds"
	@echo "  help       - Show this help"