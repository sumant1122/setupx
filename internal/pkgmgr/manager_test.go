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
		versions map[string]string
		expected []string
	}{
		{
			name:     "AptWithVersion",
			mgr:      AptManager{},
			pkgs:     []string{"neovim"},
			versions: map[string]string{"neovim": "0.9.0"},
			expected: []string{"sudo", "apt", "install", "-y", "neovim=0.9.0"},
		},
		{
			name:     "BrewWithVersion",
			mgr:      BrewManager{},
			pkgs:     []string{"go"},
			versions: map[string]string{"go": "1.21"},
			expected: []string{"brew", "install", "go@1.21"},
		},
		{
			name:     "WingetWithVersion",
			mgr:      WingetManager{},
			pkgs:     []string{"Neovim.Neovim"},
			versions: map[string]string{"Neovim.Neovim": "0.9.1"},
			expected: []string{"winget", "install", "Neovim.Neovim", "--version", "0.9.1"},
		},
		{
			name:     "AptNoVersion",
			mgr:      AptManager{},
			pkgs:     []string{"git"},
			versions: nil,
			expected: []string{"sudo", "apt", "install", "-y", "git"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.mgr.InstallCommand(tt.pkgs, tt.versions)
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
