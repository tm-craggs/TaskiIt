#!/usr/bin/env bash

set -e

VERSION="v1.0.1"
ARCHIVE="tidytask-linux-x86.tar.gz"
URL="https://github.com/tm-craggs/tidytask/releases/download/$VERSION/$ARCHIVE"

echo "Downloading TidyTask $VERSION..."
curl -fLO "$URL" || {
  echo "Failed to download $ARCHIVE. Check the URL or version."
  exit 1
}

echo "Extracting..."
[ -f tidytask ] && rm tidytask
tar -xzf "$ARCHIVE"

echo "Installing to /usr/local/bin..."
chmod +x tidytask
sudo mv tidytask /usr/local/bin

echo "Cleaning up..."
rm "$ARCHIVE"

echo "Installation successful"
