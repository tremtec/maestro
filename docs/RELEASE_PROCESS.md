# Release Process Documentation

## Overview

This document describes how to release new versions of Maestro with built Go binaries for multiple platforms.

## Release Methods

### Method 1: GitHub Actions (Recommended)

Automatically triggered when you push a tag.

```bash
# 1. Update version in relevant files (if needed)
# Edit version constants in code if applicable

# 2. Commit changes
git add .
git commit -m "chore: prepare for v0.1.0 release"

# 3. Create and push tag
git tag v0.1.0
git push origin v0.1.0

# For pre-releases (alpha, beta, rc)
git tag v0.1.0-rc.1
git push origin v0.1.0-rc.1
```

The GitHub Actions workflow will:
1. Build binaries for all platforms (Linux, macOS, Windows)
2. Create tar.gz archives with checksums
3. Create a GitHub Release with release notes
4. Attach all assets to the release

### Method 2: Local Build (Testing)

Use the Makefile for local testing:

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Create full release locally
make release-local

# The release artifacts will be in dist/release/
```

### Method 3: GoReleaser (Advanced)

For advanced releases with additional features (Homebrew, Docker, etc.):

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Test the release (dry run)
goreleaser release --snapshot --clean

# Actual release (requires GITHUB_TOKEN)
export GITHUB_TOKEN=your_token
goreleaser release --clean
```

## Versioning

We follow [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH` (e.g., `v1.2.3`)
- Pre-releases: `v1.2.3-alpha.1`, `v1.2.3-beta.2`, `v1.2.3-rc.1`

### Tag Format

- Stable: `v0.1.0`, `v1.0.0`
- Pre-release: `v0.1.0-rc.1`, `v0.2.0-beta.1`

## Supported Platforms

| OS | Architecture | Binary Name |
|----|--------------|-------------|
| Linux | amd64 | maestro-linux-amd64 |
| Linux | arm64 | maestro-linux-arm64 |
| Linux | 386 | maestro-linux-386 |
| macOS | amd64 | maestro-darwin-amd64 |
| macOS | arm64 | maestro-darwin-arm64 |
| Windows | amd64 | maestro-windows-amd64.exe |
| Windows | 386 | maestro-windows-386.exe |

## Release Artifacts

Each release includes:

1. **Platform-specific tar.gz archives**
   - Contains the binary + README + LICENSE
   - Named: `maestro-{version}-{os}-{arch}.tar.gz`

2. **SHA256SUMS.txt**
   - Checksums for all archives
   - Use to verify downloads: `sha256sum -c SHA256SUMS.txt`

3. **Release Notes**
   - Auto-generated from git commits
   - Includes installation instructions
   - Platform compatibility matrix

## Installation from Release

### macOS

```bash
# Intel Macs
curl -L -o maestro.tar.gz "https://github.com/tremtec/maestro/releases/download/v0.1.0/maestro-v0.1.0-darwin-amd64.tar.gz"

# Apple Silicon Macs
curl -L -o maestro.tar.gz "https://github.com/tremtec/maestro/releases/download/v0.1.0/maestro-v0.1.0-darwin-arm64.tar.gz"

tar -xzf maestro.tar.gz
mv maestro /usr/local/bin/
maestro --version
```

### Linux

```bash
# AMD64
curl -L -o maestro.tar.gz "https://github.com/tremtec/maestro/releases/download/v0.1.0/maestro-v0.1.0-linux-amd64.tar.gz"

# ARM64 (AWS Graviton, Raspberry Pi, etc.)
curl -L -o maestro.tar.gz "https://github.com/tremtec/maestro/releases/download/v0.1.0/maestro-v0.1.0-linux-arm64.tar.gz"

tar -xzf maestro.tar.gz
sudo mv maestro /usr/local/bin/
maestro --version
```

### Windows

```powershell
# Download the .tar.gz or .zip from releases page
# Extract to a directory in your PATH
# Or use scoop/chocolatey (if added)
```

### Go Install

```bash
# Install specific version
go install github.com/tremtec/maestro@v0.1.0

# Install latest
go install github.com/tremtec/maestro@latest
```

## Verification

### Checksum Verification

```bash
# Download the binary and checksums file
curl -L -o maestro.tar.gz "https://github.com/tremtec/maestro/releases/download/v0.1.0/maestro-v0.1.0-linux-amd64.tar.gz"
curl -L -o SHA256SUMS.txt "https://github.com/tremtec/maestro/releases/download/v0.1.0/SHA256SUMS.txt"

# Verify
sha256sum -c SHA256SUMS.txt

# Should output:
# maestro-v0.1.0-linux-amd64.tar.gz: OK
```

### Binary Verification

```bash
# Check version
maestro --version

# Check help
maestro --help
```

## Post-Release Checklist

After releasing:

- [ ] Verify GitHub Release was created
- [ ] Check all platform binaries are attached
- [ ] Verify checksums file is present
- [ ] Test installation on at least one platform
- [ ] Update documentation if needed
- [ ] Announce in relevant channels (if applicable)
- [ ] Close related milestone/issues

## Troubleshooting

### Build Failures

**Issue:** Build fails on specific platform

**Solution:** Check that CGO is disabled (CGO_ENABLED=0) and all dependencies support the target platform.

### Missing Assets

**Issue:** Release created but no binaries attached

**Solution:** Check GitHub Actions logs for upload errors. Ensure `permissions: contents: write` is set.

### Checksum Mismatch

**Issue:** Checksums don't match

**Solution:** Re-download both files. Check for partial downloads or network issues.

## CI/CD Integration

### Example: Testing Before Release

Add to your CI workflow:

```yaml
name: Pre-release Checks

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24.2'
      
      - run: make ci
      
      - run: make build-all
        # Ensures cross-compilation works
```

### Automated Releases

For fully automated releases on merge to main:

```yaml
name: Auto Release

on:
  push:
    branches: [main]

jobs:
  release:
    if: "contains(github.event.head_commit.message, 'release:')"
    uses: ./.github/workflows/release.yml
```

## Resources

- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Releases Documentation](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases)
- [Semantic Versioning](https://semver.org/)
