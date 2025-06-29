# メイン実行部設計

## 概要

アプリケーションのエントリーポイント（main.go）とセットアップ実行制御（internal/setup）の設計。

## main.go - エントリーポイント

```go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "setup-go/internal/common"
    "setup-go/internal/installer"
    "setup-go/internal/setup"
    "time"
)

func main() {
    // コマンドライン引数の解析
    var (
        verbose = flag.Bool("verbose", false, "詳細なログを出力")
        dryRun  = flag.Bool("dry-run", false, "実際にはインストールせず、実行内容のみ表示")
        timeout = flag.Duration("timeout", 10*time.Minute, "全体のタイムアウト時間")
        help    = flag.Bool("help", false, "ヘルプを表示")
    )
    flag.Parse()

    if *help {
        printHelp()
        return
    }

    // ログレベル設定
    logger := common.NewColorLogger()
    if *verbose {
        logger.Debug("詳細ログモードが有効です")
    }

    // 依存関係の初期化
    cmdRunner := common.NewSystemCommandRunner()
    checker := common.NewSystemChecker(cmdRunner)
    pathHelper := common.NewPathHelper()

    // ドライランモードの場合はモックを使用
    if *dryRun {
        logger.Info("ドライランモードで実行します")
        cmdRunner = common.NewMockCommandRunner()
    }

    // セットアップランナーの初期化
    runner := setup.NewRunner(logger, cmdRunner, checker, pathHelper)

    // コンテキストとタイムアウト設定
    ctx, cancel := context.WithTimeout(context.Background(), *timeout)
    defer cancel()

    // セットアップ実行
    if err := runner.Run(ctx); err != nil {
        logger.Error("セットアップに失敗しました: %v", err)
        os.Exit(1)
    }

    logger.Success("セットアップが完了しました！")
    printPostInstallInstructions()
}

func printHelp() {
    fmt.Println("Claude Code セットアップツール")
    fmt.Println("")
    fmt.Println("使用方法:")
    fmt.Println("  setup-go [オプション]")
    fmt.Println("")
    fmt.Println("オプション:")
    flag.PrintDefaults()
    fmt.Println("")
    fmt.Println("例:")
    fmt.Println("  setup-go                    # 通常のセットアップ")
    fmt.Println("  setup-go --verbose         # 詳細ログ付きセットアップ")
    fmt.Println("  setup-go --dry-run         # ドライラン（実際にはインストールしない）")
    fmt.Println("  setup-go --timeout 5m      # 5分でタイムアウト")
}

func printPostInstallInstructions() {
    fmt.Println("")
    fmt.Println("=== セットアップ完了 ===")
    fmt.Println("")
    fmt.Println("次の手順を実行してください:")
    fmt.Println("")
    fmt.Println("1. Claude Code API キーの設定:")
    fmt.Println("   export ANTHROPIC_API_KEY='your-api-key'")
    fmt.Println("")
    fmt.Println("2. PATHの設定（必要に応じて）:")
    fmt.Println("   export PATH=\"$PATH:/opt/nvim-linux-x86_64/bin\"")
    fmt.Println("")
    fmt.Println("3. 新しいシェルセッションを開始:")
    fmt.Println("   exec $SHELL")
    fmt.Println("")
    fmt.Println("4. ツールの使用:")
    fmt.Println("   claude-code --help")
    fmt.Println("   nvim")
    fmt.Println("")
    fmt.Println("注意: デフォルトシェルの変更は再ログイン後に有効になります")
}
```

## internal/setup/runner.go - セットアップ実行制御

```go
package setup

import (
    "context"
    "fmt"
    "setup-go/internal/common"
    "setup-go/internal/installer"
    "sync"
    "time"
)

// Runner はセットアップ全体の実行を制御
type Runner struct {
    logger     common.Logger
    cmdRunner  common.CommandRunner
    checker    common.Checker
    pathHelper *common.PathHelper
    installers []installer.Installer
}

// NewRunner は新しいRunnerインスタンスを作成
func NewRunner(
    logger common.Logger,
    cmdRunner common.CommandRunner,
    checker common.Checker,
    pathHelper *common.PathHelper,
) *Runner {
    return &Runner{
        logger:     logger,
        cmdRunner:  cmdRunner,
        checker:    checker,
        pathHelper: pathHelper,
    }
}

// Run はセットアップを実行
func (r *Runner) Run(ctx context.Context) error {
    r.logger.Info("Claude Code セットアップを開始...")
    
    // インストーラーの初期化
    if err := r.initializeInstallers(); err != nil {
        return fmt.Errorf("インストーラーの初期化に失敗: %w", err)
    }
    
    // 事前チェック
    if err := r.preInstallCheck(ctx); err != nil {
        return fmt.Errorf("事前チェックに失敗: %w", err)
    }
    
    // インストール実行
    if err := r.runInstallers(ctx); err != nil {
        return fmt.Errorf("インストールに失敗: %w", err)
    }
    
    // 事後チェック
    if err := r.postInstallCheck(ctx); err != nil {
        r.logger.Error("事後チェックでエラーが発見されました: %v", err)
        // 事後チェックは警告として扱い、エラーは返さない
    }
    
    return nil
}

// initializeInstallers はインストーラーを初期化
func (r *Runner) initializeInstallers() error {
    r.installers = []installer.Installer{
        installer.NewConfigCloner(r.cmdRunner, r.logger, r.checker, r.pathHelper),
        installer.NewNodeJSInstaller(r.cmdRunner, r.logger, r.checker),
        installer.NewClaudeCodeInstaller(r.cmdRunner, r.logger, r.checker),
        installer.NewNeovimInstaller(r.cmdRunner, r.logger, r.checker, r.pathHelper),
        installer.NewFishInstaller(r.cmdRunner, r.logger, r.checker),
    }
    
    r.logger.Debug("インストーラーを初期化しました: %d個", len(r.installers))
    return nil
}

// preInstallCheck は事前チェックを実行
func (r *Runner) preInstallCheck(ctx context.Context) error {
    r.logger.Info("事前チェックを実行中...")
    
    // システム要件チェック
    if err := r.checkSystemRequirements(ctx); err != nil {
        return err
    }
    
    // ネットワーク接続チェック
    if err := r.checkNetworkConnectivity(ctx); err != nil {
        return err
    }
    
    // 権限チェック
    if err := r.checkPermissions(ctx); err != nil {
        return err
    }
    
    r.logger.Success("事前チェックが完了しました")
    return nil
}

// runInstallers はインストーラーを順次実行
func (r *Runner) runInstallers(ctx context.Context) error {
    for i, inst := range r.installers {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        
        r.logger.Info("(%d/%d) %s のセットアップを開始...", i+1, len(r.installers), inst.Name())
        
        start := time.Now()
        if err := inst.Install(ctx); err != nil {
            return fmt.Errorf("%s のインストールに失敗: %w", inst.Name(), err)
        }
        
        duration := time.Since(start)
        r.logger.Debug("%s のセットアップが完了しました (所要時間: %v)", inst.Name(), duration)
    }
    
    return nil
}

// runInstallersParallel はインストーラーを並列実行（将来的な拡張用）
func (r *Runner) runInstallersParallel(ctx context.Context) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(r.installers))
    
    for _, inst := range r.installers {
        wg.Add(1)
        go func(installer installer.Installer) {
            defer wg.Done()
            if err := installer.Install(ctx); err != nil {
                errChan <- fmt.Errorf("%s: %w", installer.Name(), err)
            }
        }(inst)
    }
    
    wg.Wait()
    close(errChan)
    
    // エラーがあれば最初のエラーを返す
    if err := <-errChan; err != nil {
        return err
    }
    
    return nil
}

// postInstallCheck は事後チェックを実行
func (r *Runner) postInstallCheck(ctx context.Context) error {
    r.logger.Info("事後チェックを実行中...")
    
    var errors []error
    
    for _, inst := range r.installers {
        if installed, err := inst.IsInstalled(ctx); err != nil {
            errors = append(errors, fmt.Errorf("%s の確認に失敗: %w", inst.Name(), err))
        } else if !installed {
            errors = append(errors, fmt.Errorf("%s がインストールされていません", inst.Name()))
        } else {
            if version, err := inst.Version(ctx); err == nil {
                r.logger.Info("%s: %s", inst.Name(), version)
            } else {
                r.logger.Info("%s: インストール済み", inst.Name())
            }
        }
    }
    
    if len(errors) > 0 {
        for _, err := range errors {
            r.logger.Error("%v", err)
        }
        return fmt.Errorf("%d個のエラーが見つかりました", len(errors))
    }
    
    r.logger.Success("事後チェックが完了しました")
    return nil
}

// checkSystemRequirements はシステム要件をチェック
func (r *Runner) checkSystemRequirements(ctx context.Context) error {
    // Ubuntu/Debian系であることを確認
    if !r.checker.FileExists("/etc/debian_version") {
        return fmt.Errorf("このツールはUbuntu/Debian系でのみ動作します")
    }
    
    // 必要なコマンドの存在確認
    requiredCommands := []string{"curl", "git", "sudo", "tar"}
    for _, cmd := range requiredCommands {
        if !r.checker.CommandExists(cmd) {
            return fmt.Errorf("必要なコマンドが見つかりません: %s", cmd)
        }
    }
    
    r.logger.Debug("システム要件チェックが完了しました")
    return nil
}

// checkNetworkConnectivity はネットワーク接続をチェック
func (r *Runner) checkNetworkConnectivity(ctx context.Context) error {
    testURLs := []string{
        "https://deb.nodesource.com",
        "https://github.com",
        "https://registry.npmjs.org",
    }
    
    for _, url := range testURLs {
        if err := r.cmdRunner.Run(ctx, "curl", "-sf", "--max-time", "10", url); err != nil {
            return fmt.Errorf("ネットワーク接続に問題があります: %s", url)
        }
    }
    
    r.logger.Debug("ネットワーク接続チェックが完了しました")
    return nil
}

// checkPermissions は権限をチェック
func (r *Runner) checkPermissions(ctx context.Context) error {
    // sudo権限の確認
    if err := r.cmdRunner.Run(ctx, "sudo", "-n", "true"); err != nil {
        return fmt.Errorf("sudo権限が必要です。パスワードなしでsudoを実行できるように設定してください")
    }
    
    // 書き込み権限の確認
    homeDir, err := r.pathHelper.HomeDir()
    if err != nil {
        return fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
    }
    
    if !r.checker.DirExists(homeDir) {
        return fmt.Errorf("ホームディレクトリが存在しません: %s", homeDir)
    }
    
    r.logger.Debug("権限チェックが完了しました")
    return nil
}
```

## go.mod - Go モジュール定義

```go
module setup-go

go 1.21

require (
    // 標準ライブラリのみ使用
    // 必要に応じて外部ライブラリを追加
)
```

## 特徴

1. **堅牢なエラーハンドリング**: 各段階でエラーを適切に処理
2. **タイムアウト制御**: 全体およびコマンド単位でのタイムアウト
3. **ドライランモード**: 実際にインストールせずに動作確認
4. **詳細ログ**: デバッグ用の詳細ログ出力
5. **事前・事後チェック**: インストール前後の環境確認
6. **並列実行対応**: 将来的な並列インストール対応
7. **ユーザーフレンドリー**: 分かりやすいヘルプと完了メッセージ

## 使用例

```bash
# 通常のセットアップ
go run main.go

# 詳細ログ付き
go run main.go --verbose

# ドライラン
go run main.go --dry-run

# カスタムタイムアウト
go run main.go --timeout 15m
```