package cmd

import (
	"setupx/internal/config"
	"setupx/internal/models"
	"setupx/internal/pkgmgr"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain [pkg]",
	Short: "Explain what command would be run for a package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]
		
		var cfg *models.Config
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

		var targetPkg string
		if cfg != nil {
			targetPkg = cfg.GetPackageName(pkg, string(osName))
		} else {
			targetPkg = pkg
		}

		installCmd := mgr.InstallCommand([]string{targetPkg})
		
		fmt.Printf("Package: %s\n", pkg)
		fmt.Printf("Detected OS: %s\n", osName)
		fmt.Printf("Mapped Name: %s\n", targetPkg)
		fmt.Printf("Command: %s\n", pkgmgr.FormatCommand(installCmd))
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
