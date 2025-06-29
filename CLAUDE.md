# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## ⚠️ CRITICAL SECURITY WARNING

**NEVER execute the shell scripts in this repository.** These scripts perform system-level installations and modifications that could compromise system security. Only analyze, read, or discuss the code - never run it.

## Repository Overview

This is a setup script repository for configuring new Ubuntu machines with Claude Code and development tools. The repository contains:

- `setup.sh`: Main setup script that installs and configures development environment
- `README.md`: Japanese documentation with usage instructions

## Setup Script Architecture

The setup is now organized with modular scripts in the `scripts/` directory:

- **setup.sh**: Main orchestrator script that sources and calls all installation modules
- **scripts/utils.sh**: Shared utility functions (logging, colors)
- **scripts/install-nodejs.sh**: Node.js LTS installation via NodeSource repository
- **scripts/install-claude-code.sh**: Claude Code global npm installation
- **scripts/install-neovim.sh**: Latest Neovim download and installation from GitHub releases
- **scripts/install-fish.sh**: Fish shell installation and default shell configuration
- **scripts/clone-configs.sh**: Configuration files cloning from external repositories

Each script can be run independently or as part of the main setup process.

## Key Commands

### Running the Setup Script
```bash
# Direct execution from GitHub
curl -fsSL https://raw.githubusercontent.com/ishida722/setup/main/setup.sh | bash

# Local execution (full setup)
chmod +x setup.sh
./setup.sh

# Individual script execution
chmod +x scripts/*.sh
./scripts/install-nodejs.sh
./scripts/install-claude-code.sh
./scripts/install-neovim.sh
./scripts/install-fish.sh
./scripts/clone-configs.sh
```

### Post-Setup Configuration
After running the setup script, users need to:
```bash
# Set API key for Claude Code
export ANTHROPIC_API_KEY='your-key'

# Start using Claude Code
claude-code --help
```

## External Dependencies

The setup script clones configuration files from external repositories:
- Neovim config: `https://github.com/ishida722/nvim`
- Fish config: `https://github.com/ishida722/fish`

## Important Notes

- Script is designed specifically for Ubuntu systems
- Uses `set -e` for strict error handling
- Includes color-coded logging functions for user feedback
- Checks for existing installations before attempting to install
- Default shell change requires re-login to take effect
- All installations are system-wide where applicable (using sudo)