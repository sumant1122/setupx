package models

import "testing"

func TestGetPackageName(t *testing.T) {
	cfg := &Config{
		Mappings: map[string]PackageDetail{
			"neovim": {
				Windows: "Neovim.Neovim",
				Linux:   "neovim",
				Mac:     "neovim",
			},
		},
	}

	tests := []struct {
		pkg      string
		os       string
		expected string
	}{
		{"neovim", "windows", "Neovim.Neovim"},
		{"neovim", "linux", "neovim"},
		{"neovim", "mac", "neovim"},
		{"git", "windows", "git"}, // Should return the name itself if no mapping
	}

	for _, tt := range tests {
		result := cfg.GetPackageName(tt.pkg, tt.os)
		if result != tt.expected {
			t.Errorf("GetPackageName(%s, %s) = %s; expected %s", tt.pkg, tt.os, result, tt.expected)
		}
	}
}
