package cmd

import (
	"setupx/internal/config"
	"setupx/internal/models"
	"setupx/internal/pkgmgr"
	"setupx/internal/runner"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [pkg]",
	Short: "Install a specific package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]
		
		var cfg *models.Config
		
		// Let's assume setupx.yaml is optional here or we just use it for mapping
		cfg, err := config.LoadConfig("setupx.yaml")
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatalf("Error loading config: %v", err)
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
			log.Fatalf("Error: %v", err)
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
			return
		}

		installCmd := mgr.InstallCommand([]string{targetPkg}, versions)
		if err := run.Run(installCmd); err != nil {
			log.Fatalf("Error executing install: %v", err)
		}

		if dryRun {
			fmt.Printf("[Dry-run] Would install %s\n", targetPkg)
		} else {
			fmt.Printf("[Success] %s installed\n", targetPkg)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
