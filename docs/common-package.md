# 共通パッケージ設計 (internal/common)

## 概要

全てのインストーラーで共通して使用される機能を提供するパッケージ群。

## パッケージ構成

### command.go - コマンド実行

```go
package common

import (
    "context"
    "os/exec"
)

// CommandRunner はコマンド実行の抽象化
type CommandRunner interface {
    // Run はコマンドを実行し、エラーを返す
    Run(ctx context.Context, name string, args ...string) error
    
    // RunWithOutput はコマンドを実行し、出力とエラーを返す
    RunWithOutput(ctx context.Context, name string, args ...string) (string, error)
    
    // CommandExists はコマンドが存在するかチェック
    CommandExists(name string) bool
}

// SystemCommandRunner は実際のシステムコマンドを実行
type SystemCommandRunner struct{}

func NewSystemCommandRunner() *SystemCommandRunner {
    return &SystemCommandRunner{}
}

func (r *SystemCommandRunner) Run(ctx context.Context, name string, args ...string) error {
    cmd := exec.CommandContext(ctx, name, args...)
    return cmd.Run()
}

func (r *SystemCommandRunner) RunWithOutput(ctx context.Context, name string, args ...string) (string, error) {
    cmd := exec.CommandContext(ctx, name, args...)
    out, err := cmd.Output()
    return string(out), err
}

func (r *SystemCommandRunner) CommandExists(name string) bool {
    _, err := exec.LookPath(name)
    return err == nil
}

// MockCommandRunner はテスト用のモック
type MockCommandRunner struct {
    Commands map[string]string // コマンド -> 期待する出力
    Errors   map[string]error  // コマンド -> 返すエラー
}

func NewMockCommandRunner() *MockCommandRunner {
    return &MockCommandRunner{
        Commands: make(map[string]string),
        Errors:   make(map[string]error),
    }
}
```

### logger.go - ログ出力

```go
package common

import (
    "fmt"
    "log/slog"
    "os"
)

// Logger はログ出力の抽象化
type Logger interface {
    Info(msg string, args ...any)
    Success(msg string, args ...any)
    Error(msg string, args ...any)
    Debug(msg string, args ...any)
}

// ColorLogger はカラー付きログ出力
type ColorLogger struct {
    logger *slog.Logger
}

const (
    ColorReset  = "\033[0m"
    ColorGreen  = "\033[32m"
    ColorBlue   = "\033[34m"
    ColorRed    = "\033[31m"
    ColorYellow = "\033[33m"
)

func NewColorLogger() *ColorLogger {
    handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    })
    return &ColorLogger{
        logger: slog.New(handler),
    }
}

func (l *ColorLogger) Info(msg string, args ...any) {
    fmt.Printf("%s[INFO]%s %s\n", ColorBlue, ColorReset, fmt.Sprintf(msg, args...))
}

func (l *ColorLogger) Success(msg string, args ...any) {
    fmt.Printf("%s[SUCCESS]%s %s\n", ColorGreen, ColorReset, fmt.Sprintf(msg, args...))
}

func (l *ColorLogger) Error(msg string, args ...any) {
    fmt.Printf("%s[ERROR]%s %s\n", ColorRed, ColorReset, fmt.Sprintf(msg, args...))
}

func (l *ColorLogger) Debug(msg string, args ...any) {
    l.logger.Debug(fmt.Sprintf(msg, args...))
}
```

### checker.go - 存在確認

```go
package common

import (
    "os"
    "path/filepath"
)

// Checker はファイル・ディレクトリ・コマンドの存在確認
type Checker interface {
    FileExists(path string) bool
    DirExists(path string) bool
    CommandExists(name string) bool
}

// SystemChecker は実際のシステムをチェック
type SystemChecker struct {
    cmdRunner CommandRunner
}

func NewSystemChecker(cmdRunner CommandRunner) *SystemChecker {
    return &SystemChecker{
        cmdRunner: cmdRunner,
    }
}

func (c *SystemChecker) FileExists(path string) bool {
    info, err := os.Stat(path)
    return err == nil && !info.IsDir()
}

func (c *SystemChecker) DirExists(path string) bool {
    info, err := os.Stat(path)
    return err == nil && info.IsDir()
}

func (c *SystemChecker) CommandExists(name string) bool {
    return c.cmdRunner.CommandExists(name)
}

// PathHelper はパス操作のヘルパー
type PathHelper struct{}

func NewPathHelper() *PathHelper {
    return &PathHelper{}
}

func (p *PathHelper) HomeDir() (string, error) {
    return os.UserHomeDir()
}

func (p *PathHelper) ConfigDir() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".config"), nil
}

func (p *PathHelper) EnsureDir(path string) error {
    return os.MkdirAll(path, 0755)
}
```

## 使用例

```go
// 依存性注入でテスタブルな構造
type SomeInstaller struct {
    cmdRunner CommandRunner
    logger    Logger
    checker   Checker
}

func NewSomeInstaller(cmdRunner CommandRunner, logger Logger, checker Checker) *SomeInstaller {
    return &SomeInstaller{
        cmdRunner: cmdRunner,
        logger:    logger,
        checker:   checker,
    }
}

func (i *SomeInstaller) Install(ctx context.Context) error {
    if i.checker.CommandExists("some-command") {
        i.logger.Info("some-command は既にインストール済み")
        return nil
    }
    
    i.logger.Info("some-command をインストール中...")
    if err := i.cmdRunner.Run(ctx, "apt-get", "install", "-y", "some-command"); err != nil {
        i.logger.Error("インストールに失敗: %v", err)
        return err
    }
    
    i.logger.Success("some-command をインストールしました")
    return nil
}
```

## テスト性

- インターフェースによる抽象化でモック可能
- 依存性注入によりテスト用の実装を差し込み可能
- 各機能が独立しており単体テストが容易