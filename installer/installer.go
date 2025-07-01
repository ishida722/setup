package installer

import (
	"fmt"
	"os"
	"os/exec"
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
	info       *color.Color
	success    *color.Color
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

// CheckInstalled checks if the software is already installed
func (cmd *InstallCommand) CheckInstalled() (bool, string, error) {
	return checkInstalled(cmd.CheckCommands)
}

// Install executes the installation process, including check for existing installation
func (cmd *InstallCommand) Install() InstallResult {
	logger := NewLogger()
	
	// Check if already installed
	installed, version, err := cmd.CheckInstalled()
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
	installFunc := cmd.InstallFunc
	if installFunc == nil {
		installFunc = defaultInstallFunc
	}
	
	if err := installFunc(); err != nil {
		logger.Error(fmt.Sprintf("%s のカスタムインストール処理でエラー: %v", cmd.Name, err))
		return InstallResult{
			AlreadyInstalled: false,
			Success:          false,
			Version:          "",
			Err:              err,
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

// ExecuteInstall executes the installation process for a given InstallCommand (legacy function)
func ExecuteInstall(cmd InstallCommand) InstallResult {
	return cmd.Install()
}