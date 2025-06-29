#!/bin/bash

# 設定ファイルクローンスクリプト

# 共通ユーティリティを読み込み
source "$(dirname "$0")/utils.sh"

clone_configs() {
    mkdir -p ~/.config
    
    # Neovim設定
    if [ -d ~/.config/nvim ]; then
        log "Neovim設定は既に存在します"
    else
        log "Neovim設定をクローン中..."
        git clone https://github.com/ishida722/nvim ~/.config/nvim
        success "Neovim設定をクローンしました"
    fi
    
    # Fish設定
    if [ -d ~/.config/fish ]; then
        log "Fish設定は既に存在します"
    else
        log "Fish設定をクローン中..."
        git clone https://github.com/ishida722/fish ~/.config/fish
        success "Fish設定をクローンしました"
    fi
}

# スクリプトが直接実行された場合
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    clone_configs
fi