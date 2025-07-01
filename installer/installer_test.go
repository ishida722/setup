package installer

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

// Test helper to capture output
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputChan <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = old

	return <-outputChan
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}
	if logger.info == nil {
		t.Error("Logger.info is nil")
	}
	if logger.success == nil {
		t.Error("Logger.success is nil")
	}
	if logger.errorColor == nil {
		t.Error("Logger.errorColor is nil")
	}
}

func TestLoggerMethods(t *testing.T) {
	logger := NewLogger()

	tests := []struct {
		name     string
		method   func(string)
		message  string
		expected string
	}{
		{
			name:     "Log method",
			method:   logger.Log,
			message:  "test info message",
			expected: "test info message",
		},
		{
			name:     "Success method",
			method:   logger.Success,
			message:  "test success message",
			expected: "test success message",
		},
		{
			name:     "Error method",
			method:   logger.Error,
			message:  "test error message",
			expected: "test error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.method(tt.message)
			})
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestCheckInstalled(t *testing.T) {
	tests := []struct {
		name          string
		checkCommands []string
		expectError   bool
		expectFound   bool
	}{
		{
			name:          "Command exists - echo",
			checkCommands: []string{"echo --version", "echo test"},
			expectError:   false,
			expectFound:   true,
		},
		{
			name:          "Command does not exist",
			checkCommands: []string{"nonexistentcommand --version"},
			expectError:   false,
			expectFound:   false,
		},
		{
			name:          "Multiple commands - first fails, second succeeds",
			checkCommands: []string{"nonexistentcommand --version", "echo test"},
			expectError:   false,
			expectFound:   true,
		},
		{
			name:          "Empty command list",
			checkCommands: []string{},
			expectError:   false,
			expectFound:   false,
		},
		{
			name:          "Command with empty string",
			checkCommands: []string{""},
			expectError:   false,
			expectFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, version, err := checkInstalled(tt.checkCommands)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if found != tt.expectFound {
				t.Errorf("Expected found=%v, got found=%v", tt.expectFound, found)
			}
			if tt.expectFound && version == "" {
				t.Error("Expected version string but got empty")
			}
		})
	}
}

func TestInstallCommandCheckInstalled(t *testing.T) {
	tests := []struct {
		name          string
		cmd           InstallCommand
		expectError   bool
		expectFound   bool
	}{
		{
			name: "Valid command exists",
			cmd: InstallCommand{
				CheckCommands: []string{"echo test"},
				Name:          "Test Command",
			},
			expectError: false,
			expectFound: true,
		},
		{
			name: "Command does not exist",
			cmd: InstallCommand{
				CheckCommands: []string{"nonexistentcommand"},
				Name:          "Nonexistent Command",
			},
			expectError: false,
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, version, err := tt.cmd.CheckInstalled()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if found != tt.expectFound {
				t.Errorf("Expected found=%v, got found=%v", tt.expectFound, found)
			}
			if tt.expectFound && version == "" {
				t.Error("Expected version string but got empty")
			}
		})
	}
}

func TestInstallCommandInstall_AlreadyInstalled(t *testing.T) {
	cmd := InstallCommand{
		CheckCommands:   []string{"echo v1.0.0"},
		InstallCommands: []string{"echo 'should not run'"},
		Name:            "Test Already Installed",
		Description:     "Test software that is already installed",
	}

	result := cmd.Install()

	if !result.Success {
		t.Error("Expected success=true for already installed software")
	}
	if !result.AlreadyInstalled {
		t.Error("Expected AlreadyInstalled=true")
	}
	if result.Err != nil {
		t.Errorf("Expected no error, got: %v", result.Err)
	}
	if result.Version == "" {
		t.Error("Expected version string")
	}
}

func TestInstallCommandInstall_NotInstalled(t *testing.T) {
	cmd := InstallCommand{
		CheckCommands:   []string{"nonexistentcommand --version"},
		InstallCommands: []string{"echo 'installing'"},
		Name:            "Test Not Installed",
		Description:     "Test software that needs installation",
	}

	// This test will succeed because echo command exists and will be used for version check
	result := cmd.Install()

	if !result.Success {
		t.Errorf("Expected success=true, got error: %v", result.Err)
	}
	if result.AlreadyInstalled {
		t.Error("Expected AlreadyInstalled=false")
	}
}

func TestInstallCommandInstall_WithCustomFunction(t *testing.T) {
	customFuncCalled := false
	cmd := InstallCommand{
		CheckCommands:   []string{"nonexistentcommand --version"},
		InstallCommands: []string{"echo 'installing'"},
		InstallFunc: func() error {
			customFuncCalled = true
			return nil
		},
		Name:        "Test Custom Function",
		Description: "Test software with custom install function",
	}

	result := cmd.Install()

	if !customFuncCalled {
		t.Error("Expected custom install function to be called")
	}
	if !result.Success {
		t.Errorf("Expected success=true, got error: %v", result.Err)
	}
}

func TestInstallCommandInstall_CustomFunctionError(t *testing.T) {
	expectedError := errors.New("custom function error")
	cmd := InstallCommand{
		CheckCommands:   []string{"nonexistentcommand --version"},
		InstallCommands: []string{"echo 'installing'"},
		InstallFunc: func() error {
			return expectedError
		},
		Name:        "Test Custom Function Error",
		Description: "Test software with failing custom install function",
	}

	result := cmd.Install()

	if result.Success {
		t.Error("Expected success=false due to custom function error")
	}
	if result.AlreadyInstalled {
		t.Error("Expected AlreadyInstalled=false")
	}
	if result.Err == nil {
		t.Error("Expected error from custom function")
	}
}

func TestRunInstallCommands(t *testing.T) {
	tests := []struct {
		name         string
		commands     []string
		expectError  bool
	}{
		{
			name:        "Valid commands",
			commands:    []string{"echo test1", "echo test2"},
			expectError: false,
		},
		{
			name:        "Empty command list",
			commands:    []string{},
			expectError: false,
		},
		{
			name:        "Command with empty string",
			commands:    []string{""},
			expectError: false,
		},
		{
			name:        "Invalid command",
			commands:    []string{"nonexistentcommand"},
			expectError: true,
		},
		{
			name:        "Mixed valid and invalid commands",
			commands:    []string{"echo test", "nonexistentcommand"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runInstallCommands(tt.commands)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name          string
		checkCommands []string
		expectError   bool
		expectVersion bool
	}{
		{
			name:          "Valid version command",
			checkCommands: []string{"echo v1.0.0"},
			expectError:   false,
			expectVersion: true,
		},
		{
			name:          "Invalid version command",
			checkCommands: []string{"nonexistentcommand --version"},
			expectError:   false,
			expectVersion: false,
		},
		{
			name:          "Empty command list",
			checkCommands: []string{},
			expectError:   false,
			expectVersion: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := getVersion(tt.checkCommands)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tt.expectVersion && version == "" {
				t.Error("Expected version string but got empty")
			}
			if !tt.expectVersion && version != "" {
				t.Errorf("Expected empty version but got: %s", version)
			}
		})
	}
}

func TestExecuteInstall_Legacy(t *testing.T) {
	cmd := InstallCommand{
		CheckCommands:   []string{"echo v1.0.0"},
		InstallCommands: []string{"echo 'installing'"},
		Name:            "Test Legacy Function",
		Description:     "Test legacy ExecuteInstall function",
	}

	result := ExecuteInstall(cmd)

	// Should behave the same as cmd.Install()
	if !result.Success {
		t.Error("Expected success=true for already installed software")
	}
	if !result.AlreadyInstalled {
		t.Error("Expected AlreadyInstalled=true")
	}
	if result.Err != nil {
		t.Errorf("Expected no error, got: %v", result.Err)
	}
}

func TestDefaultInstallFunc(t *testing.T) {
	err := defaultInstallFunc()
	if err != nil {
		t.Errorf("Expected defaultInstallFunc to return nil, got: %v", err)
	}
}

// Integration test with real commands
func TestIntegration_RealCommands(t *testing.T) {
	// Test with a command that should exist on most systems
	cmd := InstallCommand{
		CheckCommands: []string{"echo --help"},
		Name:          "Echo Command",
		Description:   "Built-in echo command",
	}

	found, version, err := cmd.CheckInstalled()
	if err != nil {
		t.Errorf("Unexpected error checking echo command: %v", err)
	}
	if !found {
		t.Error("Expected echo command to be found")
	}
	if version == "" {
		t.Error("Expected version output from echo --help")
	}

	// Test install (should detect as already installed)
	result := cmd.Install()
	if !result.Success {
		t.Errorf("Expected success, got error: %v", result.Err)
	}
	if !result.AlreadyInstalled {
		t.Error("Expected echo to be detected as already installed")
	}
}

// Benchmark tests
func BenchmarkCheckInstalled(b *testing.B) {
	checkCmds := []string{"echo test"}
	for i := 0; i < b.N; i++ {
		checkInstalled(checkCmds)
	}
}

func BenchmarkInstallCommandCheckInstalled(b *testing.B) {
	cmd := InstallCommand{
		CheckCommands: []string{"echo test"},
		Name:          "Benchmark Test",
	}
	for i := 0; i < b.N; i++ {
		cmd.CheckInstalled()
	}
}