# トラブルシューティングガイド

このドキュメントでは、Ubuntu開発環境セットアップ時によく発生する問題と解決方法について説明します。

## Ansibleプレイブック実行時のエラー

### aptキャッシュ更新エラー

#### 症状
Ansibleプレイブック実行時に以下のエラーが発生する：

```
TASK [Install basic dependencies and Fish shell] ***********************************************************************
fatal: [localhost]: FAILED! => {"changed": false, "msg": "Failed to update apt cache: unknown reason"}
```

#### よくある原因

**1. サードパーティリポジトリのGPGキーエラー（最も一般的）**
- Steamリポジトリのキーが見つからない: `NO_PUBKEY F24AEA9FB05498B7`
- その他のソフトウェアリポジトリの署名が無効

**2. ネットワーク接続の問題**
- リポジトリサーバーにアクセスできない
- プロキシ設定の問題

**3. ディスク容量不足**
- `/var/cache/apt/` の容量不足
- ルートファイルシステムの容量不足

#### 診断手順

**1. apt updateを直接実行してエラーを確認**
```bash
sudo apt update
```

**2. GPGキーエラーを確認**
```bash
sudo apt update 2>&1 | grep -i "NO_PUBKEY\|GPG"
```

**3. ネットワーク接続を確認**
```bash
ping -c 3 archive.ubuntu.com
```

**4. ディスク容量を確認**
```bash
df -h /var/cache/apt/
```

#### 解決方法

**SteamリポジトリのGPGキーエラーの場合：**

```bash
# 正しいSteam GPGキーを追加
curl -fsSL https://repo.steampowered.com/steam/archive/stable/steam.gpg | sudo gpg --dearmor -o /usr/share/keyrings/steam.gpg

# Steamリポジトリ設定を修正
sudo sh -c 'echo "deb [arch=amd64,i386 signed-by=/usr/share/keyrings/steam.gpg] https://repo.steampowered.com/steam stable steam" > /etc/apt/sources.list.d/steam-stable.list'

# aptキャッシュを更新
sudo apt update
```

**一般的なリポジトリ問題の場合：**

```bash
# 問題のあるリポジトリを一時的に無効化
sudo mv /etc/apt/sources.list.d/問題のあるリポジトリ.list /etc/apt/sources.list.d/問題のあるリポジトリ.list.bak

# aptキャッシュを更新
sudo apt update

# 正しい設定でリポジトリを再追加
```

**ネットワーク問題の場合：**

```bash
# DNSを確認
nslookup archive.ubuntu.com

# プロキシ設定を確認（必要に応じて）
echo $http_proxy
echo $https_proxy
```

**ディスク容量不足の場合：**

```bash
# aptキャッシュをクリア
sudo apt clean

# 不要なパッケージを削除
sudo apt autoremove

# ディスク容量を再確認
df -h
```

#### 予防策

- **手動ソフトウェアインストール時は必ずGPGキーを正しく設定する**
- **可能な限り公式パッケージリポジトリを使用する**
- **定期的なシステムメンテナンスでディスク容量不足を防ぐ**
- **サードパーティリポジトリ追加時は信頼できるソースのみを使用する**

## その他のよくある問題

### Node.js インストールエラー

**症状：** NodeSourceリポジトリの追加に失敗

**解決方法：**
```bash
# NodeSourceの公式スクリプトを手動実行
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### Neovim バイナリダウンロードエラー

**症状：** GitHub Releasesからのダウンロードに失敗

**解決方法：**
```bash
# 手動でNeovimをダウンロード・インストール
cd /tmp
wget https://github.com/neovim/neovim/releases/latest/download/nvim-linux64.tar.gz
sudo tar -C /opt -xzf nvim-linux64.tar.gz
sudo ln -sf /opt/nvim-linux64/bin/nvim /usr/local/bin/nvim
```

### 権限エラー

**症状：** `become: yes` 使用時にsudoパスワードを求められる

**解決方法：**
```bash
# パスワードを求められる場合は --ask-become-pass オプションを使用
ansible-playbook playbook.yml --ask-become-pass
```

## サポート

問題が解決しない場合：

1. **エラーメッセージ全体をコピー**してサポートに連絡
2. **実行環境の情報**を提供：
   - Ubuntu バージョン: `lsb_release -a`
   - Python バージョン: `python3 --version`
   - Ansible バージョン: `ansible --version`
3. **実行したコマンド**と**出力結果**を記録

---

最終更新: 2025年7月14日