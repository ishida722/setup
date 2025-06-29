#!/bin/bash

# Claude Code セットアップスクリプト (Ubuntu用)
set -e

# スクリプトディレクトリの取得
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_DIR="$SCRIPT_DIR/scripts"

# 共通ユーティリティを読み込み
source "$SCRIPTS_DIR/utils.sh"

# 各インストールスクリプトを読み込み
source "$SCRIPTS_DIR/install-nodejs.sh"
source "$SCRIPTS_DIR/install-claude-code.sh"
source "$SCRIPTS_DIR/install-neovim.sh"
source "$SCRIPTS_DIR/install-fish.sh"
source "$SCRIPTS_DIR/clone-configs.sh"

# メイン実行
main() {
    log "Claude Code セットアップを開始..."
    clone_configs
    install_nodejs
    install_claude_code
    install_neovim
    install_fish
    
    echo ""
    success "セットアップ完了！"
    echo "使用方法: claude-code --help"
    echo "API キー設定: export ANTHROPIC_API_KEY='your-key'"
    echo "Neovim: nvim"
    echo "注意: デフォルトシェルの変更は再ログイン後に有効になります"
}

main "$@"
