# Snap Build Directory

This directory contains the configuration and scripts for building the `nos` snap package.

## Files:

- `../snapcraft.yaml` - Main snap configuration file
- `build-snap.sh` - Script to build the snap locally
- `publish-snap.sh` - Script to publish the snap to the Snap Store

## Building:

```bash
# Build the snap locally
./snap/local/build-snap.sh

# Test the snap locally
sudo snap install --dangerous ./nos_*.snap
```

## Publishing:

1. First time setup requires login to Snapcraft:
   ```bash
   snapcraft login
   ```

2. Register the snap name (one time only):
   ```bash
   snapcraft register nos
   ```

3. Build and publish:
   ```bash
   ./snap/local/publish-snap.sh
   ```