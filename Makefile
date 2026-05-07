BINARY_NAME ?= devclean
PKG ?= ./...
FORMULA_PATH ?= homebrew-tap/Formula/devclean-cli.rb

.PHONY: fmt test build run install-user uninstall-user brew-formula-update brew-formula-publish release-and-publish brew-install-local sha256-url

fmt:
	go fmt $(PKG)

test:
	go test $(PKG)

build:
	go build -ldflags "-X github.com/wangweicheng7/devclean-cli/internal/version.Version=dev -X github.com/wangweicheng7/devclean-cli/internal/version.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo none) -X github.com/wangweicheng7/devclean-cli/internal/version.Date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/$(BINARY_NAME) ./cmd/devclean

run:
	go run ./cmd/devclean $(ARGS)

install-user:
	mkdir -p "$(HOME)/bin"
	go build -o "$(HOME)/bin/$(BINARY_NAME)" ./cmd/devclean
	@echo "installed to $(HOME)/bin/$(BINARY_NAME)"

uninstall-user:
	rm -f "$(HOME)/bin/$(BINARY_NAME)"
	@echo "removed $(HOME)/bin/$(BINARY_NAME)"

sha256-url:
	@URL="$${URL:-}"; \
	if [ -z "$$URL" ]; then echo "usage: make sha256-url URL=https://..." >&2; exit 2; fi; \
	tmp="$$(mktemp -t devclean-sha.XXXXXX)"; \
	curl -L -o "$$tmp" "$$URL"; \
	shasum -a 256 "$$tmp" | awk '{print $$1}'; \
	rm -f "$$tmp"

brew-formula-update:
	@if [ -z "$${TAG:-}" ]; then echo "usage: make brew-formula-update TAG=v0.1.0" >&2; exit 2; fi
	bash scripts/update_formula_from_tag.sh "$${TAG}"

brew-formula-publish:
	bash scripts/publish_formula_to_tap.sh

release-and-publish:
	@if [ -z "$${TAG:-}" ]; then echo "usage: make release-and-publish TAG=v0.2.1" >&2; exit 2; fi
	bash scripts/release_and_publish.sh "$${TAG}"

brew-install-local:
	@echo "brew install from local formula file: $(FORMULA_PATH)" >&2
	brew install --formula --build-from-source "$(FORMULA_PATH)"

