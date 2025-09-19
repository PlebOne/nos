#!/bin/bash

# Build .deb packages for all supported architectures
# Usage: ./build-all-debs.sh

set -e

echo "🔨 Building .deb packages for all architectures..."

# Clean previous builds
rm -rf build/ dist/

# Architectures to build
ARCHS=("amd64" "arm64" "armhf")

for arch in "${ARCHS[@]}"; do
    echo ""
    echo "📦 Building for architecture: $arch"
    ./build-deb.sh "$arch"
done

echo ""
echo "✅ All .deb packages built successfully!"
echo ""
echo "📋 Generated packages:"
ls -la dist/*.deb

echo ""
echo "📋 Package sizes:"
du -h dist/*.deb

echo ""
echo "🚀 To test a package:"
echo "   sudo dpkg -i dist/nos_*_amd64.deb"
echo "   nos --help"