package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// InstallCommand defines installation command configuration
type InstallCommand struct {
	CheckCommands   []string    // List of commands to check if installed (e.g. ["node --version"])
	InstallCommands []string    // List of commands to execute for installation
	InstallFunc     func() error // Custom installation function (default: does nothing)
	Name            string      // Software name (e.g. "Node.js")
	Description     string      // Description
}

// InstallResult represents the result of an installation
type InstallResult struct {
	AlreadyInstalled bool   // Whether already installed
	Success          bool   // Whether installation succeeded
	Version          string // Installed version
	Err              error  // Error information
}

// Logger provides colored logging functions
type Logger struct {
	info    *color.Color
	success *color.Color
	errorColor *color.Color
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		info:       color.New(color.FgBlue),
		success:    color.New(color.FgGreen),
		errorColor: color.New(color.FgRed),
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

// Error prints an error message
func (l *Logger) Error(message string) {
	l.errorColor.Print("[ERROR] ")
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

// defaultInstallFunc is the default installation function that does nothing
func defaultInstallFunc() error {
	return nil
}

// checkInstalled checks if software is installed using check commands
func checkInstalled(checkCmds []string) (bool, string, error) {
	for _, cmdStr := range checkCmds {
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}
		
		cmd := exec.Command(parts[0], parts[1:]...)
		output, err := cmd.Output()
		if err == nil {
			return true, strings.TrimSpace(string(output)), nil
		}
	}
	return false, "", nil
}

// runInstallCommands executes installation commands sequentially
func runInstallCommands(installCmds []string) error {
	for _, cmdStr := range installCmds {
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}
		
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// getVersion gets version from check commands
func getVersion(checkCmds []string) (string, error) {
	installed, version, err := checkInstalled(checkCmds)
	if !installed || err != nil {
		return "", err
	}
	return version, nil
}

// ExecuteInstall executes the installation process for a given InstallCommand
func ExecuteInstall(cmd InstallCommand) InstallResult {
	logger := NewLogger()
	
	// Check if already installed
	installed, version, err := checkInstalled(cmd.CheckCommands)
	if err != nil {
		logger.Error(fmt.Sprintf("%s のインストール状態確認でエラー: %v", cmd.Name, err))
		return InstallResult{
			AlreadyInstalled: false,
			Success:          false,
			Version:          "",
			Err:              err,
		}
	}
	
	if installed {
		logger.Log(fmt.Sprintf("%s は既にインストール済み: %s", cmd.Name, version))
		return InstallResult{
			AlreadyInstalled: true,
			Success:          true,
			Version:          version,
			Err:              nil,
		}
	}
	
	// Install if not installed
	logger.Log(fmt.Sprintf("%s をインストール中...", cmd.Name))
	
	// Execute install commands
	if err := runInstallCommands(cmd.InstallCommands); err != nil {
		logger.Error(fmt.Sprintf("%s のインストールコマンド実行でエラー: %v", cmd.Name, err))
		return InstallResult{
			AlreadyInstalled: false,
			Success:          false,
			Version:          "",
			Err:              err,
		}
	}
	
	// Execute custom install function
	if cmd.InstallFunc != nil {
		if err := cmd.InstallFunc(); err != nil {
			logger.Error(fmt.Sprintf("%s のカスタムインストール処理でエラー: %v", cmd.Name, err))
			return InstallResult{
				AlreadyInstalled: false,
				Success:          false,
				Version:          "",
				Err:              err,
			}
		}
	}
	
	// Get version after installation
	finalVersion, err := getVersion(cmd.CheckCommands)
	if err != nil {
		logger.Error(fmt.Sprintf("%s のバージョン取得でエラー: %v", cmd.Name, err))
		return InstallResult{
			AlreadyInstalled: false,
			Success:          false,
			Version:          "",
			Err:              err,
		}
	}
	
	logger.Success(fmt.Sprintf("%s をインストールしました: %s", cmd.Name, finalVersion))
	return InstallResult{
		AlreadyInstalled: false,
		Success:          true,
		Version:          finalVersion,
		Err:              nil,
	}
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
	// Example usage of the new InstallCommand system
	nodeCmd := InstallCommand{
		CheckCommands:   []string{"node --version", "nodejs --version"},
		InstallCommands: []string{
			"curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -",
			"sudo apt-get install -y nodejs",
		},
		InstallFunc:    defaultInstallFunc,
		Name:          "Node.js",
		Description:   "JavaScript runtime",
	}
	
	claudeCmd := InstallCommand{
		CheckCommands:   []string{"claude-code --version"},
		InstallCommands: []string{"npm install -g @anthropic-ai/claude-code"},
		InstallFunc:    defaultInstallFunc,
		Name:          "Claude Code",
		Description:   "AI-powered coding assistant",
	}
	
	// Execute installations using the new system
	fmt.Println("=== Testing New Installation System ===")
	nodeResult := ExecuteInstall(nodeCmd)
	claudeResult := ExecuteInstall(claudeCmd)
	
	fmt.Printf("Node.js Result: Success=%v, AlreadyInstalled=%v, Version=%s\n", 
		nodeResult.Success, nodeResult.AlreadyInstalled, nodeResult.Version)
	fmt.Printf("Claude Code Result: Success=%v, AlreadyInstalled=%v, Version=%s\n", 
		claudeResult.Success, claudeResult.AlreadyInstalled, claudeResult.Version)
	
	fmt.Println("\n=== Using Legacy Setup ===")
	setup := NewSetup()
	if err := setup.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}