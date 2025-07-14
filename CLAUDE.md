# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Purpose

This repository provides automated Ubuntu development environment setup through two approaches:
- **Bash script** (`setup.sh`) - Traditional shell-based setup
- **Ansible playbook** (`playbook.yml`) - Declarative configuration management (recommended)

Both install the same development tools: Node.js LTS, Claude Code, Neovim, Fish shell, Yazi file manager, Lazygit, and GitHub CLI, plus personal dotfiles from external repositories.

## Commands

### Testing and Validation
```bash
# Test Ansible playbook (dry run)
ansible-playbook playbook.yml --check

# Run with verbose output
ansible-playbook playbook.yml -v

# Test bash script locally (requires Ubuntu)
bash setup.sh
```

### Development Commands
```bash
# Lint YAML files
yamllint playbook.yml

# Validate Ansible syntax
ansible-playbook playbook.yml --syntax-check
```

## Architecture

### Core Design Principles
- **Idempotency**: Both implementations check for existing installations before proceeding
- **Error resilience**: Continue setup even if individual components fail
- **Modular structure**: Each tool installation is independently handled

### Implementation Approaches

**Bash Script Architecture** (`setup.sh`):
- Function-based modular design with `run_with_error_handling()` wrapper
- Color-coded logging system (INFO/SUCCESS/ERROR)
- Manual idempotency through command existence checks
- Sequential installation with error tolerance

**Ansible Playbook Architecture** (`playbook.yml`):
- Declarative task definitions with built-in idempotency
- Uses `creates`, `force: no`, and module-level checks
- Leverages Ansible facts (`ansible_architecture`, `ansible_env`) for system detection
- Proper privilege escalation with `become` directives

### Key Design Patterns

**Installation Flow**:
1. System package installation via apt
2. External repository setup (Node.js, GitHub CLI)
3. Binary installations from GitHub releases (Neovim, Yazi)
4. Configuration cloning from external git repositories
5. Shell configuration (Fish as default shell)

**External Dependencies**:
- Configuration repositories: `ishida722/nvim`, `ishida722/fish`, `ishida722/krapp-config`
- GitHub releases for latest binaries
- NodeSource repository for Node.js LTS

### Go Migration Specification
The `go-setup-spec.md` outlines a planned Go rewrite using `InstallCommand` structs with:
- `CheckCommands` for installation detection
- `InstallCommands` for setup execution  
- `InstallFunc` for custom installation logic
- Structured error handling and logging

## Important Implementation Notes

- Neovim installation bypasses apt due to outdated package versions
- Yazi requires architecture detection for correct binary selection
- Fish shell default setting requires user shell change (effective after re-login)
- All git clones use `force: no` to preserve existing configurations
- Ansible version assumes localhost execution with local connection

## Troubleshooting

### Ansible Playbook Failures

#### "Failed to update apt cache: unknown reason"

**Symptoms:**
```
TASK [Install basic dependencies and Fish shell] ***********************************************************************
fatal: [localhost]: FAILED! => {"changed": false, "msg": "Failed to update apt cache: unknown reason"}
```

**Common Causes:**

1. **Third-party repository GPG key issues** (most common)
   - Steam repository GPG key missing: `NO_PUBKEY F24AEA9FB05498B7`
   - Other software repositories with invalid signatures

2. **Network connectivity issues**
   - Repository servers unreachable
   - Proxy configuration problems

3. **Disk space issues**
   - `/var/cache/apt/` full
   - Root filesystem full

**Diagnosis Steps:**

1. Test apt update directly:
   ```bash
   sudo apt update
   ```

2. Check for GPG key errors in output:
   ```bash
   sudo apt update 2>&1 | grep -i "NO_PUBKEY\|GPG"
   ```

3. Verify network connectivity:
   ```bash
   ping -c 3 archive.ubuntu.com
   ```

4. Check disk space:
   ```bash
   df -h /var/cache/apt/
   ```

**Solutions:**

**For Steam repository GPG key error:**
```bash
# Add correct Steam GPG key
curl -fsSL https://repo.steampowered.com/steam/archive/stable/steam.gpg | sudo gpg --dearmor -o /usr/share/keyrings/steam.gpg

# Fix Steam repository configuration
sudo sh -c 'echo "deb [arch=amd64,i386 signed-by=/usr/share/keyrings/steam.gpg] https://repo.steampowered.com/steam stable steam" > /etc/apt/sources.list.d/steam-stable.list'

# Update apt cache
sudo apt update
```

**For general repository issues:**
```bash
# Remove problematic repository temporarily
sudo mv /etc/apt/sources.list.d/problematic-repo.list /etc/apt/sources.list.d/problematic-repo.list.bak

# Update apt cache
sudo apt update

# Re-add repository with correct configuration
```

**Prevention:**
- Always verify repository GPG keys during manual software installation
- Use official package repositories when possible
- Regular system maintenance to prevent disk space issues