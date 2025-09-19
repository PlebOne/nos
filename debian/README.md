# Debian Package (.deb) Setup

This document explains the Debian package setup for the `nos` CLI application.

## Files Created/Modified

### Build Scripts
- `build-deb.sh` - Main script to build .deb packages (updated to v1.1.4)
- `build-all-debs.sh` - Convenience script to build all architectures

### GitHub Actions
- `.github/workflows/release.yml` - Added debian job to automatically build and upload .deb packages on releases

### Documentation
- `README.md` - Updated with .deb installation instructions

## Architecture Support

The .deb packages are built for three architectures:
- **amd64** - Intel/AMD 64-bit (most common)
- **arm64** - ARM 64-bit (newer ARM processors, Apple M1, etc.)
- **armhf** - ARM 32-bit (Raspberry Pi, older ARM devices)

## Local Testing

### Build single architecture:
```bash
./build-deb.sh amd64
```

### Build all architectures:
```bash
./build-all-debs.sh
```

### Test installation:
```bash
sudo dpkg -i dist/nos_1.1.4_amd64.deb
nos --help
```

### Check package info:
```bash
dpkg -I dist/nos_1.1.4_amd64.deb
dpkg -c dist/nos_1.1.4_amd64.deb
```

### Remove package:
```bash
sudo dpkg -r nos
```

## Package Details

- **Package name**: `nos`
- **Version**: Automatically synced with Git tags
- **Maintainer**: PlebOne <contact@plebdev.org>
- **Dependencies**: `libc6`
- **Installation path**: `/usr/local/bin/nos`
- **Documentation**: `/usr/share/doc/nos/`
- **Logo**: `/usr/share/pixmaps/nos.png`

## Automated Release Process

When you create a new Git tag (e.g., `v1.1.5`), the GitHub Actions workflow will:

1. Build the Go binaries using goreleaser
2. Build .deb packages for all architectures
3. Upload .deb packages to the GitHub release
4. Generate SHA256 checksums for verification

## Installation for Users

### Download and install manually:
```bash
# Download the latest .deb package
wget https://github.com/PlebOne/nos/releases/latest/download/nos_1.1.4_amd64.deb

# Install
sudo dpkg -i nos_1.1.4_amd64.deb

# Fix dependencies if needed
sudo apt-get install -f
```

### One-liner installation:
```bash
# For amd64 systems
curl -sL https://github.com/PlebOne/nos/releases/latest/download/nos_1.1.4_amd64.deb -o nos.deb && sudo dpkg -i nos.deb && rm nos.deb
```

## Package Standards

The package follows Debian packaging standards:
- **FHS compliance**: Binary in `/usr/local/bin/`
- **Documentation**: In `/usr/share/doc/nos/`
- **Proper metadata**: control, copyright, changelog files
- **Dependencies**: Minimal (only libc6)
- **Architecture-specific**: Separate packages for different architectures

## Troubleshooting

### Build issues:
- Ensure Go 1.25+ is installed
- Check that `dpkg-dev` package is installed: `sudo apt-get install dpkg-dev`

### Installation issues:
- Missing dependencies: `sudo apt-get install -f`
- Permission issues: Ensure using `sudo` for installation
- Architecture mismatch: Download correct architecture package

### Testing in Docker:
```bash
# Test in clean Debian environment
docker run -it --rm -v $(pwd)/dist:/packages debian:bookworm bash
apt-get update && apt-get install -y /packages/nos_1.1.4_amd64.deb
nos --help
```