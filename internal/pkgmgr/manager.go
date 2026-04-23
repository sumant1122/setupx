package pkgmgr

import (
	"fmt"
	"os/exec"
	"strings"
)

type SearchResult struct {
	Name        string
	Description string
}

type PackageManager interface {
	ParseSearchOutput(out string) []SearchResult
	InstallCommand(packages []string, versions map[string]string) []string
	SearchCommand(pkg string) []string
	ExactSearchCommand(pkg string) []string
	IsInstalledCommand(pkg string) []string
}

type AptManager struct{}

func (a AptManager) ParseSearchOutput(out string) []SearchResult {
	var results []SearchResult
	lines := strings.Split(out, "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "Sorting...") || strings.HasPrefix(line, "Full Text Search...") || strings.Contains(line, "WARNING") {
			continue
		}
		parts := strings.SplitN(line, "/", 2)
		if len(parts) == 2 {
			name := parts[0]
			desc := ""
			if i+1 < len(lines) {
				desc = strings.TrimSpace(lines[i+1])
				i++
			}
			results = append(results, SearchResult{Name: name, Description: desc})
			if len(results) >= 20 {
				break
			}
		}
	}
	return results
}

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
	return []string{"apt", "search", "--names-only", fmt.Sprintf("^%s$", pkg)}
}

func (a AptManager) IsInstalledCommand(pkg string) []string {
	return []string{"dpkg", "-s", pkg}
}

type DnfManager struct{}

func (d DnfManager) ParseSearchOutput(out string) []SearchResult {
	var results []SearchResult
	lines := strings.Split(out, "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "Last metadata expiration check:") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			name := strings.TrimSpace(parts[0])
			nameParts := strings.Split(name, ".")
			results = append(results, SearchResult{Name: nameParts[0], Description: strings.TrimSpace(parts[1])})
		} else {
			parts = strings.SplitN(line, "/", 2)
			if len(parts) == 2 {
				name := parts[0]
				desc := ""
				if i+1 < len(lines) {
					desc = strings.TrimSpace(lines[i+1])
					i++
				}
				results = append(results, SearchResult{Name: name, Description: desc})
			}
		}
		if len(results) >= 20 {
			break
		}
	}
	return results
}

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
	return []string{"dnf", "search", "--names-only", pkg}
}

func (d DnfManager) IsInstalledCommand(pkg string) []string {
	return []string{"dnf", "list", "installed", pkg}
}

type BrewManager struct{}

func (b BrewManager) ParseSearchOutput(out string) []SearchResult {
	var results []SearchResult
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "==>") {
			continue
		}
		results = append(results, SearchResult{Name: line, Description: "N/A"})
		if len(results) >= 20 {
			break
		}
	}
	return results
}

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

func (b BrewManager) IsInstalledCommand(pkg string) []string {
	return []string{"brew", "list", pkg}
}

type WingetManager struct{}

func (w WingetManager) ParseSearchOutput(out string) []SearchResult {
	var results []SearchResult
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Name") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "No package found") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			id := fields[1]
			for _, f := range fields {
				if strings.Count(f, ".") >= 1 && !strings.Contains(f, ":") {
					id = f
					break
				}
			}
			results = append(results, SearchResult{Name: id, Description: "Available on Windows"})
			if len(results) >= 20 {
				break
			}
		}
	}
	return results
}

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
	return []string{"winget", "search", pkg}
}

func (w WingetManager) ExactSearchCommand(pkg string) []string {
	return []string{"winget", "search", "--id", pkg, "--exact"}
}

func (w WingetManager) IsInstalledCommand(pkg string) []string {
	return []string{"winget", "list", "--id", pkg, "--exact"}
}

type ScoopManager struct{}

func (s ScoopManager) ParseSearchOutput(out string) []SearchResult {
	var results []SearchResult
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Results from") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "Name") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			results = append(results, SearchResult{Name: fields[0], Description: "Available in Scoop"})
			if len(results) >= 20 {
				break
			}
		}
	}
	return results
}

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

func (s ScoopManager) IsInstalledCommand(pkg string) []string {
	return []string{"scoop", "list", pkg}
}

type PacmanManager struct{}

func (p PacmanManager) ParseSearchOutput(out string) []SearchResult {
	var results []SearchResult
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			results = append(results, SearchResult{Name: line, Description: "N/A"})
			if len(results) >= 20 {
				break
			}
		}
	}
	return results
}

func (p PacmanManager) InstallCommand(packages []string, versions map[string]string) []string {
	var targets []string
	for _, p := range packages {
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

func (p PacmanManager) IsInstalledCommand(pkg string) []string {
	return []string{"pacman", "-Qi", pkg}
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
