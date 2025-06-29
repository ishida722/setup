#!/bin/bash

# Node.jsインストールスクリプト

# 共通ユーティリティを読み込み
source "$(dirname "$0")/utils.sh"

install_nodejs() {
    if command -v node &> /dev/null; then
        log "Node.js は既にインストール済み: $(node --version)"
    else
        log "Node.jsをインストール中..."
        curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
        sudo apt-get install -y nodejs
        success "Node.jsをインストールしました: $(node --version)"
    fi
}

# スクリプトが直接実行された場合
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_nodejs
fi