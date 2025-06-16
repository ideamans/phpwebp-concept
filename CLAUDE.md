# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PHPWebP Concept is a PHP-based middleware for automatic WebP image conversion on Apache web servers. It provides bidirectional conversion between traditional formats (JPEG/PNG/GIF) and WebP based on browser capabilities.

## Common Development Commands

### Setup
```bash
yarn install          # Install Node.js dependencies
```

### Development Server
```bash
PHP=7.4 yarn dev     # Start development server with specific PHP version
```

### Testing
```bash
# Go tests (new)
make test            # Test with PHP version specified in PHP_VERSION env var
make test-all        # Test all PHP versions sequentially
PHP_VERSION=8.1 go test -v ./...  # Test specific PHP version

# JavaScript tests (legacy - to be removed)
yarn test:auto       # Run tests against all PHP versions
PHP=8.1 yarn test:auto  # Test specific PHP version
```

### Building
```bash
make build           # Create release package (default: v1.0.0)
VERSION=v1.0.1 make build  # Build with specific version
```

## Architecture

### Core Components

1. **PHP Application** (`/wwwroot/phpwebp-concept/`)
   - `compress.php`: Converts JPEG/PNG/GIF → WebP for supported browsers
   - `decompress.php`: Converts WebP → PNG for non-WebP browsers
   - `common.php`: Shared utilities for request parsing and WebP command execution
   - `bin/`: Platform-specific WebP binaries (cwebp, dwebp, gif2webp, webpinfo)

2. **Request Flow**
   - Apache `.htaccess` intercepts image requests
   - Routes to `compress.php` or `decompress.php` based on Accept headers
   - Converts images using WebP binaries via stdin (prevents command injection)
   - Caches converted images in system temp directory
   - Returns appropriate Content-Type headers

3. **Caching Strategy**
   - Cache key: SHA1 of file path + modification time + file size
   - Stores converted images in system temp directory
   - Zero-byte files mark conversion failures to avoid repeated attempts
   - Size comparison: only serves WebP if smaller than original

### Security Considerations

- Path validation prevents directory traversal attacks
- Binary execution uses stdin to avoid shell injection
- Proper error handling with appropriate HTTP status codes

## Testing

The project uses AVA framework for testing with Docker Compose to test across PHP versions 5.4 - 8.3. Tests verify:
- Format conversion functionality
- Browser compatibility handling
- Error conditions and security
- Performance optimizations

## CI/CD Notes for Go Development

### Windows Runner Issues

When running Go commands in GitHub Actions on Windows runners, avoid using `./...` pattern directly as it can be misinterpreted. Windows may parse `./...` as `.txt` or other extensions.

**Problem Example:**

```yaml
# This may fail on Windows
run: go test -c -o storemanager.test ./...
```

**Solutions:**

1. Use explicit package path:

```yaml
run: go test -c -o storemanager.test github.com/ideamans/lightfile6-batch-store-manager
```

2. Use `.` for current directory:

```yaml
run: go test -c -o storemanager.test .
```

3. Use PowerShell with proper escaping:

```yaml
shell: pwsh
run: go test -c -o storemanager.test .\...
```

4. Set working directory explicitly:

```yaml
working-directory: ${{ github.workspace }}
run: go test -c -o storemanager.test ./...
```

This is a common issue when developing cross-platform Go applications and should be considered when writing GitHub Actions workflows.