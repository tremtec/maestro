#!/bin/bash
# validate-workflow.sh
# Validates GitHub Actions workflow files

echo "Validating GitHub Actions Workflow..."
echo ""

# Check if workflow file exists
if [ ! -f ".github/workflows/release.yml" ]; then
    echo "✗ Workflow file not found: .github/workflows/release.yml"
    exit 1
fi

echo "✓ Workflow file exists"

# Basic YAML syntax check
if command -v yq &> /dev/null; then
    echo "✓ yq found, validating YAML syntax..."
    if yq eval '.name' .github/workflows/release.yml > /dev/null 2>&1; then
        echo "✓ YAML syntax is valid"
    else
        echo "✗ YAML syntax error"
        exit 1
    fi
else
    echo "⚠ yq not installed, skipping YAML validation"
fi

# Check for required secrets
echo ""
echo "Checking workflow configuration..."

# Verify key sections exist
if grep -q "permissions:" .github/workflows/release.yml; then
    echo "✓ Permissions section found"
else
    echo "⚠ Permissions section not found (may cause issues)"
fi

if grep -q "contents: write" .github/workflows/release.yml; then
    echo "✓ Write permissions configured"
else
    echo "⚠ Write permissions not explicitly set"
fi

if grep -q "strategy:" .github/workflows/release.yml; then
    echo "✓ Matrix build strategy configured"
fi

if grep -q "build:" .github/workflows/release.yml; then
    echo "✓ Build job configured"
fi

if grep -q "release:" .github/workflows/release.yml; then
    echo "✓ Release job configured"
fi

if grep -q "softprops/action-gh-release" .github/workflows/release.yml; then
    echo "✓ Release action configured"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  Validation Summary                                          "
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Workflow: .github/workflows/release.yml"
echo ""
echo "Triggers:"
grep -A 2 "on:" .github/workflows/release.yml | grep -E "(push|tags)" | head -4
echo ""
echo "Jobs:"
grep "^[[:space:]]*[a-z-]*:" .github/workflows/release.yml | grep -v "^#" | sed 's/^/  - /'
echo ""
echo "✓ Workflow validation complete"
echo ""
echo "To test locally:"
echo "  1. make test-release"
echo "  2. Or run: ./scripts/test-release-locally.sh"
echo ""
echo "To trigger actual release:"
echo "  git tag v0.1.0"
echo "  git push origin v0.1.0"
echo ""
