#!/usr/bin/env bash

set -e

detect_os() {
    unameOut="$(uname -s)"
    case "${unameOut}" in
        Linux*)     os="linux";;
        Darwin*)    os="darwin";;
        CYGWIN*)    os="windows";;
        MINGW*)     os="windows";;
        MSYS*)      os="windows";;
        *)          os="unknown"
    esac

    echo "Detected OS: $os"

    if [ "$os" = "unknown" ]; then
        echo "Unsupported OS."
        exit 1
    fi
}

install_buf() {
    echo "Installing Buf..."

    if [ "$os" = "darwin" ]; then
        echo "Using Homebrew for Buf installation..."
        brew install bufbuild/buf/buf
        echo "Buf installed: $(buf --version)"
        return
    fi

    if [ "$os" = "linux" ]; then
        echo "Installing Buf for Linux (binary download)..."
        curl -sSL \
            "https://github.com/bufbuild/buf/releases/latest/download/buf-linux-x86_64" \
            -o /usr/local/bin/buf
        chmod +x /usr/local/bin/buf
        echo "Buf installed: $(buf --version)"
        return
    fi

    if [ "$os" = "windows" ]; then
        echo "Please install Buf manually for Windows:"
        echo "https://github.com/bufbuild/buf/releases"
        return
    fi
}

install_go_dependencies() {
    echo "Installing go dependencies..."
    go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
    go get -tool go.uber.org/mock/mockgen@latest
    go get -tool github.com/securego/gosec/v2/cmd/gosec@latest
}

install_npm_dependencies() {
    npm install -g openapi-format@1.16.0
    npm install -g @redocly/cli@1.9.1
}

detect_os
install_buf
install_go_dependencies
install_npm_dependencies
echo "All dependencies installed successfully!"