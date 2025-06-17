.PHONY: test test-all clean

# Test single PHP version (specified by PHP_VERSION environment variable)
test:
	go test -v ./...

# Test all PHP versions sequentially
test-all:
	@echo "Testing all PHP versions..."
	@if [ -f .php-versions ]; then \
		php_versions=$$(grep -v '^#' .php-versions | grep -v '^$$'); \
	else \
		echo "Error: .php-versions file not found"; \
		exit 1; \
	fi; \
	for php_version in $$php_versions; do \
		echo "Testing PHP version: $$php_version"; \
		PHP_VERSION=$$php_version go test ./... || exit 1; \
	done
	@echo "All tests passed!"

# Build release package
build:
	@# Get version from command line argument or default to v1.0.0
	@VERSION=$${VERSION:-v1.0.0}; \
	RELEASE="phpwebp-concept-$$VERSION"; \
	WORKING_DIR="built/$$RELEASE"; \
	BUILT_ZIP="built/$$RELEASE.zip"; \
	echo "Building $$RELEASE..."; \
	rm -rf "$$WORKING_DIR" "$$BUILT_ZIP"; \
	mkdir -p "$$WORKING_DIR"; \
	cp -a wwwroot/phpwebp-concept "$$WORKING_DIR/phpwebp-concept"; \
	cp -a wwwroot/.htaccess "$$WORKING_DIR/htaccess-example.txt"; \
	cd built && zip -r "$$RELEASE.zip" "$$RELEASE"; \
	rm -rf "$$WORKING_DIR"; \
	echo "Built $$BUILT_ZIP"

# Clean up
clean:
	go clean -testcache
	docker ps -a | grep 'php.*apache' | awk '{print $$1}' | xargs -r docker rm -f || true
	rm -rf built/