#!/bin/bash

# Claude Code セットアップスクリプト (Ubuntu用)
# 開発ルール: 可能な限りaptなどのパッケージマネージャーを使用する
# 例外: NeoVimはaptのバージョンが古いため、バイナリからインストール
set -e

# カラー出力
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Node.jsのインストール
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

# Claude Codeのインストール
install_claude_code() {
    if command -v claude-code &> /dev/null; then
        log "Claude Code は既にインストール済み"
    else
        log "Claude Codeをインストール中..."
        sudo npm install -g @anthropic-ai/claude-code
        success "Claude Codeをインストールしました"
    fi
}

# Neovimのインストール
install_neovim() {
    if command -v nvim &> /dev/null; then
        log "Neovim は既にインストール済み: $(nvim --version | head -1)"
    else
        curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz
        sudo rm -rf /opt/nvim
        sudo tar -C /opt -xzf nvim-linux-x86_64.tar.gz
        export PATH="$PATH:/opt/nvim-linux-x86_64/bin"
        success "neovimをインストールしました"
    fi
}

# Yaziのインストール
install_yazi() {
    if command -v yazi &> /dev/null; then
        log "Yazi は既にインストール済み: $(yazi --version)"
    else
        log "Yaziをインストール中..."
        
        # 必須依存関係をインストール
        sudo apt-get update
        sudo apt-get install -y file
        
        # 推奨依存関係をインストール
        sudo apt-get install -y ffmpeg p7zip-full jq poppler-utils fd-find ripgrep fzf zoxide imagemagick xclip
        
        # Nerd Fontsがインストール済みか確認（推奨）
        if ! fc-list | grep -i "nerd" &> /dev/null; then
            log "推奨: Nerd Fontsをインストールしてください"
        fi
        
        # Yaziバイナリをダウンロード
        ARCH=$(uname -m)
        if [ "$ARCH" = "x86_64" ]; then
            YAZI_ARCH="x86_64-unknown-linux-gnu"
        elif [ "$ARCH" = "aarch64" ]; then
            YAZI_ARCH="aarch64-unknown-linux-gnu"
        else
            log "エラー: サポートされていないアーキテクチャ: $ARCH"
            return 1
        fi
        
        # 最新リリースのURLを取得
        YAZI_URL=$(curl -s https://api.github.com/repos/sxyazi/yazi/releases/latest | grep "browser_download_url.*${YAZI_ARCH}.tar.gz" | cut -d '"' -f 4)
        
        if [ -z "$YAZI_URL" ]; then
            log "エラー: Yaziの最新リリースURLを取得できませんでした"
            return 1
        fi
        
        # ダウンロードとインストール
        cd /tmp
        curl -L "$YAZI_URL" -o yazi.tar.gz
        tar -xzf yazi.tar.gz
        cd yazi-*
        sudo cp yazi /usr/local/bin/
        sudo cp ya /usr/local/bin/
        sudo chmod +x /usr/local/bin/yazi /usr/local/bin/ya
        cd /
        rm -rf /tmp/yazi*
        
        success "Yaziとその依存関係をインストールしました: $(yazi --version)"
    fi
}

# Lazygitのインストール
install_lazygit() {
    if command -v lazygit &> /dev/null; then
        log "Lazygit は既にインストール済み: $(lazygit --version)"
    else
        log "Lazygitをインストール中..."
        sudo apt-get update
        sudo apt-get install -y lazygit
        success "Lazygitをインストールしました: $(lazygit --version)"
    fi
}

# Fishのインストールとデフォルトシェル設定
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

# 設定ファイルのクローン
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

# メイン実行
main() {
    log "Claude Code セットアップを開始..."
    clone_configs
    install_nodejs
    install_claude_code
    install_neovim
    install_yazi
    install_lazygit
    install_fish
    
    echo ""
    success "セットアップ完了！"
    echo "使用方法: claude-code --help"
    echo "API キー設定: export ANTHROPIC_API_KEY='your-key'"
    echo "Neovim: nvim"
    echo "Yazi: yazi"
    echo "Lazygit: lazygit"
    echo "注意: デフォルトシェルの変更は再ログイン後に有効になります"
}

main "$@"
