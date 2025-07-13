# Ubuntu開発環境セットアップスクリプト

新しいUbuntuマシンにClaude Codeと開発ツールを一括でセットアップするスクリプト集です。

## 🚀 クイックスタート

### Bashスクリプト版

```bash
curl -fsSL https://raw.githubusercontent.com/ishida722/setup/main/setup.sh | bash
```

### Ansible Playbook版

```bash
# Ansibleのインストール
sudo apt update
sudo apt install -y ansible

# Playbookの直接実行
curl -fsSL https://raw.githubusercontent.com/ishida722/setup/main/playbook.yml | ansible-playbook /dev/stdin
```

## 📦 インストールされるツール

- **Node.js LTS** - NodeSourceリポジトリからの最新LTS版
- **Claude Code** - AnthropicのAI開発ツール（npm経由）
- **Neovim** - GitHub ReleasesからのNeovim最新版
- **Fish Shell** - 高機能なシェル（デフォルトシェルに設定）
- **Yazi** - 高機能なファイルマネージャー
- **Lazygit** - Git用のTUIツール
- **設定ファイル** - NeovimとFishの個人設定を外部リポジトリからクローン

## 🛠️ ファイル構成

```
.
├── setup.sh          # Bashスクリプト版（従来の方法）
├── playbook.yml       # Ansible Playbook版（推奨）
└── README.md          # このファイル
```

## 🔄 Ansible vs Bashスクリプト

| 特徴 | Ansible Playbook | Bashスクリプト |
|------|------------------|----------------|
| **冪等性** | ✅ 自動保証 | ⚠️ 手動実装 |
| **エラーハンドリング** | ✅ 高機能 | ⚠️ 基本的 |
| **再実行安全性** | ✅ 完全対応 | ⚠️ 条件分岐で対応 |
| **保守性** | ✅ 宣言的 | ⚠️ 手続き型 |
| **学習コスト** | ⚠️ 中程度 | ✅ 低い |
| **依存関係** | ⚠️ Ansible必要 | ✅ bash標準 |

**推奨**: 開発環境では冪等性と安全性の高いAnsible Playbook版を使用

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

## 🎯 Ansible Playbookの特徴

### 冪等性
- 何度実行しても同じ結果
- 既にインストール済みのツールはスキップ
- 設定ファイルが存在する場合は上書きしない

### タスク実行例
```bash
# 詳細出力で実行
ansible-playbook playbook.yml -v

# 特定のタスクのみ実行（例：Yaziのインストールをスキップ）
ansible-playbook playbook.yml --skip-tags yazi

# ドライラン（実際の変更なし）
ansible-playbook playbook.yml --check
```

## 📁 外部設定リポジトリ

以下の設定ファイルが自動でクローンされます：

- **Neovim設定**: [ishida722/nvim](https://github.com/ishida722/nvim) → `~/.config/nvim/`
- **Fish設定**: [ishida722/fish](https://github.com/ishida722/fish) → `~/.config/fish/`
- **Krapp設定**: [ishida722/krapp-config](https://github.com/ishida722/krapp-config) → `~/.config/krapp/`

## 🐛 トラブルシューティング

### Ansible Playbook版

```bash
# Python3とpipが必要
sudo apt install -y python3-pip

# Ansibleのアップデート
pip3 install --upgrade ansible

# 権限エラーの場合
sudo ansible-playbook playbook.yml
```

