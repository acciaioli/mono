#!/bin/bash

# exit when any command fails
set -e

declare -a GOOS=("linux" "darwin" "windows")
declare -a GOARCH=("amd64")

if [[ -z "${VERSION}" ]]; then
  echo "VERSION environment variable not exported"
  exit 1
fi

for goos in "${GOOS[@]}"; do
  for goarch in "${GOARCH[@]}"; do
    bin="bin/mono.$goos-$goarch"
    if [[ "$bin" == *"windows"* ]]; then
      bin="$bin.exe"
    fi
    cmd="GOOS=$goos GOARCH=$goarch go build -ldflags \"-X main.version=$VERSION\" -a -o $bin ./cmd/mono"
    echo "$cmd"
    eval "$cmd"
  done
done
