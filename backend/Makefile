# ──────────────────────────────────────────────────────────────
# Универсальный Makefile — работает из корня проекта И из backend/
# ──────────────────────────────────────────────────────────────

APP_NAME    := avito_test_app
BUILD_DIR   := ../bin               # ← именно так — всегда относительно backend/

# Определяем, где мы находимся
ifeq ($(wildcard go.mod),)
    # Мы в корне проекта → go.mod лежит в backend/
    BACKEND_DIR := .
else
    # Мы уже внутри backend/
    BACKEND_DIR := .
endif

ifeq ($(OS),Windows_NT)
    BINARY_NAME := $(BUILD_DIR)/$(APP_NAME).exe
    MKDIR       := mkdir
    RM          := del /Q 2>nul ||:
else
    BINARY_NAME := $(BUILD_DIR)/$(APP_NAME)
    MKDIR       := mkdir -p
    RM          := rm -f
endif

$(BUILD_DIR):
	@$(MKDIR) "$(BUILD_DIR)"

.PHONY: all build run test tidy clean help build-all

## all: tidy + test + build
all: tidy test build

## build: собрать бинарник для текущей ОС
build: $(BUILD_DIR)
	@echo "→ Building for $$(go env GOOS)/$$(go env GOARCH) → $(BINARY_NAME)"
	cd $(BACKEND_DIR) && go build -o "$(BINARY_NAME)" ./cmd
	@echo "Build complete"

## run: собрать и запустить
run: build
	$(BINARY_NAME)

## test: запустить тесты
test:
	cd $(BACKEND_DIR) && go test ./... -cover -v

## tidy: почистить зависимости
tidy:
	cd $(BACKEND_DIR) && go mod tidy

## clean: удалить бинарники
clean:
	$(RM) $(BUILD_DIR)/*

## build-linux: Linux amd64
build-linux: $(BUILD_DIR)
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o ../$(BUILD_DIR)/$(APP_NAME)-linux ./cmd

## build-windows: Windows amd64
build-windows: $(BUILD_DIR)
	cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 go build -o ../$(BUILD_DIR)/$(APP_NAME).exe ./cmd

## build-darwin: macOS Intel + Apple Silicon
build-darwin: $(BUILD_DIR)
	cd $(BACKEND_DIR) && GOOS=darwin GOARCH=amd64 go build -o ../$(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd
	cd $(BACKEND_DIR) && GOOS=darwin GOARCH=arm64 go build -o ../$(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./cmd

## build-all: собрать всё
build-all: build-linux build-windows build-darwin
	@echo "Cross-compilation complete"

## help: показать справку
help:
	@echo "Available targets:"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'