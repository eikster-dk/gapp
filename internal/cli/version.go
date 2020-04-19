package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	// Build, Major, Minor, Patch and Label are set at build time
	Major, Minor, Patch, Label string
)

func init() {
	if Major == "" {
		Major = "0"
	}
	if Minor == "" {
		Minor = "0"
	}
	if Patch == "" {
		Patch = "0"
	}
	if Label == "" {
		Label = "dev"
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the gapp version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gapp version %s.%s.%s-%s \n", Major, Minor, Patch, Label)
	},
}
