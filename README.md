# Ubuntu開発環境セットアップスクリプト

新しいUbuntuマシンにClaude Codeと開発ツールを一括でセットアップするスクリプト集です。

## 🚀 クイックスタート

```bash
curl -fsSL https://raw.githubusercontent.com/ishida722/setup/main/setup.sh | bash
```

## 📦 インストールされるツール

- **Node.js LTS** - NodeSourceリポジトリからの最新LTS版
- **Claude Code** - AnthropicのAI開発ツール（npm経由）
- **Neovim** - GitHub ReleasesからのNeovim最新版
- **Fish Shell** - 高機能なシェル（デフォルトシェルに設定）
- **Yazi** - 高機能なファイルマネージャー
- **Lazygit** - Git用のTUIツール
- **設定ファイル** - NeovimとFishの個人設定を外部リポジトリからクローン

## 🛠️ スクリプト構成

メインの`setup.sh`からモジュール化されたスクリプトを呼び出します：

```
scripts/
├── utils.sh                 # 共通ユーティリティ（ログ出力、色設定）
├── install-nodejs.sh        # Node.js LTS インストール
├── install-claude-code.sh   # Claude Code インストール
├── install-neovim.sh        # Neovim インストール
├── install-fish.sh          # Fish Shell インストール・設定
└── clone-configs.sh         # 設定ファイルのクローン
```

## ⚙️ システム要件

- **OS**: Ubuntu 20.04 LTS以降
- **権限**: sudo権限が必要
- **ネットワーク**: インターネット接続が必要

## 🔧 セットアップ後の設定

### Claude Code API キーの設定

```bash
# Anthropic API キーの設定
export ANTHROPIC_API_KEY='your-api-key-here'

# 永続化（.bashrcや.profileに追記）
echo 'export ANTHROPIC_API_KEY="your-api-key-here"' >> ~/.bashrc
source ~/.bashrc
```

### Fish Shell の有効化

```bash
# 再ログインするか、新しいターミナルを開いてFishを起動
fish
```

## 📁 外部設定リポジトリ

以下の設定ファイルが自動でクローンされます：

- **Neovim設定**: [ishida722/nvim](https://github.com/ishida722/nvim) → `~/.config/nvim/`
- **Fish設定**: [ishida722/fish](https://github.com/ishida722/fish) → `~/.config/fish/`

