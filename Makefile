APP_NAME := cleaner
BIN_DIR := bin
CMD := ./cmd/cleaner
GOCACHE_DIR := .gocache
PREFIX ?= /usr/local
USER_BIN_DIR ?= $(HOME)/bin
INSTALL_PATH ?= $(PREFIX)/bin/$(APP_NAME)

.PHONY: fmt test run build install install-user uninstall clean

fmt:
	gofmt -w ./cmd ./internal

test:
	mkdir -p $(GOCACHE_DIR)
	GOCACHE=$(CURDIR)/$(GOCACHE_DIR) go test ./...

run:
	mkdir -p $(GOCACHE_DIR)
	GOCACHE=$(CURDIR)/$(GOCACHE_DIR) go run $(CMD) $(ARGS)

build:
	mkdir -p $(BIN_DIR) $(GOCACHE_DIR)
	GOCACHE=$(CURDIR)/$(GOCACHE_DIR) go build -o $(BIN_DIR)/$(APP_NAME) $(CMD)

install: build
	install -d "$(PREFIX)/bin"
	install -m 0755 "$(BIN_DIR)/$(APP_NAME)" "$(INSTALL_PATH)"

install-user: build
	install -d "$(USER_BIN_DIR)"
	install -m 0755 "$(BIN_DIR)/$(APP_NAME)" "$(USER_BIN_DIR)/$(APP_NAME)"

uninstall:
	rm -f "$(INSTALL_PATH)" "$(USER_BIN_DIR)/$(APP_NAME)"

clean:
	rm -rf $(BIN_DIR) $(GOCACHE_DIR)
