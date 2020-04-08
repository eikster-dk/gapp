package cli

import (
	"github.com/eikc/gapp/internal/authentication"
	"github.com/spf13/cobra"
)

func loginCmd() *cobra.Command {
	var initialization = &cobra.Command{
		Use:   "login",
		Short: "login provides the mechanisme to authenticate with github",
		RunE: func(cmd *cobra.Command, args []string) error {
			return authentication.SaveUser()
		},
	}

	return initialization
}