package cmd

import (
	"fmt"
	"log"
	"os"
	"setupx/internal/config"
	"setupx/internal/pkgmgr"
	"setupx/internal/runner"
	"strings"

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

		run := &runner.Runner{DryRun: dryRun}

		// 1. Try Exact Search first
		exactCmd := mgr.ExactSearchCommand(pkg)
		out, _ := run.RunOutput(exactCmd)
		results := mgr.ParseSearchOutput(out)

		// 2. If no exact results, or very few, try broad search
		if len(results) == 0 {
			broadCmd := mgr.SearchCommand(pkg)
			out, _ = run.RunOutput(broadCmd)
			results = mgr.ParseSearchOutput(out)
		}

		if len(results) == 0 {
			fmt.Printf("No results found for '%s'.\n", pkg)
			return
		}

		// 3. Render Table
		renderTable(results)
	},
}

func renderTable(results []pkgmgr.SearchResult) {
	const nameWidth = 35
	fmt.Printf("%-35s %s\n", "NAME", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 70))
	for _, r := range results {
		name := r.Name
		if len(name) > nameWidth {
			name = name[:nameWidth-3] + "..."
		}
		desc := r.Description
		if len(desc) > 32 {
			desc = desc[:29] + "..."
		}
		fmt.Printf("%-35s %s\n", name, desc)
	}
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
