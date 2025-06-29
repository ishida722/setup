package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Logger provides colored logging functions
type Logger struct {
	info    *color.Color
	success *color.Color
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		info:    color.New(color.FgBlue),
		success: color.New(color.FgGreen),
	}
}

// Log prints an info message
func (l *Logger) Log(message string) {
	l.info.Print("[INFO] ")
	fmt.Println(message)
}

// Success prints a success message
func (l *Logger) Success(message string) {
	l.success.Print("[SUCCESS] ")
	fmt.Println(message)
}

// Setup handles the main setup process
type Setup struct {
	logger *Logger
}

// NewSetup creates a new setup instance
func NewSetup() *Setup {
	return &Setup{
		logger: NewLogger(),
	}
}

// commandExists checks if a command exists in PATH
func (s *Setup) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// runCommand executes a command and returns its output
func (s *Setup) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// getCommandOutput executes a command and returns its output as string
func (s *Setup) getCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

// CloneConfigs clones configuration files from GitHub
func (s *Setup) CloneConfigs() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	
	configDir := filepath.Join(homeDir, ".config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Clone Neovim config
	nvimDir := filepath.Join(configDir, "nvim")
	if _, err := os.Stat(nvimDir); os.IsNotExist(err) {
		s.logger.Log("Neovim設定をクローン中...")
		if err := s.runCommand("git", "clone", "https://github.com/ishida722/nvim", nvimDir); err != nil {
			return err
		}
		s.logger.Success("Neovim設定をクローンしました")
	} else {
		s.logger.Log("Neovim設定は既に存在します")
	}

	// Clone Fish config
	fishDir := filepath.Join(configDir, "fish")
	if _, err := os.Stat(fishDir); os.IsNotExist(err) {
		s.logger.Log("Fish設定をクローン中...")
		if err := s.runCommand("git", "clone", "https://github.com/ishida722/fish", fishDir); err != nil {
			return err
		}
		s.logger.Success("Fish設定をクローンしました")
	} else {
		s.logger.Log("Fish設定は既に存在します")
	}

	return nil
}

// InstallNodeJS installs Node.js LTS
func (s *Setup) InstallNodeJS() error {
	if s.commandExists("node") {
		version, _ := s.getCommandOutput("node", "--version")
		s.logger.Log(fmt.Sprintf("Node.js は既にインストール済み: %s", version))
		return nil
	}

	s.logger.Log("Node.jsをインストール中...")
	
	// Download and execute NodeSource setup script
	cmd := exec.Command("bash", "-c", "curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Install nodejs
	if err := s.runCommand("sudo", "apt-get", "install", "-y", "nodejs"); err != nil {
		return err
	}

	version, _ := s.getCommandOutput("node", "--version")
	s.logger.Success(fmt.Sprintf("Node.jsをインストールしました: %s", version))
	return nil
}

// InstallClaudeCode installs Claude Code globally via npm
func (s *Setup) InstallClaudeCode() error {
	if s.commandExists("claude-code") {
		s.logger.Log("Claude Code は既にインストール済み")
		return nil
	}

	s.logger.Log("Claude Codeをインストール中...")
	if err := s.runCommand("npm", "install", "-g", "@anthropic-ai/claude-code"); err != nil {
		return err
	}

	s.logger.Success("Claude Codeをインストールしました")
	return nil
}

// InstallNeovim installs the latest Neovim from GitHub releases
func (s *Setup) InstallNeovim() error {
	if s.commandExists("nvim") {
		output, _ := s.getCommandOutput("nvim", "--version")
		version := strings.Split(output, "\n")[0]
		s.logger.Log(fmt.Sprintf("Neovim は既にインストール済み: %s", version))
		return nil
	}

	s.logger.Log("Neovimをインストール中...")
	
	// Download Neovim
	if err := s.runCommand("curl", "-LO", "https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz"); err != nil {
		return err
	}

	// Remove existing installation
	if err := s.runCommand("sudo", "rm", "-rf", "/opt/nvim"); err != nil {
		// Ignore error if directory doesn't exist
	}

	// Extract to /opt
	if err := s.runCommand("sudo", "tar", "-C", "/opt", "-xzf", "nvim-linux-x86_64.tar.gz"); err != nil {
		return err
	}

	// Clean up downloaded file
	os.Remove("nvim-linux-x86_64.tar.gz")

	s.logger.Success("Neovimをインストールしました")
	return nil
}

// InstallFish installs Fish shell and sets it as default
func (s *Setup) InstallFish() error {
	if s.commandExists("fish") {
		version, _ := s.getCommandOutput("fish", "--version")
		s.logger.Log(fmt.Sprintf("Fish は既にインストール済み: %s", version))
	} else {
		s.logger.Log("Fishをインストール中...")
		if err := s.runCommand("sudo", "apt-get", "update"); err != nil {
			return err
		}
		if err := s.runCommand("sudo", "apt-get", "install", "-y", "fish"); err != nil {
			return err
		}
		s.logger.Success("Fishをインストールしました")
	}

	// Set Fish as default shell
	fishPath, err := exec.LookPath("fish")
	if err != nil {
		return err
	}

	currentShell := os.Getenv("SHELL")
	if currentShell != fishPath {
		s.logger.Log("デフォルトシェルをfishに変更中...")
		if err := s.runCommand("chsh", "-s", fishPath); err != nil {
			return err
		}
		s.logger.Success("デフォルトシェルをfishに変更しました（再ログイン後に有効）")
	} else {
		s.logger.Log("デフォルトシェルは既にfishです")
	}

	return nil
}

// Run executes the main setup process
func (s *Setup) Run() error {
	s.logger.Log("Claude Code セットアップを開始...")

	if err := s.CloneConfigs(); err != nil {
		return fmt.Errorf("設定ファイルのクローンに失敗: %w", err)
	}

	if err := s.InstallNodeJS(); err != nil {
		return fmt.Errorf("Node.jsのインストールに失敗: %w", err)
	}

	if err := s.InstallClaudeCode(); err != nil {
		return fmt.Errorf("Claude Codeのインストールに失敗: %w", err)
	}

	if err := s.InstallNeovim(); err != nil {
		return fmt.Errorf("Neovimのインストールに失敗: %w", err)
	}

	if err := s.InstallFish(); err != nil {
		return fmt.Errorf("Fishのインストールに失敗: %w", err)
	}

	fmt.Println()
	s.logger.Success("セットアップ完了！")
	fmt.Println("使用方法: claude-code --help")
	fmt.Println("API キー設定: export ANTHROPIC_API_KEY='your-key'")
	fmt.Println("Neovim: nvim")
	fmt.Println("注意: デフォルトシェルの変更は再ログイン後に有効になります")

	return nil
}

func main() {
	setup := NewSetup()
	if err := setup.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}