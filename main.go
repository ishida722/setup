package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"claude-setup/installer"
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

// InstallNodeJS installs Node.js LTS using the new installer package
func (s *Setup) InstallNodeJS() error {
	nodeCmd := installer.InstallCommand{
		CheckCommands:   []string{"node --version", "nodejs --version"},
		InstallCommands: []string{
			"curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -",
			"sudo apt-get install -y nodejs",
		},
		Name:        "Node.js",
		Description: "JavaScript runtime",
	}
	
	result := nodeCmd.Install()
	if result.Err != nil {
		return result.Err
	}
	return nil
}

// InstallClaudeCode installs Claude Code using the new installer package
func (s *Setup) InstallClaudeCode() error {
	claudeCmd := installer.InstallCommand{
		CheckCommands:   []string{"claude-code --version"},
		InstallCommands: []string{"npm install -g @anthropic-ai/claude-code"},
		Name:            "Claude Code",
		Description:     "AI-powered coding assistant",
	}
	
	result := claudeCmd.Install()
	if result.Err != nil {
		return result.Err
	}
	return nil
}

// InstallNeovim installs Neovim using the new installer package
func (s *Setup) InstallNeovim() error {
	neovimCmd := installer.InstallCommand{
		CheckCommands: []string{"nvim --version"},
		InstallCommands: []string{
			"curl -LO https://github.com/neovim/neovim/releases/latest/download/nvim-linux-x86_64.tar.gz",
		},
		InstallFunc: func() error {
			// Custom installation logic for Neovim
			if err := exec.Command("sudo", "rm", "-rf", "/opt/nvim").Run(); err != nil {
				// Ignore error if directory doesn't exist
			}
			
			if err := exec.Command("sudo", "tar", "-C", "/opt", "-xzf", "nvim-linux-x86_64.tar.gz").Run(); err != nil {
				return err
			}
			
			// Clean up downloaded file
			os.Remove("nvim-linux-x86_64.tar.gz")
			return nil
		},
		Name:        "Neovim",
		Description: "Modern Vim-based text editor",
	}
	
	result := neovimCmd.Install()
	if result.Err != nil {
		return result.Err
	}
	return nil
}

// InstallFish installs Fish shell using the new installer package
func (s *Setup) InstallFish() error {
	fishCmd := installer.InstallCommand{
		CheckCommands:   []string{"fish --version"},
		InstallCommands: []string{
			"sudo apt-get update",
			"sudo apt-get install -y fish",
		},
		InstallFunc: func() error {
			// Set Fish as default shell
			fishPath, err := exec.LookPath("fish")
			if err != nil {
				return err
			}

			currentShell := os.Getenv("SHELL")
			if currentShell != fishPath {
				logger := installer.NewLogger()
				logger.Log("デフォルトシェルをfishに変更中...")
				if err := exec.Command("chsh", "-s", fishPath).Run(); err != nil {
					return err
				}
				logger.Success("デフォルトシェルをfishに変更しました（再ログイン後に有効）")
			}
			return nil
		},
		Name:        "Fish Shell",
		Description: "User-friendly command line shell",
	}
	
	result := fishCmd.Install()
	if result.Err != nil {
		return result.Err
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
	// Example usage of the new installer package
	fmt.Println("=== Testing New Installer Package ===")
	
	nodeCmd := installer.InstallCommand{
		CheckCommands:   []string{"node --version", "nodejs --version"},
		InstallCommands: []string{
			"curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -",
			"sudo apt-get install -y nodejs",
		},
		Name:        "Node.js",
		Description: "JavaScript runtime",
	}
	
	claudeCmd := installer.InstallCommand{
		CheckCommands:   []string{"claude-code --version"},
		InstallCommands: []string{"npm install -g @anthropic-ai/claude-code"},
		Name:            "Claude Code",
		Description:     "AI-powered coding assistant",
	}
	
	// Test CheckInstalled method
	fmt.Println("\n--- Testing CheckInstalled method ---")
	nodeInstalled, nodeVersion, err := nodeCmd.CheckInstalled()
	if err != nil {
		fmt.Printf("Node.js check error: %v\n", err)
	} else {
		fmt.Printf("Node.js installed: %v, version: %s\n", nodeInstalled, nodeVersion)
	}
	
	claudeInstalled, claudeVersion, err := claudeCmd.CheckInstalled()
	if err != nil {
		fmt.Printf("Claude Code check error: %v\n", err)
	} else {
		fmt.Printf("Claude Code installed: %v, version: %s\n", claudeInstalled, claudeVersion)
	}
	
	// Test Install method
	fmt.Println("\n--- Testing Install method ---")
	nodeResult := nodeCmd.Install()
	claudeResult := claudeCmd.Install()
	
	fmt.Printf("Node.js Result: Success=%v, AlreadyInstalled=%v, Version=%s\n", 
		nodeResult.Success, nodeResult.AlreadyInstalled, nodeResult.Version)
	fmt.Printf("Claude Code Result: Success=%v, AlreadyInstalled=%v, Version=%s\n", 
		claudeResult.Success, claudeResult.AlreadyInstalled, claudeResult.Version)
	
	fmt.Println("\n=== Using Setup with Installer Package ===")
	setup := NewSetup()
	if err := setup.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}