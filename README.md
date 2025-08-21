<div align="center">
  <img src="nos.jpeg" alt="nos logo" width="200" height="200">
  
  # nos - Simple Nostr CLI Client
  
  [![GitHub Release](https://img.shields.io/github/release/timdev/nos.svg?style=flat-square)](https://github.com/timdev/nos/releases/latest)
  [![Go Report Card](https://goreportcard.com/badge/github.com/timdev/nos?style=flat-square)](https://goreportcard.com/report/github.com/timdev/nos)
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

Download the latest .deb package from the [releases page](https://github.com/timdev/nos/releases/latest):

```bash
wget https://github.com/timdev/nos/releases/latest/download/nos_0.9_amd64.deb
sudo dpkg -i nos_0.9_amd64.deb
```

### Using Go

```bash
go install github.com/timdev/nos@latest
```

### Build from source

```bash
git clone https://github.com/timdev/nos.git
cd nos
go build -o nos
sudo mv nos /usr/local/bin/
```

### Other platforms

Download the appropriate binary for your platform from the [releases page](https://github.com/timdev/nos/releases/latest).

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
