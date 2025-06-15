#!/bin/bash

# Update json-joy reference to latest version

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"
JSONJOY_DIR="$PROJECT_ROOT/.reference/json-joy"

echo "Updating json-joy reference..."

# Check if submodule exists
if [ ! -e "$JSONJOY_DIR/.git" ]; then
    echo "Error: json-joy submodule not found. Run: git submodule update --init --recursive"
    exit 1
fi

cd "$JSONJOY_DIR"

# Get current version
CURRENT_COMMIT=$(git rev-parse HEAD)
CURRENT_TAG=$(git describe --tags --exact-match 2>/dev/null || echo "latest")

echo "Current: $CURRENT_TAG ($CURRENT_COMMIT)"

# Fetch and check for updates
git fetch origin
LATEST_COMMIT=$(git rev-parse origin/master)

if [ "$CURRENT_COMMIT" = "$LATEST_COMMIT" ]; then
    echo "Already up to date!"
    exit 0
fi

echo "New changes available!"

# Update to latest
git pull origin master

# Get new version
NEW_COMMIT=$(git rev-parse HEAD)
NEW_TAG=$(git describe --tags --exact-match 2>/dev/null || echo "latest")

echo "Updated to: $NEW_TAG ($NEW_COMMIT)"

# Return to project root and commit
cd "$PROJECT_ROOT"
git add .reference/json-joy

# Generate commit message
COMMIT_MSG="chore: update json-joy reference to $NEW_TAG"
git commit -m "$COMMIT_MSG"

echo "json-joy reference updated successfully!"
