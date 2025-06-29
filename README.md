# Ubuntu開発環境セットアップスクリプト

新しいUbuntuマシンにClaude Codeと開発ツールを一括でセットアップするスクリプト集です。

## 🚀 クイックスタート

```bash
# GitHubから直接実行
curl -fsSL https://raw.githubusercontent.com/ishida722/setup/main/setup.sh | bash

# ローカルで実行
git clone https://github.com/ishida722/setup.git
cd setup
chmod +x setup.sh
./setup.sh
```

## 📦 インストールされるツール

- **Node.js LTS** - NodeSourceリポジトリからの最新LTS版
- **Claude Code** - AnthropicのAI開発ツール（npm経由）
- **Neovim** - GitHub ReleasesからのNeovim最新版
- **Fish Shell** - 高機能なシェル（デフォルトシェルに設定）
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

### 個別実行

各スクリプトは個別でも実行可能です：

```bash
chmod +x scripts/*.sh

# Node.js のみインストール
./scripts/install-nodejs.sh

# Claude Code のみインストール
./scripts/install-claude-code.sh

# Neovim のみインストール
./scripts/install-neovim.sh

# Fish Shell のみセットアップ
./scripts/install-fish.sh

# 設定ファイルのみクローン
./scripts/clone-configs.sh
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

## 🔍 トラブルシューティング

### インストールが失敗する場合

```bash
# ログを確認
tail -f /var/log/setup.log

# パッケージリストを更新
sudo apt update

# 手動で再実行
./setup.sh
```

### Claude Code が動作しない場合

```bash
# インストール確認
claude-code --version

# npm グローバルパッケージ一覧
npm list -g --depth=0

# 再インストール
npm uninstall -g claude-code
npm install -g claude-code
```

## 🛡️ セキュリティ注意事項

- スクリプトはsudo権限でシステムファイルを変更します
- 実行前にスクリプト内容を確認することを推奨します
- 信頼できる環境でのみ実行してください

## 📝 ログ

インストール処理のログは以下に出力されます：
- 標準出力：カラー付きのリアルタイム進捗表示
- エラーログ：`/var/log/setup-error.log`（エラー発生時）

## 🤝 貢献

バグ報告や機能改善の提案は[Issues](https://github.com/ishida722/setup/issues)までお願いいたします。

## 📄 ライセンス

MIT License
