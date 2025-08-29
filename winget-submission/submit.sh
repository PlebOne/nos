#!/bin/bash

# Winget Package Submission Helper Script
# This script helps submit nos to the Microsoft winget repository

set -e

echo "🚀 nos Winget Package Submission Helper"
echo "======================================"
echo

# Check if winget-pkgs directory exists
if [ ! -d "../winget-pkgs" ]; then
    echo "📥 Cloning Microsoft winget-pkgs repository..."
    echo "Please run these commands first:"
    echo
    echo "1. Fork https://github.com/microsoft/winget-pkgs"
    echo "2. Clone your fork:"
    echo "   cd .."
    echo "   git clone https://github.com/YOUR_USERNAME/winget-pkgs.git"
    echo "   cd nos"
    echo
    echo "Then run this script again."
    exit 1
fi

echo "📁 Copying manifest files to winget-pkgs repository..."

# Create directory structure
mkdir -p ../winget-pkgs/manifests/p/PlebOne/nos/1.1.1

# Copy manifest files
cp winget-submission/manifests/p/PlebOne/nos/1.1.1/* ../winget-pkgs/manifests/p/PlebOne/nos/1.1.1/

echo "✅ Files copied successfully!"
echo

cd ../winget-pkgs

echo "🌿 Creating branch for submission..."
git checkout -b add-plebOne-nos-1.1.1

echo "📝 Adding manifest files..."
git add manifests/p/PlebOne/nos/1.1.1/

echo "💾 Committing changes..."
git commit -m "Add PlebOne.nos version 1.1.1

- Package: nos - A beautiful command-line client for posting to Nostr
- Publisher: PlebOne
- Version: 1.1.1
- License: MIT
- Architectures: x64, ARM64
- Installer Type: Portable (zip)

Release URL: https://github.com/PlebOne/nos/releases/tag/v1.1.1"

echo "🚀 Pushing to your fork..."
git push origin add-plebOne-nos-1.1.1

echo
echo "🎉 Success! Next steps:"
echo "========================"
echo "1. Go to your winget-pkgs fork on GitHub"
echo "2. Click 'Compare & pull request'"
echo "3. Use title: Add PlebOne.nos version 1.1.1"
echo "4. Copy the PR description from winget-submission/README.md"
echo "5. Submit the pull request"
echo
echo "Expected approval time: 1-7 days"
echo "Once approved, users can install with: winget install PlebOne.nos"
echo

cd ../nos
