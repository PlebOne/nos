#!/bin/bash

# nos debian package build script
VERSION="0.9"
ARCH="amd64"
PACKAGE_NAME="nos"
MAINTAINER="Tim Dev <tim@example.com>"
DESCRIPTION="Beautiful Nostr CLI client with interactive menu"

# Create build directory structure
BUILD_DIR="build/nos_${VERSION}_${ARCH}"
mkdir -p "${BUILD_DIR}/DEBIAN"
mkdir -p "${BUILD_DIR}/usr/local/bin"
mkdir -p "${BUILD_DIR}/usr/share/doc/nos"
mkdir -p "${BUILD_DIR}/usr/share/pixmaps"

# Build the binary
echo "Building nos binary..."
go build -ldflags="-s -w" -o "${BUILD_DIR}/usr/local/bin/nos"

# Copy documentation
cp README.md "${BUILD_DIR}/usr/share/doc/nos/"
cp nos.jpeg "${BUILD_DIR}/usr/share/pixmaps/nos.png"

# Create control file
cat > "${BUILD_DIR}/DEBIAN/control" << EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Description: ${DESCRIPTION}
 nos is a beautiful command-line client for posting to Nostr,
 featuring an interactive menu-driven interface, secure key storage,
 multi-relay support, and a clean UI built with Charm.sh.
Depends: libc6
Homepage: https://github.com/timdev/nos
EOF

# Create copyright file
cat > "${BUILD_DIR}/usr/share/doc/nos/copyright" << EOF
Format: https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/
Upstream-Name: nos
Upstream-Contact: Tim Dev <tim@example.com>
Source: https://github.com/timdev/nos

Files: *
Copyright: 2024 Tim Dev
License: MIT
 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:
 .
 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.
 .
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 SOFTWARE.
EOF

# Create changelog
cat > "${BUILD_DIR}/usr/share/doc/nos/changelog" << EOF
nos (${VERSION}) stable; urgency=medium

  * Initial release
  * Interactive menu-driven interface
  * Secure key storage using system keyring
  * Multi-relay posting support
  * Relay management features
  * Post verification
  * Beautiful UI with Charm.sh

 -- ${MAINTAINER}  $(date -R)
EOF

gzip -9 "${BUILD_DIR}/usr/share/doc/nos/changelog"

# Set permissions
chmod 755 "${BUILD_DIR}/usr/local/bin/nos"
chmod 644 "${BUILD_DIR}/usr/share/doc/nos/"*
chmod 644 "${BUILD_DIR}/usr/share/pixmaps/nos.png"

# Build the package
echo "Building debian package..."
dpkg-deb --build "${BUILD_DIR}"

# Move package to dist directory
mkdir -p dist
mv "build/nos_${VERSION}_${ARCH}.deb" "dist/"

echo "Package built: dist/nos_${VERSION}_${ARCH}.deb"
