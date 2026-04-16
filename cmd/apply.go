package cmd

import (
	"setupx/internal/config"
	"setupx/internal/models"
	"setupx/internal/pkgmgr"
	"setupx/internal/runner"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var configURL string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Install all packages from setupx.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg *models.Config
		var err error

		if configURL != "" {
			cfg, err = config.LoadConfigFromURL(configURL)
		} else {
			cfg, err = config.LoadConfig("setupx.yaml")
		}

		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		osName := pkgmgr.DetectOS()
		mgr, err := pkgmgr.GetManager(osName, cfg.PackageManager)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		if len(cfg.Packages) == 0 {
			fmt.Println("No packages to install.")
			return
		}

		run := &runner.Runner{DryRun: dryRun}
		for _, p := range cfg.Packages {
			targetPkg := cfg.GetPackageName(p, string(osName))
			
			// Build version map for this specific package
			versions := make(map[string]string)
			if detail, ok := cfg.Mappings[p]; ok && detail.Version != "" {
				versions[targetPkg] = detail.Version
			}

			installCmd := mgr.InstallCommand([]string{targetPkg}, versions)
			
			if err := run.Run(installCmd); err != nil {
				fmt.Printf("[Warning] Failed to install %s: %v\n", p, err)
			} else if dryRun {
				fmt.Printf("[Dry-run] Would install %s\n", p)
			} else {
				fmt.Printf("[Success] %s installed\n", p)
			}
		}
	},
}

func init() {
	applyCmd.Flags().StringVarP(&configURL, "url", "u", "", "URL to a remote setupx.yaml (e.g. GitHub Gist)")
	rootCmd.AddCommand(applyCmd)
}
