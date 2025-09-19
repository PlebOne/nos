#!/bin/bash

# Build script for nos snap package
# This script builds the snap locally for testing

set -e

echo "🔨 Building nos snap package..."

# Check if snapcraft is installed
if ! command -v snapcraft &> /dev/null; then
    echo "❌ Error: snapcraft is not installed"
    echo "Install with: sudo snap install snapcraft --classic"
    exit 1
fi

# Navigate to the project root
cd "$(dirname "$0")/.."

# Clean any previous builds
echo "🧹 Cleaning previous builds..."
snapcraft clean --destructive-mode || true

# Build the snap
echo "📦 Building snap package..."
snapcraft --destructive-mode

# Find the generated snap file
SNAP_FILE=$(ls -t *.snap 2>/dev/null | head -n1)

if [ -n "$SNAP_FILE" ]; then
    echo "✅ Successfully built: $SNAP_FILE"
    echo ""
    echo "📋 To install locally for testing:"
    echo "   sudo snap install --dangerous ./$SNAP_FILE"
    echo ""
    echo "📋 To test the installation:"
    echo "   nos --help"
    echo ""
    echo "📋 To remove the test installation:"
    echo "   sudo snap remove nos"
else
    echo "❌ Error: No snap file found after build"
    exit 1
fi