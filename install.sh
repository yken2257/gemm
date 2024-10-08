#!/bin/bash

# If no version is specified as a command line argument, fetch the latest version.
if [ -z "$1" ]; then
  VERSION=$(curl -s https://api.github.com/repos/yken2257/gemm/releases/latest | grep -o '"tag_name": *"[^"]*"' | sed 's/"tag_name": *"//' | sed 's/"//')
  if [ -z "$VERSION" ]; then
    echo "Failed to fetch the latest version"
    exit 1
  fi
else
  VERSION=$1
fi

OS=$(uname -s)
ARCH=$(uname -m)
URL="https://github.com/yken2257/gemm/releases/download/${VERSION}/gemm_${OS}_${ARCH}.tar.gz"

echo "Start to install..."
echo "VERSION=$VERSION, OS=$OS, ARCH=$ARCH"
echo "URL=$URL"

TMP_DIR=$(mktemp -d)
curl -L $URL -o $TMP_DIR/gemm.tar.gz
tar -xzvf $TMP_DIR/gemm.tar.gz -C $TMP_DIR
sudo mv $TMP_DIR/gemm /usr/local/bin/gemm
sudo chmod +x /usr/local/bin/gemm

rm -rf $TMP_DIR

if [ -f "/usr/local/bin/gemm" ]; then
  echo "[SUCCESS] gemm $VERSION has been installed to /usr/local/bin"
else
  echo "[FAIL] gemm $VERSION is not installed."
fi