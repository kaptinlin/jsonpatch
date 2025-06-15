#!/bin/bash

# Check json-joy API changes for Go implementation compatibility

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
JSONJOY_DIR="$PROJECT_ROOT/.reference/json-joy"

echo "Checking JSON Patch API changes..."

# Check if submodule exists
if [ ! -e "$JSONJOY_DIR/.git" ]; then
    echo "Error: json-joy submodule not found. Run: git submodule update --init --recursive"
    exit 1
fi

cd "$JSONJOY_DIR"

# Get version information
CURRENT_TAG=$(git describe --tags --exact-match 2>/dev/null || echo "latest")
CURRENT_COMMIT=$(git rev-parse HEAD)

echo "Current version: $CURRENT_TAG ($CURRENT_COMMIT)"

# Check comparison version
if [ $# -eq 1 ]; then
    COMPARE_VERSION="$1"
    echo "Comparing with: $COMPARE_VERSION"
else
    COMPARE_VERSION=$(git describe --tags --abbrev=0 HEAD~1 2>/dev/null || echo "")
    if [ -z "$COMPARE_VERSION" ]; then
        echo "No previous version found. Showing current API structure."
        COMPARE_VERSION="HEAD"
    else
        echo "Comparing with previous: $COMPARE_VERSION"
    fi
fi

echo ""
echo "=== JSON PATCH API ANALYSIS ==="

# Core exports
echo "Exports:"
grep -n "export" src/json-patch/index.ts 2>/dev/null || echo "No exports found"

echo ""
echo "Type definitions:"
grep -n "export.*interface\|export.*type" src/json-patch/types.ts 2>/dev/null || echo "No types found"

echo ""
echo "Core functions:"
grep -A 2 "export.*function.*applyPatch\|export.*function.*validate" src/json-patch/index.ts 2>/dev/null || echo "Functions not found"

# Show changes if comparison version exists
if [ "$COMPARE_VERSION" != "HEAD" ] && git rev-parse "$COMPARE_VERSION" >/dev/null 2>&1; then
    echo ""
    echo "=== CHANGES ANALYSIS ==="

    echo "Modified files:"
    git diff --name-only "$COMPARE_VERSION" HEAD -- src/json-patch/ | grep -E '\.(ts|js)$' || echo "No files changed"

    echo ""
    echo "API changes in core files:"

    CORE_FILES=(
        "src/json-patch/index.ts"
        "src/json-patch/types.ts"
        "src/json-patch/applyPatch/index.ts"
    )

    for file in "${CORE_FILES[@]}"; do
        if git diff --quiet "$COMPARE_VERSION" HEAD -- "$file"; then
            echo "  ✓ $file - No changes"
        else
            echo "  • $file - Modified"
            git diff "$COMPARE_VERSION" HEAD -- "$file" | grep -E '^[+-].*export|^[+-].*interface|^[+-].*function' | head -5
        fi
    done
fi

echo ""
echo "=== GO IMPLEMENTATION CHECKLIST ==="
echo "Items to verify:"
echo "  □ Operation struct matches TypeScript interface"
echo "  □ ApplyPatch function signature matches"
echo "  □ All operation types supported"
echo "  □ Error handling is consistent"
echo "  □ Type definitions are complete"

cd "$PROJECT_ROOT"
