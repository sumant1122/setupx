package models

type Config struct {
	Packages       []string                 `yaml:"packages"`
	PackageManager string                   `yaml:"package_manager,omitempty"`
	Mappings       map[string]PackageDetail `yaml:"mappings"`
}

type PackageDetail struct {
	Linux   string `yaml:"linux,omitempty"`
	Mac     string `yaml:"mac,omitempty"`
	Windows string `yaml:"windows,omitempty"`
	Version string `yaml:"version,omitempty"`
	PostInstall []string `yaml:"post_install,omitempty"`
}

func (c *Config) GetPackageName(pkg string, osName string) string {
	detail, ok := c.Mappings[pkg]
	if !ok {
		return pkg // Default to the name itself
	}

	switch osName {
	case "linux":
		if detail.Linux != "" {
			return detail.Linux
		}
	case "mac":
		if detail.Mac != "" {
			return detail.Mac
		}
	case "windows":
		if detail.Windows != "" {
			return detail.Windows
		}
	}
	return pkg
}
