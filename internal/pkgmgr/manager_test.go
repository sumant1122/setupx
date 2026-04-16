package pkgmgr

import (
	"reflect"
	"testing"
)

func TestInstallCommands(t *testing.T) {
	tests := []struct {
		name     string
		mgr      PackageManager
		pkgs     []string
		expected []string
	}{
		{
			name:     "Apt",
			mgr:      AptManager{},
			pkgs:     []string{"neovim", "git"},
			expected: []string{"sudo", "apt", "install", "-y", "neovim", "git"},
		},
		{
			name:     "Brew",
			mgr:      BrewManager{},
			pkgs:     []string{"neovim", "git"},
			expected: []string{"brew", "install", "neovim", "git"},
		},
		{
			name:     "Winget",
			mgr:      WingetManager{},
			pkgs:     []string{"Neovim.Neovim", "Git.Git"},
			expected: []string{"winget", "install", "Neovim.Neovim", "Git.Git"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mgr.InstallCommand(tt.pkgs)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("%s.InstallCommand() = %v; expected %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetManager(t *testing.T) {
	tests := []struct {
		name     string
		os       OSName
		override string
		expected reflect.Type
	}{
		{"MacBrew", Mac, "", reflect.TypeOf(BrewManager{})},
		{"LinuxOverrideDnf", Linux, "dnf", reflect.TypeOf(DnfManager{})},
		{"WindowsOverrideScoop", Windows, "scoop", reflect.TypeOf(ScoopManager{})},
		{"LinuxOverridePacman", Linux, "pacman", reflect.TypeOf(PacmanManager{})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr, err := GetManager(tt.os, tt.override)
			if err != nil {
				t.Fatalf("GetManager() error = %v", err)
			}
			if reflect.TypeOf(mgr) != tt.expected {
				t.Errorf("GetManager() = %v; expected %v", reflect.TypeOf(mgr), tt.expected)
			}
		})
	}
}
