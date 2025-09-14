# ==========================
# Swagger Docs
# ==========================
DOCS_DIR=./docs
SWAGGER_MAIN=./cmd/server/main.go
SWAGGER_OUT=$(DOCS_DIR)

.PHONY: swagger-gen
swagger-gen:
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "❌ 'swag' CLI not found. Install with:"; \
		echo "   go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi
	@mkdir -p $(DOCS_DIR)
	@echo "🔄 Generating Swagger documentation..."
	@swag init -g $(SWAGGER_MAIN) -o $(SWAGGER_OUT) --parseDependency --parseInternal
	@echo "✅ Swagger generated at $(DOCS_DIR)/swagger.json"

.PHONY: swagger-clean
swagger-clean:
	@rm -rf $(DOCS_DIR)
	@echo "🧹 Removed $(DOCS_DIR)"

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  swagger-gen       Generate Swagger documentation (requires swag CLI)"
	@echo "  swagger-clean     Remove generated Swagger docs"
