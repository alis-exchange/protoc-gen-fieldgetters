#!/bin/bash

# Variables
VERSION="v0.0.4"
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Define architecture mappings if necessary
if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
fi

# Default installation directory
INSTALL_DIR="$GOPATH/bin"

# Download URL
BINARY_URL="https://github.com/your-username/protoc-gen-fieldgetters/releases/download/$VERSION/protoc-gen-go-fieldgetters-$PLATFORM-$ARCH"

# Download the binary
echo "Downloading protoc-gen-go-fieldgetters $VERSION for $PLATFORM/$ARCH..."
curl -L -o protoc-gen-go-fieldgetters "$BINARY_URL"

# Make the binary executable
chmod +x protoc-gen-go-fieldgetters

# Install the binary to the specified directory or default
echo "Installing protoc-gen-go-fieldgetters to $INSTALL_DIR..."
sudo mv protoc-gen-go-fieldgetters "$INSTALL_DIR/"

# Verify installation
echo "Installation complete. Verifying..."
protoc-gen-go-fieldgetters --version
