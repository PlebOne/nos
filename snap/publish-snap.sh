#!/bin/bash

# Publish script for nos snap package
# This script builds and publishes the snap to the Snap Store

set -e

echo "🚀 Publishing nos snap to Snap Store..."

# Check if snapcraft is installed
if ! command -v snapcraft &> /dev/null; then
    echo "❌ Error: snapcraft is not installed"
    echo "Install with: sudo snap install snapcraft --classic"
    exit 1
fi

# Check if logged in to snapcraft
if ! snapcraft whoami &> /dev/null; then
    echo "❌ Error: Not logged in to snapcraft"
    echo "Login with: snapcraft login"
    exit 1
fi

# Navigate to the project root
cd "$(dirname "$0")/.."

# Prompt for channel (stable, candidate, beta, edge)
echo "📋 Select release channel:"
echo "  1) edge (latest development)"
echo "  2) beta (pre-release testing)"
echo "  3) candidate (release candidate)"
echo "  4) stable (production release)"
read -p "Enter choice (1-4) [default: 1]: " choice

case $choice in
    1|"") CHANNEL="edge" ;;
    2) CHANNEL="beta" ;;
    3) CHANNEL="candidate" ;;
    4) CHANNEL="stable" ;;
    *) echo "❌ Invalid choice"; exit 1 ;;
esac

echo "📦 Building and publishing to $CHANNEL channel..."

# Clean any previous builds
echo "🧹 Cleaning previous builds..."
snapcraft clean --destructive-mode || true

# Build the snap
echo "📦 Building snap package..."
snapcraft --destructive-mode

# Find the generated snap file
SNAP_FILE=$(ls -t *.snap 2>/dev/null | head -n1)

if [ -z "$SNAP_FILE" ]; then
    echo "❌ Error: No snap file found after build"
    exit 1
fi

echo "📤 Uploading $SNAP_FILE to $CHANNEL channel..."

# Upload and release
snapcraft upload "$SNAP_FILE" --release="$CHANNEL"

echo "✅ Successfully published $SNAP_FILE to $CHANNEL channel!"
echo ""
echo "📋 The snap is now available via:"
if [ "$CHANNEL" = "stable" ]; then
    echo "   sudo snap install nos"
else
    echo "   sudo snap install nos --$CHANNEL"
fi
echo ""
echo "📋 Monitor the release at:"
echo "   https://snapcraft.io/nos"