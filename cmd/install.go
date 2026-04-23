package cmd

import (
	"fmt"
	"os"
	"setupx/internal/config"
	"setupx/internal/models"
	"setupx/internal/pkgmgr"
	"setupx/internal/runner"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [pkg]",
	Short: "Install a specific package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := args[0]

		var cfg *models.Config

		// Let's assume setupx.yaml is optional here or we just use it for mapping
		cfg, err := config.LoadConfig("setupx.yaml")
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("Error loading config: %w", err)
			}
			// If it doesn't exist, we'll just use the name as is
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

		run := &runner.Runner{DryRun: dryRun}

		// Check if already installed (Idempotency)
		if !dryRun && run.Check(mgr.IsInstalledCommand(targetPkg)) {
			fmt.Printf("[Skipped] %s is already installed\n", pkg)
			return nil
		}

		installCmd := mgr.InstallCommand([]string{targetPkg}, versions)
		if err := run.Run(installCmd); err != nil {
			return fmt.Errorf("Error executing install: %w", err)
		}

		if dryRun {
			fmt.Printf("[Dry-run] Would install %s\n", targetPkg)
		} else {
			fmt.Printf("[Success] %s installed\n", targetPkg)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
