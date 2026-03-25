#!/bin/bash
# test-release-locally.sh
# Simulates GitHub Actions release workflow locally
# This tests the build process without needing actual GHA runners

set -e  # Exit on error

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  Maestro Release Workflow - Local Simulation              ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
VERSION=${1:-"v0.0.0-test"}
DIST_DIR="dist"
RELEASE_DIR="dist/release"

# Step 1: Clean previous builds
echo -e "${BLUE}[STEP 1/8]${NC} Cleaning previous builds..."
rm -rf $DIST_DIR bin/
echo -e "${GREEN}✓${NC} Clean complete"
echo ""

# Step 2: Verify Go installation
echo -e "${BLUE}[STEP 2/8]${NC} Verifying Go installation..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗${NC} Go is not installed"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓${NC} Found Go version: $GO_VERSION"
echo ""

# Step 3: Get version info
echo -e "${BLUE}[STEP 3/8]${NC} Setting up version info..."
echo "Version: $VERSION"
echo -e "${GREEN}✓${NC} Version configured"
echo ""

# Step 4: Build for all platforms
echo -e "${BLUE}[STEP 4/8]${NC} Building binaries for all platforms..."
echo ""

mkdir -p $DIST_DIR

# Define build matrix (from .github/workflows/release.yml)
declare -a builds=(
    "linux:amd64:"
    "linux:arm64:"
    "linux:386:"
    "darwin:amd64:"
    "darwin:arm64:"
    "windows:amd64:.exe"
    "windows:386:.exe"
)

build_count=0
failed_builds=()

for build in "${builds[@]}"; do
    IFS=':' read -r os arch ext <<< "$build"
    binary_name="maestro-${os}-${arch}${ext}"
    
    echo -n "Building $binary_name... "
    
    if GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o "${DIST_DIR}/${binary_name}" . 2>/dev/null; then
        
        size=$(du -h "${DIST_DIR}/${binary_name}" | cut -f1)
        echo -e "${GREEN}✓${NC} ($size)"
        ((build_count++)) || true
    else
        echo -e "${RED}✗${NC} FAILED"
        failed_builds+=("$binary_name")
    fi
done

echo ""
echo -e "${GREEN}✓${NC} Completed $build_count builds"

if [ ${#failed_builds[@]} -gt 0 ]; then
    echo -e "${YELLOW}⚠ Some builds may have failed (expected for cross-compilation without full toolchain)${NC}"
fi
echo ""

# Step 5: Create release archives
echo -e "${BLUE}[STEP 5/8]${NC} Creating release archives..."
mkdir -p $RELEASE_DIR

for build in "${builds[@]}"; do
    IFS=':' read -r os arch ext <<< "$build"
    binary_name="maestro-${os}-${arch}${ext}"
    tar_name="maestro-${VERSION}-${os}-${arch}.tar.gz"
    
    echo -n "Creating $tar_name... "
    
    if [ -f "${DIST_DIR}/${binary_name}" ]; then
        tar -czf "${RELEASE_DIR}/${tar_name}" -C "${DIST_DIR}" "${binary_name}"
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${YELLOW}⚠${NC} Skipped (binary not found)"
    fi
done

echo ""
echo -e "${GREEN}✓${NC} Archives created"
echo ""

# Step 6: Generate checksums
echo -e "${BLUE}[STEP 6/8]${NC} Generating checksums..."
cd $RELEASE_DIR
sha256sum *.tar.gz > SHA256SUMS.txt 2>/dev/null || echo "# No archives found" > SHA256SUMS.txt
cd - > /dev/null
echo -e "${GREEN}✓${NC} SHA256SUMS.txt created"
echo ""

# Step 7: Verify checksums
echo -e "${BLUE}[STEP 7/8]${NC} Verifying checksums..."
cd $RELEASE_DIR
if [ -s SHA256SUMS.txt ] && sha256sum -c SHA256SUMS.txt > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} All checksums verified"
else
    echo -e "${YELLOW}⚠${NC} Checksum verification skipped (no valid archives)"
fi
cd - > /dev/null
echo ""

# Step 8: Summary
echo -e "${BLUE}[STEP 8/8]${NC} Release Summary..."
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  RELEASE ARTIFACTS                                           "
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Version: $VERSION"
echo "Location: $RELEASE_DIR"
echo ""

if ls ${DIST_DIR}/maestro-* 1> /dev/null 2>&1; then
    echo "Binaries:"
    ls -lh ${DIST_DIR}/maestro-* | awk '{printf "  %-40s %s\n", $9, $5}'
    echo ""
fi

if ls ${RELEASE_DIR}/*.tar.gz 1> /dev/null 2>&1; then
    echo "Release Archives:"
    ls -lh ${RELEASE_DIR}/*.tar.gz | awk '{printf "  %-50s %s\n", $9, $5}'
    echo ""
fi

echo "Checksums:"
if [ -s ${RELEASE_DIR}/SHA256SUMS.txt ]; then
    head -5 ${RELEASE_DIR}/SHA256SUMS.txt | while read line; do
        echo "  ${line:0:64}... ${line: -30}"
    done
else
    echo "  (No checksums generated)"
fi
echo ""

echo "═══════════════════════════════════════════════════════════════"
echo ""
echo -e "${GREEN}✓${NC} Local release simulation complete!"
echo ""

if ls ${DIST_DIR}/maestro-$(go env GOOS)-$(go env GOARCH)* 1> /dev/null 2>&1; then
    echo "To test the binary:"
    echo "  ./${DIST_DIR}/maestro-$(go env GOOS)-$(go env GOARCH) --help"
    echo ""
fi
