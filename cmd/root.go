package cmd

import (
	"github.com/spf13/cobra"
)

var dryRun bool
var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "setupx",
	Version: Version,
	Short:   "setupx is a cross-platform package manager",
	Long:    `A simple, fast tool to install packages across macOS, Linux, and Windows.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry run (don't execute commands)")
}
