#!/bin/bash

set -ex


TAG=""
if [ $# -eq 1 ]; then
    TAG="$1"
fi

# build dir
SCRIPT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
BUILD_DIR="$SCRIPT_DIR/build"
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# checksum file
CHECKSUM_FILE_NAME="SHA256SUMS.txt"
CHECKSUM_FILE_PATH="$BUILD_DIR/$CHECKSUM_FILE_NAME"
touch "$CHECKSUM_FILE_PATH"

PROJECT="BedrockWormhole"

PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    "windows/386"
    "windows/arm64"
    "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%/*}"
    ARCH="${PLATFORM#*/}"

    echo "Building for $OS/$ARCH..."

    TEMP_DIR=$(mktemp -d)

    if [ "$OS" = "windows" ]; then
        BINARY_NAME="$PROJECT.exe"
    else
        BINARY_NAME="$PROJECT"
    fi

    GOOS=$OS GOARCH=$ARCH go build -ldflags "-s -w" -o "$TEMP_DIR/$BINARY_NAME"

    [ -f README.md ] && cp README.md "$TEMP_DIR/"
    [ -f LICENSE.md ] && cp LICENSE.md "$TEMP_DIR/"

    COMPRESSED_NAME="${PROJECT}_${TAG:=dev}_${OS}_${ARCH}"
    if [ "$OS" = "windows" ]; then
        (cd "$TEMP_DIR" && zip -r "$BUILD_DIR/${COMPRESSED_NAME}.zip" .)
        (cd "$BUILD_DIR" && sha256sum "${COMPRESSED_NAME}.zip" >> "$CHECKSUM_FILE_NAME")
    else
        tar -czf "$BUILD_DIR/${COMPRESSED_NAME}.tgz" -C "$TEMP_DIR" .
        (cd "$BUILD_DIR" && sha256sum "${COMPRESSED_NAME}.tgz" >> "$CHECKSUM_FILE_NAME")
    fi

    rm -rf "$TEMP_DIR"
done

echo "Build complete! Artifacts are in the $BUILD_DIR directory"
echo "SHA256 checksums are available in $CHECKSUM_FILE"

# make a tag
if [ -n "$TAG" ]; then
    git tag -d "$TAG" || true
    git push --delete origin "$TAG" || true
    git tag "$TAG"
    git push origin "$TAG"
else
    echo "Warning: No tag was provided - git tag was not created or pushed"
fi
