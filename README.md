# setupx 📦

A simple, fast, and dependency-light cross-platform package manager for developers. Define your packages once in a YAML file and install them across macOS, Linux, and Windows using native package managers.

## Features

- **Cross-Platform**: Maps generic package names to native package managers (`brew`, `apt`, `dnf`, `winget`, `scoop`).
- **OS Detection**: Automatically detects your operating system and selects the right tool.
- **Remote Gist Support**: Fetch and apply configurations directly from a URL (e.g., GitHub Gists).
- **Dry-Run Mode**: Preview commands without executing them using the `--dry-run` flag.
- **Explain Mode**: Understand exactly how a package name is mapped and what command will run.
- **Simple Configuration**: Uses a clean `setupx.yaml` for package lists and custom mappings.

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed, then:

```bash
git clone https://github.com/youruser/setupx.git
cd setupx
go build -o setupx main.go
sudo mv setupx /usr/local/bin/ # Optional: move to path
```

## Configuration (`setupx.yaml`)

Create a `setupx.yaml` file in your project root to define your environment:

```yaml
packages:
  - neovim
  - git
  - fzf

mappings:
  neovim:
    windows: Neovim.Neovim
    linux: neovim
    mac: neovim
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

### 🔍 Explain a Package
See how `setupx` maps a package name to your current OS:
```bash
setupx explain neovim
```

### 📦 Install a Specific Package
Install a package directly (it will use mappings from `setupx.yaml` if available):
```bash
setupx install ripgrep
```

### 🔍 Search for a Package
Find the correct package ID from your native package manager:
```bash
setupx search neovim
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

| OS | Default Manager |
|---|---|
| **macOS** | Homebrew (`brew`) |
| **Linux** | `apt` (Default), `dnf` |
| **Windows** | `winget` (Default), `scoop` |

## Development

Run tests:
```bash
go test ./...
```

Build:
```bash
go build -o setupx main.go
```

## License
MIT
