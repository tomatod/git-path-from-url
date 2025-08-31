#!/bin/bash
set -e

BIN_DIR=bin

OS_LIST=("linux" "darwin" "windows")
ARCH_LIST=("amd64" "arm64")

SRC="main.go"

rm -drf $BIN_DIR
mkdir -p $BIN_DIR

for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    dir="$BIN_DIR/$os-$arch"
    mkdir -p "$dir"
    
    ext=""
    if [ "$os" = "windows" ]; then
      ext=".exe"
    fi
    
    output="$dir/git-path-from-url$ext"
    
    echo "Building $output..."
    GOOS=$os GOARCH=$arch go build -o "$output" $SRC
  done
done
