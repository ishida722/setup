#!/bin/bash

# Neovimインストールスクリプト

# 共通ユーティリティを読み込み
source "$(dirname "$0")/utils.sh"

install_neovim() {
    if command -v nvim &> /dev/null; then
        log "Neovim は既にインストール済み: $(nvim --version | head -1)"
    else
        log "Neovimをインストール中..."
        curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz
        sudo rm -rf /opt/nvim
        sudo tar -C /opt -xzf nvim-linux-x86_64.tar.gz
        export PATH="$PATH:/opt/nvim-linux-x86_64/bin"
        success "neovimをインストールしました"
        
        # 一時ファイルを削除
        rm -f nvim-linux-x86_64.tar.gz
    fi
}

# スクリプトが直接実行された場合
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_neovim
fi