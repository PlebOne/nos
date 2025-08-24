# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

**nos** is a simple, beautiful command-line client for posting to Nostr (Notes and Other Stuff Transmitted by Relays). It's written in Go and features an interactive terminal UI built with Charm.sh components.

### Core Architecture

- **Single-file architecture**: The entire application is contained in `main.go` (~1000 lines)
- **Nostr protocol**: Uses the `go-nostr` library for protocol implementation
- **UI Components**: Built with Charm.sh ecosystem (`huh`, `lipgloss`) for forms and styling
- **Key Storage**: Secure key management via system keyring (`go-keyring`)
- **Multi-relay support**: Posts to multiple Nostr relays simultaneously

### Key Components

1. **Interactive Menu System**: Main menu-driven interface with options for setup, posting, verification, and relay management
2. **Key Management**: Secure nsec (private key) storage using OS keyring services
3. **Relay Management**: Configurable relay lists with fallback to defaults
4. **Post Verification**: Ability to check if posts appear on configured relays

## Common Development Tasks

### Building and Running
```bash
# Build the binary
go build -o nos

# Run directly from source
go run main.go

# Quick post message
./nos "Hello Nostr world!"

# Interactive menu mode
./nos
```

### Testing the Application
```bash
# Run with interactive menu to test all features
go run main.go

# Test direct posting (will prompt for setup if first time)
go run main.go "Test message"

# Test specific commands
go run main.go verify
go run main.go relay list
go run main.go reset
```

### Creating Debian Package
```bash
# Build debian package using the provided script
./build-deb.sh

# Package will be created in dist/ directory
```

### Module Management
```bash
# Update dependencies
go mod tidy

# Upgrade go-nostr (main dependency)
go get -u github.com/nbd-wtf/go-nostr

# Upgrade Charm.sh components
go get -u github.com/charmbracelet/huh
go get -u github.com/charmbracelet/lipgloss
```

## Code Architecture Details

### Main Application Flow

The application operates in two modes:
1. **Interactive Menu Mode** (default when no args): Shows a TUI menu for all operations
2. **Command Mode**: Direct CLI commands for specific actions

### Key Functions and Their Purposes

- `showMainMenu()`: Main interactive loop with account status display
- `postToNostr()`: Core posting logic with multi-relay publishing
- `handleVerify()`: Post verification across all configured relays
- `showRelayMenu()`: Interactive relay management
- `getActiveRelays()`: Returns current relay list (custom or defaults)

### State Management

Application state is managed through the OS keyring:
- `nsec` key: Private key storage
- `relay-list` key: JSON-encoded custom relay list

### Default Relay Configuration

The application includes sensible defaults:
- wss://relay.damus.io
- wss://nos.lol
- wss://relay.nostr.band
- wss://relay.current.fyi
- wss://relay.snort.social
- wss://relay.primal.net

### Error Handling Patterns

- Network timeouts: 10s for connections, 5s for publishing
- Graceful degradation: Continue if some relays fail
- User-friendly error messages with styled output
- Validation at form input level using huh validators

## Development Guidelines

### UI/UX Consistency
- Use predefined lipgloss styles: `titleStyle`, `successStyle`, `errorStyle`, `infoStyle`
- Maintain consistent messaging patterns ("✓ success", "→ processing", etc.)
- Always provide "Press Enter to continue..." for interactive flows

### Nostr Protocol Considerations
- Always validate nsec keys before storage
- Generate proper event IDs and signatures
- Use appropriate event kinds (KindTextNote for posts)
- Include proper timestamps using `nostr.Now()`

### Security Practices
- Private keys never appear in logs or output
- Use EchoMode(huh.EchoModePassword) for key input
- Validate key format before storage
- Secure deletion on reset operations

### Testing Considerations
- The app requires system keyring access for full testing
- Network connectivity needed for relay operations
- Consider using test relays for development
- Account setup is required for most functionality

## File Structure and Dependencies

### Core Dependencies
- `github.com/nbd-wtf/go-nostr`: Nostr protocol implementation and NIP-19 encoding
- `github.com/charmbracelet/huh`: Interactive forms and prompts
- `github.com/charmbracelet/lipgloss`: Terminal styling and colors
- `github.com/zalando/go-keyring`: Cross-platform secure key storage

### Build Artifacts
- `nos`: Compiled binary
- `dist/nos_*.deb`: Debian package (created by build-deb.sh)
- `build/`: Temporary build directory for packaging

### Configuration
- No config files - uses keyring for persistent state
- Relay configuration stored as JSON in keyring
- Self-contained with sensible defaults
