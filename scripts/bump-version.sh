#!/bin/bash

# Version bump script for MongoDB Replica Set project
# Usage: ./scripts/bump-version.sh [major|minor|patch]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if version type is provided
if [ $# -eq 0 ]; then
    print_error "Version type not specified. Usage: $0 [major|minor|patch]"
    exit 1
fi

VERSION_TYPE=$1

# Validate version type
if [[ ! "$VERSION_TYPE" =~ ^(major|minor|patch)$ ]]; then
    print_error "Invalid version type. Must be major, minor, or patch"
    exit 1
fi

# Read current version
if [ ! -f "VERSION" ]; then
    print_error "VERSION file not found"
    exit 1
fi

CURRENT_VERSION=$(cat VERSION)
print_status "Current version: $CURRENT_VERSION"

# Parse version components
IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

# Bump version based on type
case $VERSION_TYPE in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"
print_status "New version: $NEW_VERSION"

# Update VERSION file
echo "$NEW_VERSION" > VERSION
print_status "Updated VERSION file"

# Update package.json files if they exist
if [ -f "app-node/package.json" ]; then
    # Update Node.js app version
    sed -i "s/\"version\": \"[^\"]*\"/\"version\": \"$NEW_VERSION\"/" app-node/package.json
    print_status "Updated app-node/package.json version"
fi

# Create git tag
if git rev-parse --git-dir > /dev/null 2>&1; then
    # Check if we have uncommitted changes
    if [ -n "$(git status --porcelain)" ]; then
        print_warning "You have uncommitted changes. Please commit them before creating a tag."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
    
    # Create and push tag
    git add VERSION
    if [ -f "app-node/package.json" ]; then
        git add app-node/package.json
    fi
    
    git commit -m "Bump version to $NEW_VERSION"
    git tag -a "v$NEW_VERSION" -m "Release version $NEW_VERSION"
    
    print_status "Created git tag v$NEW_VERSION"
    print_status "To push the tag, run: git push origin v$NEW_VERSION"
else
    print_warning "Not in a git repository. Skipping git operations."
fi

print_status "Version bump completed successfully!"
print_status "New version: $NEW_VERSION" 