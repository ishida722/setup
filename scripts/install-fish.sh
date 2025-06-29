#!/bin/bash

# Fishシェルインストールスクリプト

# 共通ユーティリティを読み込み
source "$(dirname "$0")/utils.sh"

install_fish() {
    if command -v fish &> /dev/null; then
        log "Fish は既にインストール済み: $(fish --version)"
    else
        log "Fishをインストール中..."
        sudo apt-get update
        sudo apt-get install -y fish
        success "Fishをインストールしました"
    fi
    
    # デフォルトシェルをfishに変更
    FISH_PATH=$(which fish)
    if [ "$SHELL" != "$FISH_PATH" ]; then
        log "デフォルトシェルをfishに変更中..."
        chsh -s "$FISH_PATH"
        success "デフォルトシェルをfishに変更しました（再ログイン後に有効）"
    else
        log "デフォルトシェルは既にfishです"
    fi
}

# スクリプトが直接実行された場合
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_fish
fi