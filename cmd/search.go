package cmd

import (
	"setupx/internal/config"
	"setupx/internal/pkgmgr"
	"setupx/internal/runner"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [pkg]",
	Short: "Search for a package in the native package manager",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]
		
		cfg, err := config.LoadConfig("setupx.yaml")
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("Error loading config: %v", err)
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

		searchCmd := mgr.SearchCommand(pkg)
		run := &runner.Runner{DryRun: dryRun}
		
		if dryRun {
			// In search, dry run should probably still show the command
			log.Printf("[Dry-run] Would run search: %v", searchCmd)
			return
		}

		if err := run.Run(searchCmd); err != nil {
			log.Fatalf("Error executing search: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
