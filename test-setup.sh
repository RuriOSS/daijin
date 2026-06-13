#!/bin/bash
# Test script to create mock containers for testing the TUI

set -e

CONTAINERS_DIR="${PREFIX:-/data/data/com.termux/files/usr}/var/daijin/containers"
HOME_DIR="${HOME:-/data/data/com.termux/files/home}"

echo "Creating test containers..."

# Create containers directory
mkdir -p "$CONTAINERS_DIR"

# Create test container directories
mkdir -p "$HOME_DIR/alpine-edge-test"
mkdir -p "$HOME_DIR/ubuntu-22.04-test"
mkdir -p "$HOME_DIR/debian-bookworm-test"

# Create mock config files
cat > "$CONTAINERS_DIR/alpine-edge-test.conf" <<EOF
backend="proot"
container_dir="$HOME_DIR/alpine-edge-test"
EOF

cat > "$CONTAINERS_DIR/ubuntu-22.04-test.conf" <<EOF
backend="ruri"
container_dir="$HOME_DIR/ubuntu-22.04-test"
extra_args="--verbose"
EOF

cat > "$CONTAINERS_DIR/debian-bookworm-test.conf" <<EOF
backend="proot"
container_dir="$HOME_DIR/debian-bookworm-test"
EOF

echo "Test containers created:"
ls -lh "$CONTAINERS_DIR"

echo ""
echo "Now you can run: ./daijin"
