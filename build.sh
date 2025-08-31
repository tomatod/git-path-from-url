#!/bin/bash
set -e

DIST_DIR=dist

OS_LIST=("linux" "darwin" "windows")
ARCH_LIST=("amd64" "arm64")

rm -rf $DIST_DIR
mkdir -p $DIST_DIR

for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    dir="$DIST_DIR/$os-$arch"
    mkdir -p "$dir"
    
    ext=""
    if [ "$os" = "windows" ]; then
      ext=".exe"
    fi
    
    output="$dir/git-path-from-url$ext"
    
    echo "Building $output..."
    GOOS=$os GOARCH=$arch go build -o "$output"

    zip_file="$dir/git-path-from-url.zip"
    echo "Zipping $zip_file..."
    zip -j "$zip_file" "$output"
  done
done

echo "All builds and zip packages completed."
