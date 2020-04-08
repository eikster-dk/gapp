package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os/exec"
)

func Do(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{Use: "gapp", SilenceUsage: true}
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(actionsCmd())
	rootCmd.AddCommand(loginCmd())
	rootCmd.AddCommand(repoCmd())

	// secrets should probably be a sub command under actions
	rootCmd.AddCommand(secretCmd())


	err := rootCmd.Execute()
	if err != nil {
		return 0
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}

	return 1
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the gapp version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Coming soon :-)")
	},
}
