# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PHPWebP Concept is a PHP-based middleware for automatic WebP image conversion on Apache web servers. It provides bidirectional conversion between traditional formats (JPEG/PNG/GIF) and WebP based on browser capabilities.

## Version Management

### PHP Versions
- Supported PHP versions are managed in `.php-versions` file (required)
- The CI/CD workflow reads PHP versions from `.php-versions` automatically
- Comments (lines starting with #) and empty lines are ignored
- New PHP versions from https://hub.docker.com/_/php are added after passing tests
- PHP versions don't need to be stable releases - any version available on Docker Hub without beta/alpha/rc tags is acceptable
- To add/remove PHP versions, edit `.php-versions` file (one version per line)

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

# Test individual PHP versions
PHP_VERSION=7.4 go test -v ./...   # Test PHP 7.4
PHP_VERSION=8.1 go test -v ./...   # Test PHP 8.1
PHP_VERSION=8.2 go test -v ./...   # Test PHP 8.2
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
2. Create a new branch: `git checkout -b php-8.4` (use the actual PHP version)
3. Add the version to `.php-versions` file
4. Test the new version: `PHP_VERSION=8.4 go test -v ./...`
5. Fix any issues found during testing
6. Commit changes: `git commit -m "Add support for PHP 8.4"`
7. Push branch: `git push -u origin php-8.4`
8. Create pull request with title: "Add support for PHP 8.4"

### Updating libwebp
1. Check current version in `.libwebp-version`
2. Create a new branch: `git checkout -b libwebp-1.5.0` (use the actual libwebp version)
3. Download new version from https://developers.google.com/speed/webp/docs/precompiled
4. Extract binaries for each architecture:
   - Linux x86_64 → `wwwroot/phpwebp-concept/bin/linux-x86_64/`
   - Linux ARM64 → `wwwroot/phpwebp-concept/bin/linux-aarch64/`
5. **IMPORTANT**: Set execute permissions on all binaries:
   ```bash
   chmod +x wwwroot/phpwebp-concept/bin/linux-aarch64/*
   chmod +x wwwroot/phpwebp-concept/bin/linux-x86_64/*
   ```
6. Update `.libwebp-version` with new version number
7. Run full test suite: `make test-all` to test all PHP versions
8. Commit changes: `git commit -m "Update libwebp to 1.5.0"`
9. Push branch: `git push -u origin libwebp-1.5.0`
10. Create pull request with title: "Update libwebp to 1.5.0"

### Updating Both PHP Version and libwebp
When updating both PHP version and libwebp simultaneously:
1. Create a combined branch: `git checkout -b php-8.4-libwebp-1.5.0`
2. Follow steps from "Adding New PHP Versions" (steps 3-5)
3. Follow steps from "Updating libwebp" (steps 3-7)
4. Commit all changes: `git commit -m "Add PHP 8.4 support and update libwebp to 1.5.0"`
5. Push branch: `git push -u origin php-8.4-libwebp-1.5.0`
6. Create pull request with title: "Add PHP 8.4 support and update libwebp to 1.5.0"

### Adding New Architecture Support
1. Create directory: `wwwroot/phpwebp-concept/bin/[os-architecture]/`
2. Add compiled binaries: `cwebp`, `dwebp`, `gif2webp`, `webpinfo`
3. **IMPORTANT**: Make binaries executable: `chmod +x wwwroot/phpwebp-concept/bin/[os-architecture]/*`
4. Test on target architecture

## Important Notes

### Binary File Permissions
When updating any program files under `wwwroot/phpwebp-concept/bin/`, always ensure to set execute permissions:
```bash
chmod +x wwwroot/phpwebp-concept/bin/[architecture]/*
```
This is required for the binaries to function properly on the server.