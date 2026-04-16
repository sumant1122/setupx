package cmd

import (
	"setupx/internal/config"
	"setupx/internal/pkgmgr"
	"setupx/internal/runner"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type pkgResult struct {
	Name        string
	Description string
}

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
		results := parseResults(out, string(osName))

		// 2. If no exact results, or very few, try broad search
		if len(results) == 0 {
			broadCmd := mgr.SearchCommand(pkg)
			out, _ = run.RunOutput(broadCmd)
			results = parseResults(out, string(osName))
		}

		if len(results) == 0 {
			fmt.Printf("No results found for '%s'.\n", pkg)
			return
		}

		// 3. Render Table
		renderTable(results)
	},
}

func parseResults(out string, osName string) []pkgResult {
	var results []pkgResult
	lines := strings.Split(out, "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "Sorting...") || strings.HasPrefix(line, "Full Text Search...") || strings.Contains(line, "WARNING") {
			continue
		}

		// Simple heuristics for common package managers
		switch osName {
		case "linux": // Apt/Dnf
			parts := strings.SplitN(line, "/", 2)
			if len(parts) == 2 {
				name := parts[0]
				desc := ""
				if i+1 < len(lines) {
					desc = strings.TrimSpace(lines[i+1])
					i++ // skip next line as it's the description
				}
				results = append(results, pkgResult{Name: name, Description: desc})
			}
		case "mac": // Brew
			// Brew search just returns names
			results = append(results, pkgResult{Name: line, Description: "N/A"})
		case "windows": // Winget
			// Winget output format: Name | Id | Version | Source
			// We want the 'Id' (column 2) because that's what is used for install/mappings
			if strings.HasPrefix(line, "Name") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "No package found") {
				continue
			}
			
			// Winget uses fixed-width columns usually, but fields can help if there are no spaces in names.
			// However, IDs never have spaces. 
			// Heuristic: The ID is usually the second 'blob' of text that looks like a package identifier (e.g. Neovim.Neovim)
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				// ID is almost always the second field in winget search
				id := fields[len(fields)-3] // Backwards heuristic often works better for winget if names have spaces
				if len(fields) == 2 {
					id = fields[1]
				} else if len(fields) >= 4 {
					// Name ID Version Source -> ID is index 1 or 2 depending on spaces
					// Let's look for the one with a dot or capitalized pattern
					id = fields[1] 
					for _, f := range fields {
						if strings.Contains(f, ".") {
							id = f
							break
						}
					}
				}
				results = append(results, pkgResult{Name: id, Description: "Available on Windows"})
			}
		default:
			results = append(results, pkgResult{Name: line, Description: ""})
		}

		if len(results) >= 20 {
			break
		}
	}
	return results
}

func renderTable(results []pkgResult) {
	const nameWidth = 20
	fmt.Printf("%-20s %s\n", "NAME", "DESCRIPTION")
	fmt.Println(strings.Repeat("-", 60))
	for _, r := range results {
		name := r.Name
		if len(name) > nameWidth {
			name = name[:nameWidth-3] + "..."
		}
		desc := r.Description
		if len(desc) > 37 {
			desc = desc[:34] + "..."
		}
		fmt.Printf("%-20s %s\n", name, desc)
	}
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
