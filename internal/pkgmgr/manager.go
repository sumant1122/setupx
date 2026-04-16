package pkgmgr

import (
	"fmt"
	"os/exec"
	"strings"
)

type PackageManager interface {
	InstallCommand(packages []string, versions map[string]string) []string
	SearchCommand(pkg string) []string
	ExactSearchCommand(pkg string) []string
}

type AptManager struct{}

func (a AptManager) InstallCommand(packages []string, versions map[string]string) []string {
	var targets []string
	for _, p := range packages {
		if v, ok := versions[p]; ok && v != "" {
			targets = append(targets, fmt.Sprintf("%s=%s", p, v))
		} else {
			targets = append(targets, p)
		}
	}
	return append([]string{"sudo", "apt", "install", "-y"}, targets...)
}

func (a AptManager) SearchCommand(pkg string) []string {
	return []string{"apt", "search", "--names-only", pkg}
}

func (a AptManager) ExactSearchCommand(pkg string) []string {
	// Using regex for exact match in apt
	return []string{"apt", "search", "--names-only", fmt.Sprintf("^%s$", pkg)}
}

type DnfManager struct{}

func (d DnfManager) InstallCommand(packages []string, versions map[string]string) []string {
	var targets []string
	for _, p := range packages {
		if v, ok := versions[p]; ok && v != "" {
			targets = append(targets, fmt.Sprintf("%s-%s", p, v))
		} else {
			targets = append(targets, p)
		}
	}
	return append([]string{"sudo", "dnf", "install", "-y"}, targets...)
}

func (d DnfManager) SearchCommand(pkg string) []string {
	return []string{"dnf", "search", "--names-only", pkg}
}

func (d DnfManager) ExactSearchCommand(pkg string) []string {
	return []string{"dnf", "search", "--names-only", pkg} // Dnf handles it fairly well or we'd use exact flags if available
}

type BrewManager struct{}

func (b BrewManager) InstallCommand(packages []string, versions map[string]string) []string {
	var targets []string
	for _, p := range packages {
		if v, ok := versions[p]; ok && v != "" {
			targets = append(targets, fmt.Sprintf("%s@%s", p, v))
		} else {
			targets = append(targets, p)
		}
	}
	return append([]string{"brew", "install"}, targets...)
}

func (b BrewManager) SearchCommand(pkg string) []string {
	return []string{"brew", "search", pkg}
}

func (b BrewManager) ExactSearchCommand(pkg string) []string {
	return []string{"brew", "search", fmt.Sprintf("/^%s$/", pkg)}
}

type WingetManager struct{}

func (w WingetManager) InstallCommand(packages []string, versions map[string]string) []string {
	var targets []string
	for _, p := range packages {
		if v, ok := versions[p]; ok && v != "" {
			targets = append(targets, p, "--version", v)
		} else {
			targets = append(targets, p)
		}
	}
	return append([]string{"winget", "install"}, targets...)
}

func (w WingetManager) SearchCommand(pkg string) []string {
	return []string{"winget", "search", "--name", pkg}
}

func (w WingetManager) ExactSearchCommand(pkg string) []string {
	return []string{"winget", "search", "--exact", pkg}
}

type ScoopManager struct{}

func (s ScoopManager) InstallCommand(packages []string, versions map[string]string) []string {
	var targets []string
	for _, p := range packages {
		if v, ok := versions[p]; ok && v != "" {
			targets = append(targets, fmt.Sprintf("%s@%s", p, v))
		} else {
			targets = append(targets, p)
		}
	}
	return append([]string{"scoop", "install"}, targets...)
}

func (s ScoopManager) SearchCommand(pkg string) []string {
	return []string{"scoop", "search", pkg}
}

func (s ScoopManager) ExactSearchCommand(pkg string) []string {
	return []string{"scoop", "search", pkg}
}

type PacmanManager struct{}

func (p PacmanManager) InstallCommand(packages []string, versions map[string]string) []string {
	var targets []string
	for _, p := range packages {
		// Pacman doesn't support easy version pinning in the install command directly 
		// without downgrading or specific archive handling. 
		// For MVP, we'll just log/ignore or use the standard.
		targets = append(targets, p)
	}
	return append([]string{"sudo", "pacman", "-S", "--noconfirm"}, targets...)
}

func (p PacmanManager) SearchCommand(pkg string) []string {
	return []string{"pacman", "-Ssq", pkg}
}

func (p PacmanManager) ExactSearchCommand(pkg string) []string {
	return []string{"pacman", "-Ssq", fmt.Sprintf("^%s$", pkg)}
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
