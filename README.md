<div align="center">
  <img src="nos.jpeg" alt="nos logo" width="200" height="200">
  
  # nos - Simple Nostr CLI Client
  
  [![GitHub Release](https://img.shields.io/badge/release-v0.9-blue.svg?style=flat-square)](https://github.com/PlebOne/nos/releases)
  [![Go Report Card](https://goreportcard.com/badge/github.com/PlebOne/nos?style=flat-square)](https://goreportcard.com/report/github.com/PlebOne/nos)
  [![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
  
  A beautiful command-line client for posting to Nostr, built with Go and Charm.sh.
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

## Installation

### Debian/Ubuntu (Recommended)

Download the latest .deb package from the [releases page](https://github.com/PlebOne/nos/releases):

```bash
wget https://github.com/PlebOne/nos/releases/download/v0.9/nos_0.9_amd64.deb
sudo dpkg -i nos_0.9_amd64.deb
```

### Using Go

```bash
go install github.com/PlebOne/nos@latest
```

### Build from source

```bash
git clone https://github.com/PlebOne/nos.git
cd nos
go build -o nos
sudo mv nos /usr/local/bin/
```

### Other platforms

Download the appropriate binary for your platform from the [releases page](https://github.com/PlebOne/nos/releases).

### Windows

Build and run on Windows using PowerShell and Go. Recommended terminal: Windows Terminal or another ANSI-capable terminal.

```powershell
# From the repo root
go build -o nos.exe .
.\nos.exe
```

Alternatively, use the included helper script to build and package a zip:

```powershell
.\build-windows.ps1
```

### Distribute via winget

You can distribute `nos` through the Windows Package Manager (winget) by submitting a manifest to the community `winget-pkgs` repository. This repo includes a small helper to generate a manifest folder for a release.

1. Build and upload your release artifacts (zip for x64/arm64). Note the public URLs and SHA256 checksums for the uploaded zips.
2. Run the helper to generate the manifest files (example):

```powershell
Set-Location .\packaging\winget
.\generate-winget-manifest.ps1 -Version "0.9.1" -InstallerUrlX64 "https://example.com/nos_0.9.1_windows_amd64.zip" -Sha256X64 "<sha256>" -InstallerUrlArm64 "https://example.com/nos_0.9.1_windows_arm64.zip" -Sha256Arm64 "<sha256>"
```

3. Commit the generated manifest folder from `packaging/winget/PlebOne/nos/<version>` and open a PR against https://github.com/microsoft/winget-pkgs following their submission guidelines.

Notes:
- winget requires the installer URLs to be accessible publicly (GitHub releases or another static host are common).
- The helper produces a `manifest.yaml` matching winget's expected YAML structure; review it before submitting.

### Signing releases (GPG)

The release workflow can sign artifacts with GPG so signatures are published alongside the release. To enable GPG signing in CI:

1. Export your GPG private key (ASCII-armored):

```bash
gpg --export-secret-keys --armor <KEY_ID> > nos-gpg-private.asc
```

2. Add the file contents as a repository secret named `GPG_PRIVATE_KEY` (copy/paste the file contents).
3. If your key has a passphrase, add it as repository secret `GPG_PASSPHRASE`. If not, set `GPG_PASSPHRASE` to an empty string.

If you prefer sigstore/cosign keyless signing, let me know and I can wire that into the workflow instead.

### Signing releases (sigstore / cosign keyless) - recommended

This repo supports keyless signing with `cosign` via GitHub Actions OIDC. Keyless signing avoids storing private keys in repository secrets and is generally a better security posture.

To use keyless signing:

1. Ensure the `release` workflow has `id-token: write` permissions (already configured).
2. The workflow installs `cosign` and uses the OIDC flow to request signing identity at runtime.

When you push a tag (e.g., `v0.9.1`), the release workflow will:
- Run goreleaser to build artifacts.
- Use `cosign` keyless mode to sign archives and checksums and publish signatures alongside the release.

If you prefer to use a specific cosign key (not keyless), tell me and I can add that flow which will require storing a private key in GitHub Secrets.

If you prefer sigstore/cosign keyless signing, let me know and I can wire that into the workflow instead.



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

For faster posting, you can also use command-line arguments:

```bash
# Post a message directly
nos "Hello Nostr world!"

# First time setup will prompt for your nsec key
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

## Dependencies

- [go-nostr](https://github.com/nbd-wtf/go-nostr) - Nostr protocol implementation
- [Charm.sh Huh](https://github.com/charmbracelet/huh) - Beautiful form inputs
- [Charm.sh Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [go-keyring](https://github.com/zalando/go-keyring) - Secure key storage
