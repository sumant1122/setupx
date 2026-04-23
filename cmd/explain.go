package cmd

import (
	"fmt"
	"os"
	"setupx/internal/config"
	"setupx/internal/models"
	"setupx/internal/pkgmgr"

	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain [pkg]",
	Short: "Explain what command would be run for a package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := args[0]

		var cfg *models.Config
		cfg, err := config.LoadConfig("setupx.yaml")
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("Error loading config: %w", err)
		}

		osName := pkgmgr.DetectOS()
		var pmOverride string
		if cfg != nil {
			pmOverride = cfg.PackageManager
		}
		mgr, err := pkgmgr.GetManager(osName, pmOverride)
		if err != nil {
			return fmt.Errorf("Error: %w", err)
		}

		var targetPkg string
		versions := make(map[string]string)
		if cfg != nil {
			targetPkg = cfg.GetPackageName(pkg, string(osName))
			if detail, ok := cfg.Mappings[pkg]; ok && detail.Version != "" {
				versions[targetPkg] = detail.Version
			}
		} else {
			targetPkg = pkg
		}

		installCmd := mgr.InstallCommand([]string{targetPkg}, versions)

		fmt.Printf("Package: %s\n", pkg)
		fmt.Printf("Detected OS: %s\n", osName)
		fmt.Printf("Mapped Name: %s\n", targetPkg)
		fmt.Printf("Command: %s\n", pkgmgr.FormatCommand(installCmd))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
