# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PHPWebP Concept is a PHP-based middleware for automatic WebP image conversion on Apache web servers. It provides bidirectional conversion between traditional formats (JPEG/PNG/GIF) and WebP based on browser capabilities.

## Version Management

### PHP Versions
- Supported PHP versions are managed in `.github/workflows/cicd.yml`
- New PHP versions from https://hub.docker.com/_/php are added after passing tests
- Currently supported: 5.6, 7.0, 7.1, 7.2, 7.3, 7.4, 8.0, 8.1, 8.2

### libwebp Updates
- Check https://developers.google.com/speed/webp/docs/precompiled for new versions
- Current version is tracked in `.libwebp-version` (currently 1.5.0)
- Update process:
  1. Download new binaries for supported architectures
  2. Replace files in `wwwroot/phpwebp-concept/bin/[architecture]/`
  3. Update `.libwebp-version` file
  4. Test with all PHP versions
- Required binaries: `cwebp`, `dwebp`, `gif2webp`, `webpinfo`
- Supported architectures: `linux-x86_64`, `linux-aarch64`

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

## Development Guidelines

### Adding New PHP Versions
1. Check https://hub.docker.com/_/php for available PHP versions
2. Add the version to `.github/workflows/cicd.yml` matrix
3. Run tests locally: `PHP_VERSION=X.X make test`
4. Create PR if tests pass

### Updating libwebp
1. Check current version in `.libwebp-version`
2. Download new version from https://developers.google.com/speed/webp/docs/precompiled
3. Extract binaries for each architecture:
   - Linux x86_64 → `wwwroot/phpwebp-concept/bin/linux-x86_64/`
   - Linux ARM64 → `wwwroot/phpwebp-concept/bin/linux-aarch64/`
4. Update `.libwebp-version` with new version number
5. Run full test suite: `make test-all`

### Adding New Architecture Support
1. Create directory: `wwwroot/phpwebp-concept/bin/[os-architecture]/`
2. Add compiled binaries: `cwebp`, `dwebp`, `gif2webp`, `webpinfo`
3. Make binaries executable: `chmod +x wwwroot/phpwebp-concept/bin/[os-architecture]/*`
4. Test on target architecture