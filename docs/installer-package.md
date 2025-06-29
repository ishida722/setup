# インストーラーパッケージ設計 (internal/installer)

## 概要

各ソフトウェアのインストール処理を個別のパッケージとして実装し、共通インターフェースで統一する。

## 共通インターフェース (installer.go)

```go
package installer

import (
    "context"
    "setup-go/internal/common"
)

// Installer は全てのインストーラーが実装するインターフェース
type Installer interface {
    // Name はインストーラーの名前を返す
    Name() string
    
    // IsInstalled はソフトウェアがインストール済みかチェック
    IsInstalled(ctx context.Context) (bool, error)
    
    // Install はソフトウェアをインストール
    Install(ctx context.Context) error
    
    // Version はインストール済みのバージョンを返す（可能な場合）
    Version(ctx context.Context) (string, error)
}

// BaseInstaller は共通機能を提供
type BaseInstaller struct {
    name      string
    cmdRunner common.CommandRunner
    logger    common.Logger
    checker   common.Checker
}

func NewBaseInstaller(
    name string,
    cmdRunner common.CommandRunner,
    logger common.Logger,
    checker common.Checker,
) *BaseInstaller {
    return &BaseInstaller{
        name:      name,
        cmdRunner: cmdRunner,
        logger:    logger,
        checker:   checker,
    }
}

func (b *BaseInstaller) Name() string {
    return b.name
}
```

## Node.js インストーラー (nodejs.go)

```go
package installer

import (
    "context"
    "fmt"
    "setup-go/internal/common"
    "strings"
)

// NodeJSInstaller はNode.jsのインストールを担当
type NodeJSInstaller struct {
    *BaseInstaller
}

func NewNodeJSInstaller(
    cmdRunner common.CommandRunner,
    logger common.Logger,
    checker common.Checker,
) *NodeJSInstaller {
    return &NodeJSInstaller{
        BaseInstaller: NewBaseInstaller("Node.js", cmdRunner, logger, checker),
    }
}

func (n *NodeJSInstaller) IsInstalled(ctx context.Context) (bool, error) {
    return n.checker.CommandExists("node"), nil
}

func (n *NodeJSInstaller) Version(ctx context.Context) (string, error) {
    if installed, _ := n.IsInstalled(ctx); !installed {
        return "", fmt.Errorf("Node.js is not installed")
    }
    
    output, err := n.cmdRunner.RunWithOutput(ctx, "node", "--version")
    if err != nil {
        return "", err
    }
    
    return strings.TrimSpace(output), nil
}

func (n *NodeJSInstaller) Install(ctx context.Context) error {
    if installed, err := n.IsInstalled(ctx); err != nil {
        return err
    } else if installed {
        version, _ := n.Version(ctx)
        n.logger.Info("Node.js は既にインストール済み: %s", version)
        return nil
    }
    
    n.logger.Info("Node.jsをインストール中...")
    
    // NodeSourceリポジトリセットアップ
    if err := n.cmdRunner.Run(ctx, "curl", "-fsSL", 
        "https://deb.nodesource.com/setup_lts.x"); err != nil {
        return fmt.Errorf("NodeSource setup failed: %w", err)
    }
    
    // Node.jsインストール
    if err := n.cmdRunner.Run(ctx, "sudo", "apt-get", "install", "-y", "nodejs"); err != nil {
        return fmt.Errorf("Node.js installation failed: %w", err)
    }
    
    version, _ := n.Version(ctx)
    n.logger.Success("Node.jsをインストールしました: %s", version)
    
    return nil
}
```

## Claude Code インストーラー (claude.go)

```go
package installer

import (
    "context"
    "fmt"
    "setup-go/internal/common"
)

// ClaudeCodeInstaller はClaude Codeのインストールを担当
type ClaudeCodeInstaller struct {
    *BaseInstaller
}

func NewClaudeCodeInstaller(
    cmdRunner common.CommandRunner,
    logger common.Logger,
    checker common.Checker,
) *ClaudeCodeInstaller {
    return &ClaudeCodeInstaller{
        BaseInstaller: NewBaseInstaller("Claude Code", cmdRunner, logger, checker),
    }
}

func (c *ClaudeCodeInstaller) IsInstalled(ctx context.Context) (bool, error) {
    return c.checker.CommandExists("claude-code"), nil
}

func (c *ClaudeCodeInstaller) Version(ctx context.Context) (string, error) {
    if installed, _ := c.IsInstalled(ctx); !installed {
        return "", fmt.Errorf("Claude Code is not installed")
    }
    
    // claude-codeはバージョン取得コマンドが異なる場合があるため適宜調整
    output, err := c.cmdRunner.RunWithOutput(ctx, "claude-code", "--version")
    if err != nil {
        return "", err
    }
    
    return strings.TrimSpace(output), nil
}

func (c *ClaudeCodeInstaller) Install(ctx context.Context) error {
    if installed, err := c.IsInstalled(ctx); err != nil {
        return err
    } else if installed {
        c.logger.Info("Claude Code は既にインストール済み")
        return nil
    }
    
    c.logger.Info("Claude Codeをインストール中...")
    
    if err := c.cmdRunner.Run(ctx, "npm", "install", "-g", "@anthropic-ai/claude-code"); err != nil {
        return fmt.Errorf("Claude Code installation failed: %w", err)
    }
    
    c.logger.Success("Claude Codeをインストールしました")
    
    return nil
}
```

## Neovim インストーラー (neovim.go)

```go
package installer

import (
    "context"
    "fmt"
    "setup-go/internal/common"
    "strings"
)

// NeovimInstaller はNeovimのインストールを担当
type NeovimInstaller struct {
    *BaseInstaller
    pathHelper *common.PathHelper
}

func NewNeovimInstaller(
    cmdRunner common.CommandRunner,
    logger common.Logger,
    checker common.Checker,
    pathHelper *common.PathHelper,
) *NeovimInstaller {
    return &NeovimInstaller{
        BaseInstaller: NewBaseInstaller("Neovim", cmdRunner, logger, checker),
        pathHelper:    pathHelper,
    }
}

func (n *NeovimInstaller) IsInstalled(ctx context.Context) (bool, error) {
    return n.checker.CommandExists("nvim"), nil
}

func (n *NeovimInstaller) Version(ctx context.Context) (string, error) {
    if installed, _ := n.IsInstalled(ctx); !installed {
        return "", fmt.Errorf("Neovim is not installed")
    }
    
    output, err := n.cmdRunner.RunWithOutput(ctx, "nvim", "--version")
    if err != nil {
        return "", err
    }
    
    lines := strings.Split(output, "\n")
    if len(lines) > 0 {
        return strings.TrimSpace(lines[0]), nil
    }
    
    return "", fmt.Errorf("could not parse version")
}

func (n *NeovimInstaller) Install(ctx context.Context) error {
    if installed, err := n.IsInstalled(ctx); err != nil {
        return err
    } else if installed {
        version, _ := n.Version(ctx)
        n.logger.Info("Neovim は既にインストール済み: %s", version)
        return nil
    }
    
    n.logger.Info("Neovimをインストール中...")
    
    // 最新版をダウンロード
    if err := n.cmdRunner.Run(ctx, "curl", "-LO", 
        "https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz"); err != nil {
        return fmt.Errorf("Neovim download failed: %w", err)
    }
    
    // 既存インストールを削除
    if err := n.cmdRunner.Run(ctx, "sudo", "rm", "-rf", "/opt/nvim"); err != nil {
        n.logger.Debug("Failed to remove existing Neovim: %v", err)
    }
    
    // 解凍とインストール
    if err := n.cmdRunner.Run(ctx, "sudo", "tar", "-C", "/opt", "-xzf", 
        "nvim-linux-x86_64.tar.gz"); err != nil {
        return fmt.Errorf("Neovim extraction failed: %w", err)
    }
    
    // 一時ファイル削除
    if err := n.cmdRunner.Run(ctx, "rm", "-f", "nvim-linux-x86_64.tar.gz"); err != nil {
        n.logger.Debug("Failed to remove temporary file: %v", err)
    }
    
    n.logger.Success("Neovimをインストールしました")
    
    return nil
}
```

## Fish シェル インストーラー (fish.go)

```go
package installer

import (
    "context"
    "fmt"
    "setup-go/internal/common"
    "strings"
)

// FishInstaller はFishシェルのインストールを担当
type FishInstaller struct {
    *BaseInstaller
}

func NewFishInstaller(
    cmdRunner common.CommandRunner,
    logger common.Logger,
    checker common.Checker,
) *FishInstaller {
    return &FishInstaller{
        BaseInstaller: NewBaseInstaller("Fish Shell", cmdRunner, logger, checker),
    }
}

func (f *FishInstaller) IsInstalled(ctx context.Context) (bool, error) {
    return f.checker.CommandExists("fish"), nil
}

func (f *FishInstaller) Version(ctx context.Context) (string, error) {
    if installed, _ := f.IsInstalled(ctx); !installed {
        return "", fmt.Errorf("Fish is not installed")
    }
    
    output, err := f.cmdRunner.RunWithOutput(ctx, "fish", "--version")
    if err != nil {
        return "", err
    }
    
    return strings.TrimSpace(output), nil
}

func (f *FishInstaller) Install(ctx context.Context) error {
    if installed, err := f.IsInstalled(ctx); err != nil {
        return err
    } else if installed {
        version, _ := f.Version(ctx)
        f.logger.Info("Fish は既にインストール済み: %s", version)
    } else {
        f.logger.Info("Fishをインストール中...")
        
        // パッケージリスト更新
        if err := f.cmdRunner.Run(ctx, "sudo", "apt-get", "update"); err != nil {
            return fmt.Errorf("apt-get update failed: %w", err)
        }
        
        // Fishインストール
        if err := f.cmdRunner.Run(ctx, "sudo", "apt-get", "install", "-y", "fish"); err != nil {
            return fmt.Errorf("Fish installation failed: %w", err)
        }
        
        f.logger.Success("Fishをインストールしました")
    }
    
    // デフォルトシェル変更
    return f.setDefaultShell(ctx)
}

func (f *FishInstaller) setDefaultShell(ctx context.Context) error {
    fishPath, err := f.cmdRunner.RunWithOutput(ctx, "which", "fish")
    if err != nil {
        return fmt.Errorf("could not find fish path: %w", err)
    }
    fishPath = strings.TrimSpace(fishPath)
    
    currentShell := os.Getenv("SHELL")
    if currentShell == fishPath {
        f.logger.Info("デフォルトシェルは既にfishです")
        return nil
    }
    
    f.logger.Info("デフォルトシェルをfishに変更中...")
    if err := f.cmdRunner.Run(ctx, "chsh", "-s", fishPath); err != nil {
        return fmt.Errorf("failed to change default shell: %w", err)
    }
    
    f.logger.Success("デフォルトシェルをfishに変更しました（再ログイン後に有効）")
    
    return nil
}
```

## 設定ファイル クローン (config.go)

```go
package installer

import (
    "context"
    "fmt"
    "path/filepath"
    "setup-go/internal/common"
)

// ConfigCloner は設定ファイルのクローンを担当
type ConfigCloner struct {
    *BaseInstaller
    pathHelper *common.PathHelper
}

func NewConfigCloner(
    cmdRunner common.CommandRunner,
    logger common.Logger,
    checker common.Checker,
    pathHelper *common.PathHelper,
) *ConfigCloner {
    return &ConfigCloner{
        BaseInstaller: NewBaseInstaller("Config Files", cmdRunner, logger, checker),
        pathHelper:    pathHelper,
    }
}

func (c *ConfigCloner) IsInstalled(ctx context.Context) (bool, error) {
    configDir, err := c.pathHelper.ConfigDir()
    if err != nil {
        return false, err
    }
    
    nvimConfigExists := c.checker.DirExists(filepath.Join(configDir, "nvim"))
    fishConfigExists := c.checker.DirExists(filepath.Join(configDir, "fish"))
    
    return nvimConfigExists && fishConfigExists, nil
}

func (c *ConfigCloner) Version(ctx context.Context) (string, error) {
    return "latest", nil // 設定ファイルはバージョン管理なし
}

func (c *ConfigCloner) Install(ctx context.Context) error {
    configDir, err := c.pathHelper.ConfigDir()
    if err != nil {
        return fmt.Errorf("could not get config directory: %w", err)
    }
    
    // .configディレクトリ作成
    if err := c.pathHelper.EnsureDir(configDir); err != nil {
        return fmt.Errorf("could not create config directory: %w", err)
    }
    
    // Neovim設定クローン
    if err := c.cloneNeovimConfig(ctx, configDir); err != nil {
        return err
    }
    
    // Fish設定クローン
    if err := c.cloneFishConfig(ctx, configDir); err != nil {
        return err
    }
    
    return nil
}

func (c *ConfigCloner) cloneNeovimConfig(ctx context.Context, configDir string) error {
    nvimConfigPath := filepath.Join(configDir, "nvim")
    
    if c.checker.DirExists(nvimConfigPath) {
        c.logger.Info("Neovim設定は既に存在します")
        return nil
    }
    
    c.logger.Info("Neovim設定をクローン中...")
    if err := c.cmdRunner.Run(ctx, "git", "clone", 
        "https://github.com/ishida722/nvim", nvimConfigPath); err != nil {
        return fmt.Errorf("Neovim config clone failed: %w", err)
    }
    
    c.logger.Success("Neovim設定をクローンしました")
    return nil
}

func (c *ConfigCloner) cloneFishConfig(ctx context.Context, configDir string) error {
    fishConfigPath := filepath.Join(configDir, "fish")
    
    if c.checker.DirExists(fishConfigPath) {
        c.logger.Info("Fish設定は既に存在します")
        return nil
    }
    
    c.logger.Info("Fish設定をクローン中...")
    if err := c.cmdRunner.Run(ctx, "git", "clone", 
        "https://github.com/ishida722/fish", fishConfigPath); err != nil {
        return fmt.Errorf("Fish config clone failed: %w", err)
    }
    
    c.logger.Success("Fish設定をクローンしました")
    return nil
}
```

## 設計の利点

1. **統一インターフェース**: 全てのインストーラーが同じメソッドを持つ
2. **独立性**: 各インストーラーは他に依存しない
3. **テスト性**: モックを使った単体テストが容易
4. **拡張性**: 新しいソフトウェアのインストーラーを簡単に追加可能
5. **エラーハンドリング**: 詳細なエラー情報を提供