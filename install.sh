#!/usr/bin/env bash

set -e

VERSION="v1.0.1"
ARCHIVE="tidytask-linux-x86.tar.gz"
URL="https://github.com/tm-craggs/tidytask/releases/download/$VERSION/$ARCHIVE"

echo "Downloading TidyTask $VERSION..."
if ! curl -fLO "$URL"; then
  echo "❌ Failed to download $ARCHIVE from $URL"
  exit 1
fi

echo "Extracting..."
[ -f tidytask ] && rm tidytask
tar -xzf "$ARCHIVE"

echo "Installing to /usr/local/bin..."
chmod +x tidytask
sudo mv tidytask /usr/local/bin

echo "Cleaning up..."
rm "$ARCHIVE"

echo "Installation successful"
