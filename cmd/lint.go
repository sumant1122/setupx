package cmd

import (
	"fmt"
	"setupx/internal/config"

	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Validate the setupx.yaml configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("setupx.yaml")
		if err != nil {
			return fmt.Errorf("failed to load or parse config: %w", err)
		}

		fmt.Println("Validating setupx.yaml...")

		if len(cfg.Packages) == 0 {
			fmt.Println("[-Warning-] No packages defined in the list.")
		} else {
			fmt.Printf("[OK] %d packages defined.\n", len(cfg.Packages))
		}

		// Check for dangling mappings
		for k := range cfg.Mappings {
			found := false
			for _, p := range cfg.Packages {
				if k == p {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("[-Warning-] Mapping for '%s' exists, but it is not in the packages list.\n", k)
			}
		}

		fmt.Println("Configuration is valid.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
