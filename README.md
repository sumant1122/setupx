# setupx 📦

[![Go Version](https://img.shields.io/github/go-mod/go-version/sumant1122/setupx)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/sumant1122/setupx/actions/workflows/go.yml/badge.svg)](https://github.com/sumant1122/setupx/actions)

Set up your development environment on macOS, Linux, and Windows using a single YAML file.

## Why setupx?

Developers often switch between macOS (work), Linux (personal/server), and Windows. Maintaining three separate setup scripts or remembering different package manager commands is a major context switch. 

**setupx** allows you to define your tools once. It handles the OS detection, name mapping, and version pinning automatically, so you can bootstrap any machine in seconds.

## setupx Demo
<img width="800" height="450" alt="setupx-gif" src="https://github.com/user-attachments/assets/3084f2c0-01ec-41f7-a0ec-0a965f3eb286" />

## Features

- **Cross-Platform**: Maps generic package names to native package managers (`brew`, `apt`, `dnf`, `pacman`, `winget`, `scoop`).
- **OS Detection**: Automatically detects your operating system and selects the right tool.
- **Idempotency**: Automatically skips packages that are already installed, making updates nearly instant.
- **Remote Gist Support**: Fetch and apply configurations directly from a URL (e.g., GitHub Gists).
- **Native Search**: Search for package IDs directly through `setupx` with clean, table-formatted results.
- **Version Pinning**: Specify exact package versions to ensure reproducible environments.
- **Dry-Run Mode**: Preview commands without executing them using the `--dry-run` flag.
- **Explain Mode**: Understand exactly how a package name is mapped and what command will run.
- **Simple Configuration**: Uses a clean `setupx.yaml` for package lists and custom mappings.

## Installation

### Option 1: Binary Download (Recommended)
Download the latest pre-built binary for your platform from the [Releases](https://github.com/sumant1122/setupx/releases) page.

1.  **Linux/Mac**: 
    ```bash
    chmod +x setupx-<os>-amd64
    sudo mv setupx-<os>-amd64 /usr/local/bin/setupx
    ```
2.  **Windows**: 
    - Rename `setupx-windows-amd64.exe` to `setupx.exe`.
    - Add the folder containing `setupx.exe` to your system **PATH**.

### Option 2: Build from Source
Ensure you have [Go](https://go.dev/doc/install) installed, then:

```bash
git clone https://github.com/sumant1122/setupx.git
cd setupx
go build -o setupx main.go
sudo mv setupx /usr/local/bin/ # Optional: move to path
```

## Configuration (`setupx.yaml`)

Create a `setupx.yaml` file in your project root to define your environment:

```yaml
package_manager: brew  # Optional: Force a specific package manager

packages:
  - neovim
  - git
  - fzf

mappings:
  neovim:
    windows: Neovim.Neovim
    linux: neovim
    mac: neovim
    version: "0.9.5"  # Optional: Pin to a specific version
  git:
    windows: Git.Git
```

## Usage

### 🚀 Apply Configuration
Install all packages defined in `setupx.yaml`:
```bash
setupx apply
```

### 🌍 Remote Configuration (Gist)
Bootstrap a new machine using a configuration stored online (e.g., GitHub Gist raw URL):
```bash
setupx apply --url https://gist.githubusercontent.com/user/id/raw/setupx.yaml
```

### 🔍 Search for a Package
Find the correct package ID from your native package manager (results are formatted in a clean table):
```bash
setupx search neovim
```

### 🔍 Explain a Package
See how `setupx` maps a package name and what command it would run:
```bash
setupx explain neovim
```

### 📦 Install a Specific Package
Install a package directly (it will use mappings from `setupx.yaml` if available):
```bash
setupx install ripgrep
```

### 🛡️ Dry Run
Preview any command without making changes:
```bash
setupx apply --dry-run
setupx install fzf -d
```

### ℹ️ Version
Check the current version:
```bash
setupx --version
```

## Supported Package Managers

| OS | Default Manager | Supported Alternatives |
|---|---|---|
| **macOS** | Homebrew (`brew`) | |
| **Linux** | `apt` | `dnf`, `pacman`, `brew` |
| **Windows** | `winget` | `scoop` |

## Development

Run tests:
```bash
go test ./...
```

Build:
```bash
go build -o setupx main.go
```

## Recent Updates

- **Architecture**: Decoupled package manager search parsing for better modularity and testability.

## Roadmap

We're constantly working to make **setupx** the best environment orchestrator. Here’s what’s coming next:

- [x] **Post-Install Hooks**: Run custom shell commands after a package is successfully installed.
- [ ] **Interactive Onboarding**: Automatically search and suggest mappings when a package is missing from your config.
- [ ] **Configuration Variables**: Use placeholders like `{{HOME}}` or `{{OS}}` in your YAML.
- [ ] **Self-Update**: A simple `setupx self-update` command to stay current.
- [ ] **Export Mode**: Generate a `setupx.yaml` by scanning your current system's installed packages.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines and our [Code of Conduct](CODE_OF_CONDUCT.md).

## License
MIT
