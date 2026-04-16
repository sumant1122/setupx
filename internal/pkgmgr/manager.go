package pkgmgr

import (
	"fmt"
	"os/exec"
	"strings"
)

type PackageManager interface {
	InstallCommand(packages []string) []string
	SearchCommand(pkg string) []string
}

type AptManager struct{}

func (a AptManager) InstallCommand(packages []string) []string {
	return append([]string{"sudo", "apt", "install", "-y"}, packages...)
}

func (a AptManager) SearchCommand(pkg string) []string {
	return []string{"apt", "search", "--names-only", pkg}
}

type DnfManager struct{}

func (d DnfManager) InstallCommand(packages []string) []string {
	return append([]string{"sudo", "dnf", "install", "-y"}, packages...)
}

func (d DnfManager) SearchCommand(pkg string) []string {
	return []string{"dnf", "search", "--names-only", pkg}
}

type BrewManager struct{}

func (b BrewManager) InstallCommand(packages []string) []string {
	return append([]string{"brew", "install"}, packages...)
}

func (b BrewManager) SearchCommand(pkg string) []string {
	return []string{"brew", "search", pkg}
}

type WingetManager struct{}

func (w WingetManager) InstallCommand(packages []string) []string {
	return append([]string{"winget", "install"}, packages...)
}

func (w WingetManager) SearchCommand(pkg string) []string {
	return []string{"winget", "search", "--name", pkg}
}

type ScoopManager struct{}

func (s ScoopManager) InstallCommand(packages []string) []string {
	return append([]string{"scoop", "install"}, packages...)
}

func (s ScoopManager) SearchCommand(pkg string) []string {
	return []string{"scoop", "search", pkg}
}

type PacmanManager struct{}

func (p PacmanManager) InstallCommand(packages []string) []string {
	return append([]string{"sudo", "pacman", "-S", "--noconfirm"}, packages...)
}

func (p PacmanManager) SearchCommand(pkg string) []string {
	return []string{"pacman", "-Ssq", pkg}
}

func GetManager(os OSName, configOverride string) (PackageManager, error) {
	if configOverride != "" {
		return managerFromName(configOverride)
	}

	switch os {
	case Linux:
		for _, name := range []string{"apt", "dnf", "pacman", "brew"} {
			if _, err := exec.LookPath(name); err == nil {
				return managerFromName(name)
			}
		}
		return AptManager{}, nil // Default fallback
	case Mac:
		return BrewManager{}, nil
	case Windows:
		for _, name := range []string{"winget", "scoop"} {
			if _, err := exec.LookPath(name); err == nil {
				return managerFromName(name)
			}
		}
		return WingetManager{}, nil // Default fallback
	default:
		return nil, fmt.Errorf("unsupported OS: %s", os)
	}
}

func managerFromName(name string) (PackageManager, error) {
	switch strings.ToLower(name) {
	case "apt":
		return AptManager{}, nil
	case "dnf":
		return DnfManager{}, nil
	case "pacman":
		return PacmanManager{}, nil
	case "brew":
		return BrewManager{}, nil
	case "winget":
		return WingetManager{}, nil
	case "scoop":
		return ScoopManager{}, nil
	default:
		return nil, fmt.Errorf("unknown package manager: %s", name)
	}
}

func FormatCommand(cmd []string) string {
	return strings.Join(cmd, " ")
}
