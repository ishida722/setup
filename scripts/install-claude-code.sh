#!/bin/bash

# Claude Codeインストールスクリプト

# 共通ユーティリティを読み込み
source "$(dirname "$0")/utils.sh"

install_claude_code() {
    if command -v claude-code &> /dev/null; then
        log "Claude Code は既にインストール済み"
    else
        log "Claude Codeをインストール中..."
        npm install -g @anthropic-ai/claude-code
        success "Claude Codeをインストールしました"
    fi
}

# スクリプトが直接実行された場合
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    install_claude_code
fi