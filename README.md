<div align="center">
  <img src="nos.jpeg" alt="nos logo" width="200" height="200">
  
  # nos - Beautiful Nostr CLI Client
  
  [![GitHub Release](https://img.shields.io/github/v/release/PlebOne/nos?style=flat-square)](https://github.com/PlebOne/nos/releases)
  [![Go Report Card](https://goreportcard.com/badge/github.com/PlebOne/nos?style=flat-square)](https://goreportcard.com/report/github.com/PlebOne/nos)
  [![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
  [![winget](https://img.shields.io/badge/winget-PlebOne.nos-blue?style=flat-square)](https://github.com/microsoft/winget-pkgs)
  
  A beautiful command-line client for posting to Nostr, built with Go and Charm.sh.
  
  🚀 **Now available on Microsoft winget!**
</div>

## ✨ Features

- 🔐 Secure key storage using system keyring
- 🎨 Beautiful terminal UI with Charm.sh
- 📡 Posts to multiple relays for reliability
- 🚀 Simple one-command posting
- ⚙️ Customizable relay management
- 🔄 Account switching with complete reset
- 📱 Interactive menu-driven interface
- ✅ Post verification across relays

## 🚀 Installation

### Linux (Debian/Ubuntu)

Download and install the .deb package:

```bash
# Download the latest .deb package from releases
wget https://github.com/PlebOne/nos/releases/latest/download/nos_1.1.7_amd64.deb

# Install the package
sudo dpkg -i nos_1.1.7_amd64.deb

# If you encounter dependency issues, run:
sudo apt-get install -f
```

**Alternative architectures:**
- ARM64: `nos_1.1.7_arm64.deb`
- ARMv7: `nos_1.1.7_armhf.deb`

### Linux (Snap Store)

Install from the Snap Store:

```bash
sudo snap install nos
```

The snap package provides automatic updates and works across all major Linux distributions.

### Windows (Recommended)

Install via Microsoft winget (Windows Package Manager):

```cmd
winget install PlebOne.nos
```

**Note**: If you encounter a security warning during winget installation, this is because the executable is not yet code-signed. You can:

1. **Alternative method**: Download manually from [releases](https://github.com/PlebOne/nos/releases) and run:
   ```powershell
   # Download the .exe file, then:
   Unblock-File -Path "path\to\nos-VERSION-windows-amd64.exe"
   ```

2. **One-liner install** (PowerShell as Administrator):
   ```powershell
   Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/PlebOne/nos/main/install.ps1'))
   ```

*Code signing is being implemented to resolve this issue permanently.*

### Manual Download

Download the latest release for your platform from the [releases page](https://github.com/PlebOne/nos/releases):

- **Windows**: Download `nos-VERSION-windows-amd64.exe` or `nos-VERSION-windows-arm64.exe`
- **Linux**: Download `nos_VERSION_linux_amd64.tar.gz` or `nos_VERSION_linux_arm64.tar.gz`
- **macOS**: Download `nos_VERSION_darwin_amd64.tar.gz` or `nos_VERSION_darwin_arm64.tar.gz`

**Windows users**: After downloading, right-click the .exe file → Properties → Check "Unblock" if present.

### Using Go

```bash
go install github.com/PlebOne/nos@latest
```

### Build from Source

```bash
git clone https://github.com/PlebOne/nos.git
cd nos
go build -o nos
# Linux/macOS: sudo mv nos /usr/local/bin/
# Windows: move nos.exe to a directory in your PATH
```



## Usage

### Interactive Menu Mode

Simply run `nos` without any arguments to enter the interactive menu:

```bash
nos
```

This opens a beautiful menu-driven interface where you can:
- Setup your account (add nsec)
- Post messages
- Verify your posts on relays
- Manage relays
- Reset your account

### Quick Command Mode

For faster posting, you can use command-line arguments or stdin:

```bash
# Post a message directly
nos "Hello Nostr world!"

# Post from stdin (best for hashtags and URLs)
echo "Check out #bitcoin at https://bitcoin.org" | nos

# Multi-line posts with stdin
echo "Line 1
Line 2
#hashtag https://example.com" | nos

# Using PowerShell (Windows)
"My #nostr post with https://example.com" | nos
```

**💡 Tip**: Use stdin (`echo "message" | nos`) when your message contains special characters like hashtags (`#`) or URLs, as it avoids shell escaping issues.

### First Time Setup

The first time you post, nos will prompt for your nsec (private key):

```bash
nos "My first post"
# Will prompt: "Enter your nsec key (starts with nsec1):"
```

### Verifying Your Posts

Check if your posts are visible on relays:

```bash
nos verify
```

This will:
- Show your npub (public key)
- Connect to each relay and check for your recent posts
- Display a summary of posts found on each relay

### Managing Relays

Customize which relays to use through the interactive menu or commands:

```bash
# Interactive relay management
nos relay

# Or use direct commands:
nos relay list                      # List current relays
nos relay add wss://relay.example.com    # Add a new relay
nos relay remove wss://relay.example.com # Remove a relay
nos relay reset                     # Reset to default relays
```

### Changing Accounts / Reset

To completely reset nos and change to a different Nostr account:

```bash
nos reset
```

This will:
- Delete your stored nsec key
- Remove any custom relay configuration
- Allow you to set up nos with a different account

## Security

Your nsec key is stored securely using your system's native keyring:
- **Linux**: Secret Service API (GNOME Keyring, KWallet)
- **macOS**: Keychain
- **Windows**: Windows Credential Manager

## Default Relays

The client posts to these relays by default:
- wss://relay.damus.io
- wss://nos.lol
- wss://relay.nostr.band
- wss://relay.current.fyi
- wss://relay.snort.social
- wss://relay.primal.net

## 🔧 Development & Release Automation

This project features a fully automated release pipeline that:

### Automated Releases
- **Multi-platform builds**: Automatically builds for Windows (x64/arm64), Linux (x64/arm64), and macOS (x64/arm64)
- **GitHub Releases**: Creates GitHub releases with pre-built binaries and checksums
- **Windows winget**: Automatically generates and commits winget manifests for Microsoft Store distribution

### Release Process
1. Create and push a version tag: `git tag v1.2.3 && git push --tags`
2. GitHub Actions automatically:
   - Builds binaries for all platforms using goreleaser
   - Extracts Windows .exe files for winget compliance
   - Uploads standalone .exe files alongside tar.gz archives
   - Generates winget manifests with correct SHA256 checksums
   - Commits manifests to the main branch for tracking

### winget Integration
- Windows users can install via: `winget install PlebOne.nos`
- Automatic manifest generation ensures Microsoft compliance
- Supports both x64 and arm64 Windows architectures
- Silent installation switches for enterprise deployment

### Contributing
When contributing new features or bug fixes:
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request
5. Once merged, maintainers can create a new release tag to trigger automation

## 📋 Requirements

- **Go 1.21+** for building from source
- **Windows Terminal** or ANSI-capable terminal for best experience on Windows
- **System keyring** for secure key storage (automatically available on all supported platforms)

## 📦 Dependencies

- [go-nostr](https://github.com/nbd-wtf/go-nostr) - Nostr protocol implementation
- [Charm.sh Huh](https://github.com/charmbracelet/huh) - Beautiful form inputs
- [Charm.sh Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [go-keyring](https://github.com/zalando/go-keyring) - Secure key storage

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/PlebOne/nos/issues)
- **Discussions**: [GitHub Discussions](https://github.com/PlebOne/nos/discussions)
- **Nostr**: Find us on Nostr using the relays listed above!
