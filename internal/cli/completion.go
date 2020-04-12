package cli

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

func (c *CLI) completion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long:  `Generate shell completion scripts for gapp CLI commands.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			shellType, err := cmd.Flags().GetString("shell")
			if err != nil {
				return err
			}

			if shellType == "" {
				return errors.New("error: the value for `--shell` is required\nsee `gapp help completion` for more information")
			}

			switch shellType {
			case "bash":
				return c.rootCmd.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return c.rootCmd.GenZshCompletion(cmd.OutOrStdout())
			case "powershell":
				return c.rootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell type %q", shellType)
			}
		},
	}

	cmd.Flags().StringP("shell", "s", "bash", "Shell type for auto completion")

	return cmd
}
