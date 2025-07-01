# Go版setup.sh 共通インストール処理仕様

## 概要

Bash版setup.shのGo移植における共通インストール処理の設計仕様書

## 現在のBash版構造分析

### ファイル構成
- `setup.sh`: メインオーケストレータ
- `scripts/utils.sh`: 共通ユーティリティ関数
- `scripts/install-*.sh`: 各ソフトウェア個別インストールスクリプト

### 共通パターン
```bash
install_software() {
    if command -v software &> /dev/null; then
        log "software は既にインストール済み: $(software --version)"
    else
        log "softwareをインストール中..."
        # インストールコマンド実行
        success "softwareをインストールしました: $(software --version)"
    fi
}
```

## Go版設計仕様

### 1. データ構造

```go
// インストールコマンド定義
type InstallCommand struct {
    CheckCommands  []string // インストール済み確認コマンドリスト (例: ["node --version"])
    InstallCommands []string // インストール実行コマンドリスト
    InstallFunc    func() error // カスタムインストール処理関数（デフォルト: 何もしない関数）
    Name          string   // ソフトウェア名 (例: "Node.js")
    Description   string   // 説明
}

// インストール結果
type InstallResult struct {
    AlreadyInstalled bool    // 既にインストール済みか
    Success         bool     // インストール成功か
    Version         string   // インストール済みバージョン
    Err             error    // エラー情報
}
```

### 2. 共通インストール処理関数

```go
// メイン処理関数
func ExecuteInstall(cmd InstallCommand) InstallResult

// 内部ヘルパー関数
func checkInstalled(checkCmds []string) (bool, string, error)
func runInstallCommands(installCmds []string) error
func getVersion(checkCmds []string) (string, error)
func defaultInstallFunc() error { return nil } // デフォルトのInstallFunc
```

### 3. 主要機能仕様

#### インストール済み確認
- `CheckCommands`リストの各コマンドを順次実行してインストール状態を判定
- いずれかのコマンドが成功した場合: インストール済み、バージョン情報取得
- すべてのコマンドが失敗した場合: 未インストールと判定

#### インストール実行
- 未インストール時のみ以下を順次実行:
  1. `InstallCommands`リストの各コマンドを順次実行
  2. `InstallFunc`カスタム関数を実行（InstallCommandsの実行後）
- コマンド実行結果の成功/失敗を判定
- エラー時は適切なエラーメッセージを返却

#### ログ処理
- カラー付きログ出力
  - INFO: 青色 `[INFO]`
  - SUCCESS: 緑色 `[SUCCESS]`
  - ERROR: 赤色 `[ERROR]`
- 統一されたフォーマット

#### エラーハンドリング
- コマンド実行エラーの適切な捕捉
- エラー詳細情報の保持
- 処理継続/停止の制御可能

### 4. 使用例

```go
// Node.jsインストール例
nodeCmd := InstallCommand{
    CheckCommands:   []string{"node --version", "nodejs --version"},
    InstallCommands: []string{
        "curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -",
        "sudo apt-get install -y nodejs",
    },
    InstallFunc:    defaultInstallFunc, // デフォルト関数（何もしない）
    Name:          "Node.js",
    Description:   "JavaScript runtime",
}

// カスタム処理が必要な場合の例
complexCmd := InstallCommand{
    CheckCommands:   []string{"custom-tool --version"},
    InstallCommands: []string{"sudo apt-get update"},
    InstallFunc: func() error {
        // 複雑なカスタムインストール処理
        // 設定ファイル作成、権限設定など
        return nil
    },
    Name:          "Custom Tool",
    Description:   "複雑なインストール処理が必要なツール",
}

result := ExecuteInstall(nodeCmd)
if result.Err != nil {
    // エラー処理
}
```

### 5. 拡張性

- 新しいソフトウェアは`InstallCommand`構造体で簡単に追加可能
- 複雑なインストール処理も`InstallCommand`の拡張で対応
- ログレベルや出力形式のカスタマイズ可能

## 実装時の考慮事項

1. **セキュリティ**: sudoコマンドの実行制御
2. **パフォーマンス**: 並行インストール処理の検討
3. **エラー復旧**: 失敗時のロールバック機能
4. **テスタビリティ**: モック可能な設計
5. **クロスプラットフォーム**: Ubuntu以外への対応準備