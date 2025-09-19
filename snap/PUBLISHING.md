# Snap Store Publishing Workflow

This document describes the complete workflow for publishing the `nos` CLI application to the Snap Store.

## Prerequisites

1. **Install snapcraft**:
   ```bash
   sudo snap install snapcraft --classic
   ```

2. **Create a Snapcraft account** at [snapcraft.io](https://snapcraft.io/)

3. **Login to snapcraft**:
   ```bash
   snapcraft login
   ```

## One-time Setup

### 1. Register the snap name

```bash
snapcraft register nos
```

**Note**: Snap names are globally unique. If `nos` is already taken, you'll need to choose a different name and update the `name` field in `snap/snapcraft.yaml`.

### 2. Set up store permissions (if needed)

For snaps that need additional permissions, you may need to request store approval:

```bash
snapcraft list-revisions nos
snapcraft status nos
```

## Publishing Process

### Local Testing

Before publishing, always test locally:

```bash
# Build and test locally
./snap/build-snap.sh

# Install and test
sudo snap install --dangerous ./nos_*.snap
nos --help
nos

# Remove test installation
sudo snap remove nos
```

### Publishing to Snap Store

Use the automated script:

```bash
./snap/publish-snap.sh
```

Or manually:

```bash
# Build
snapcraft --destructive-mode

# Upload and release to edge channel
snapcraft upload nos_*.snap --release=edge

# Promote to other channels when ready
snapcraft promote nos --from-channel=edge --to-channel=beta
snapcraft promote nos --from-channel=beta --to-channel=candidate
snapcraft promote nos --from-channel=candidate --to-channel=stable
```

## Release Channels

The Snap Store has four release channels:

- **edge**: Latest development builds (automatic from CI)
- **beta**: Pre-release testing
- **candidate**: Release candidates
- **stable**: Production releases

### Typical workflow:
1. Commit code → auto-publish to **edge**
2. Test edge build → promote to **beta**
3. Test beta build → promote to **candidate**
4. Final testing → promote to **stable**

## Continuous Integration

### GitHub Actions Workflow

Create `.github/workflows/snap.yml`:

```yaml
name: Build and Publish Snap

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Build snap
      uses: snapcore/action-build@v1
      id: build
    
    - name: Upload snap artifact
      uses: actions/upload-artifact@v4
      with:
        name: snap
        path: ${{ steps.build.outputs.snap }}
    
    - name: Publish to edge
      if: github.ref == 'refs/heads/main'
      uses: snapcore/action-publish@v1
      env:
        SNAPCRAFT_STORE_CREDENTIALS: ${{ secrets.SNAPCRAFT_TOKEN }}
      with:
        snap: ${{ steps.build.outputs.snap }}
        release: edge
    
    - name: Publish to stable
      if: startsWith(github.ref, 'refs/tags/v')
      uses: snapcore/action-publish@v1
      env:
        SNAPCRAFT_STORE_CREDENTIALS: ${{ secrets.SNAPCRAFT_TOKEN }}
      with:
        snap: ${{ steps.build.outputs.snap }}
        release: stable
```

### Setup GitHub Secrets

1. Generate snapcraft token:
   ```bash
   snapcraft export-login --snaps=nos --channels=edge,beta,candidate,stable --acls package_upload,package_release -
   ```

2. Add the token as `SNAPCRAFT_TOKEN` secret in GitHub repository settings.

## Versioning

Update version in `snap/snapcraft.yaml` for each release:

```yaml
version: '1.1.4'  # Update this for new releases
```

Version format:
- Use semantic versioning (e.g., `1.2.3`)
- Prefix with `v` for git tags (e.g., `v1.2.3`)
- Snap versions don't need the `v` prefix

## Monitoring and Maintenance

### Check snap status
```bash
snapcraft status nos
snapcraft metrics nos
```

### View snap information
```bash
snap info nos
snapcraft list-revisions nos
```

### Snap Store Dashboard

Monitor your snap at: `https://snapcraft.io/nos`

## Troubleshooting

### Common Issues

1. **Build failures**: Check `snapcraft.yaml` syntax and Go version compatibility
2. **Permission issues**: Ensure correct plugs are defined in snapcraft.yaml
3. **Upload failures**: Verify snapcraft login and network connectivity

### Debug builds
```bash
# Build with debug output
snapcraft --debug

# Clean everything
snapcraft clean --destructive-mode

# Build specific part
snapcraft build nos
```

### Getting help
- [Snapcraft documentation](https://snapcraft.io/docs)
- [Snapcraft forum](https://forum.snapcraft.io/)
- [Snapcraft GitHub issues](https://github.com/snapcore/snapcraft/issues)

## Security Considerations

- The snap uses `strict` confinement for security
- Required plugs: `network`, `network-bind`, `home`, `password-manager-service`
- Keyring access is sandboxed to the snap's private area
- Network access is required for Nostr relay connections