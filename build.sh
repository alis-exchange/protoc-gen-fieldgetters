#!/bin/bash

set -e  # Exit immediately on error
VERSION=$1
OUTPUT_DIR="dist/$VERSION"

# Create the output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# List of platforms and architectures to build for
platforms=("linux" "darwin" "windows")
architectures=("amd64" "386" "arm64")

# Build the plugin for each OS/arch combination
for OS in "${platforms[@]}"; do
  for ARCH in "${architectures[@]}"; do
    # Skip darwin/386 because it is not supported
    if [ "$OS" = "darwin" ] && [ "$ARCH" = "386" ]; then
      echo "Skipping unsupported GOOS/GOARCH pair darwin/386"
      continue
    fi

    EXT=""
    if [ "$OS" = "windows" ]; then
      EXT=".exe"
    fi

    # Set environment variables to target the OS and architecture
    GOOS=$OS GOARCH=$ARCH go build -ldflags "-X main.version=$VERSION" -o "./$OUTPUT_DIR/protoc-gen-go-fieldgetters-$OS-$ARCH$EXT" ./cmd/protoc-gen-go-fieldgetters

    echo "Built protoc-gen-go-fieldgetters@$VERSION for $OS/$ARCH"
  done
done

