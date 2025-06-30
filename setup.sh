#!/bin/bash

# Claude Code セットアップスクリプト (Ubuntu用)
# 開発ルール: 可能な限りaptなどのパッケージマネージャーを使用する
# 例外: NeoVimはaptのバージョンが古いため、バイナリからインストール

# カラー出力
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# エラーハンドリング付きで関数を実行
run_with_error_handling() {
    local func_name="$1"
    local description="$2"
    
    log "$description を開始..."
    if $func_name; then
        success "$description が完了しました"
        return 0
    else
        error "$description でエラーが発生しました。処理を続行します..."
        return 1
    fi
}

# Node.jsのインストール
install_nodejs() {
    if command -v node &> /dev/null; then
        log "Node.js は既にインストール済み: $(node --version)"
        return 0
    else
        curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash - || return 1
        sudo apt-get install -y nodejs || return 1
        log "Node.jsをインストールしました: $(node --version)"
        return 0
    fi
}

# Claude Codeのインストール
install_claude_code() {
    if command -v claude-code &> /dev/null; then
        log "Claude Code は既にインストール済み"
        return 0
    else
        sudo npm install -g @anthropic-ai/claude-code || return 1
        log "Claude Codeをインストールしました"
        return 0
    fi
}

# Neovimのインストール
install_neovim() {
    if command -v nvim &> /dev/null; then
        log "Neovim は既にインストール済み: $(nvim --version | head -1)"
        return 0
    else
        curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz || return 1
        sudo rm -rf /opt/nvim || return 1
        sudo tar -C /opt -xzf nvim-linux-x86_64.tar.gz || return 1
        export PATH="$PATH:/opt/nvim-linux-x86_64/bin"
        log "neovimをインストールしました"
        return 0
    fi
}

# Yaziのインストール
install_yazi() {
    if command -v yazi &> /dev/null; then
        log "Yazi は既にインストール済み: $(yazi --version)"
        return 0
    else
        # 必須依存関係をインストール
        sudo apt-get update || return 1
        sudo apt-get install -y file || return 1
        
        # 推奨依存関係をインストール
        sudo apt-get install -y ffmpeg p7zip-full jq poppler-utils fd-find ripgrep fzf zoxide imagemagick xclip || return 1
        
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
        YAZI_URL=$(curl -s https://api.github.com/repos/sxyazi/yazi/releases/latest | grep "browser_download_url.*${YAZI_ARCH}.tar.gz" | cut -d '"' -f 4) || return 1
        
        if [ -z "$YAZI_URL" ]; then
            log "エラー: Yaziの最新リリースURLを取得できませんでした"
            return 1
        fi
        
        # ダウンロードとインストール
        cd /tmp || return 1
        curl -L "$YAZI_URL" -o yazi.tar.gz || return 1
        tar -xzf yazi.tar.gz || return 1
        cd yazi-* || return 1
        sudo cp yazi /usr/local/bin/ || return 1
        sudo cp ya /usr/local/bin/ || return 1
        sudo chmod +x /usr/local/bin/yazi /usr/local/bin/ya || return 1
        cd / || return 1
        rm -rf /tmp/yazi* || return 1
        
        log "Yaziとその依存関係をインストールしました: $(yazi --version)"
        return 0
    fi
}

# Lazygitのインストール
install_lazygit() {
    if command -v lazygit &> /dev/null; then
        log "Lazygit は既にインストール済み: $(lazygit --version)"
        return 0
    else
        sudo apt-get update || return 1
        sudo apt-get install -y lazygit || return 1
        log "Lazygitをインストールしました: $(lazygit --version)"
        return 0
    fi
}

# Fishのインストールとデフォルトシェル設定
install_fish() {
    if command -v fish &> /dev/null; then
        log "Fish は既にインストール済み: $(fish --version)"
    else
        sudo apt-get update || return 1
        sudo apt-get install -y fish || return 1
        log "Fishをインストールしました"
    fi
    
    # デフォルトシェルをfishに変更
    FISH_PATH=$(which fish)
    if [ "$SHELL" != "$FISH_PATH" ]; then
        chsh -s "$FISH_PATH" || {
            log "デフォルトシェルの変更に失敗しましたが、処理を続行します"
            return 0
        }
        log "デフォルトシェルをfishに変更しました（再ログイン後に有効）"
    else
        log "デフォルトシェルは既にfishです"
    fi
    return 0
}

# 設定ファイルのクローン
clone_configs() {
    mkdir -p ~/.config || return 1
    
    # Neovim設定
    if [ -d ~/.config/nvim ]; then
        log "Neovim設定は既に存在します"
    else
        git clone https://github.com/ishida722/nvim ~/.config/nvim || return 1
        log "Neovim設定をクローンしました"
    fi
    
    # Fish設定
    if [ -d ~/.config/fish ]; then
        log "Fish設定は既に存在します"
    else
        git clone https://github.com/ishida722/fish ~/.config/fish || return 1
        log "Fish設定をクローンしました"
    fi
    return 0
}

# メイン実行
main() {
    log "Claude Code セットアップを開始..."
    
    # エラーが発生してもセットアップを続行
    run_with_error_handling clone_configs "設定ファイルのクローン"
    run_with_error_handling install_nodejs "Node.jsのインストール"
    run_with_error_handling install_claude_code "Claude Codeのインストール"
    run_with_error_handling install_neovim "Neovimのインストール"
    run_with_error_handling install_yazi "Yaziのインストール"
    run_with_error_handling install_lazygit "Lazygitのインストール"
    run_with_error_handling install_fish "Fishのインストール"
    
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
