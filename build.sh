#!/bin/bash

# Go版セットアップツールのビルドスクリプト

set -e

echo "Building Claude Setup (Go version)..."
go mod tidy
go build -o claude-setup main.go

echo "Build completed! Binary created: ./claude-setup"
echo "Usage: ./claude-setup"