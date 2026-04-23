package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"setupx/internal/config"
	"setupx/internal/models"
	"setupx/internal/pkgmgr"
	"setupx/internal/runner"

	"github.com/spf13/cobra"
)

var configURL string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Install all packages from setupx.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *models.Config
		var err error

		if configURL != "" {
			cfg, err = config.LoadConfigFromURL(configURL)
		} else {
			cfg, err = config.LoadConfig("setupx.yaml")
		}

		if err != nil {
			return fmt.Errorf("Error loading config: %w", err)
		}

		osName := pkgmgr.DetectOS()
		mgr, err := pkgmgr.GetManager(osName, cfg.PackageManager)
		if err != nil {
			return fmt.Errorf("Error: %w", err)
		}

		if len(cfg.Packages) == 0 {
			fmt.Println("No packages to install.")
			return nil
		}

		run := &runner.Runner{DryRun: dryRun}
		for _, p := range cfg.Packages {
			targetPkg := cfg.GetPackageName(p, string(osName))

			// 1. Check if already installed (Idempotency)
			if !dryRun && run.Check(mgr.IsInstalledCommand(targetPkg)) {
				fmt.Printf("[Skipped] %s is already installed\n", p)
				continue
			}

			// 2. Build version map for this specific package
			versions := make(map[string]string)
			if detail, ok := cfg.Mappings[p]; ok && detail.Version != "" {
				versions[targetPkg] = detail.Version
			}

			installCmd := mgr.InstallCommand([]string{targetPkg}, versions)

			if err := run.Run(installCmd); err != nil {
				fmt.Printf("[Warning] Failed to install %s: %v\n", p, err)
			} else if dryRun {
				fmt.Printf("[Dry-run] Would install %s\n", p)
				if detail, ok := cfg.Mappings[p]; ok && len(detail.PostInstall) > 0 {
					for _, hook := range detail.PostInstall {
						fmt.Printf("[Dry-run] Would run hook for %s: %s\n", p, hook)
					}
				}
			} else {
				fmt.Printf("[Success] %s installed\n", p)
				if detail, ok := cfg.Mappings[p]; ok && len(detail.PostInstall) > 0 {
					for _, hook := range detail.PostInstall {
						fmt.Printf("[Hook] Running: %s\n", hook)
						var hookCmd *exec.Cmd
						if osName == pkgmgr.Windows {
							hookCmd = exec.Command("cmd", "/c", hook)
						} else {
							hookCmd = exec.Command("sh", "-c", hook)
						}
						hookCmd.Stdout = os.Stdout
						hookCmd.Stderr = os.Stderr
						if err := hookCmd.Run(); err != nil {
							fmt.Printf("[Warning] Hook failed for %s: %v\n", p, err)
						}
					}
				}
			}
		}
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVarP(&configURL, "url", "u", "", "URL to a remote setupx.yaml (e.g. GitHub Gist)")
	rootCmd.AddCommand(applyCmd)
}
