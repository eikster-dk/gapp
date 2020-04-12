package cli

import (
	"github.com/eikc/gapp/internal/authentication"
	"github.com/spf13/cobra"
)

func loginCmd() *cobra.Command {
	var initialization = &cobra.Command{
		Use:   "login",
		Short: "Login initiates gapp by allowing you to provide a personal access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return authentication.SaveUser()
		},
	}

	return initialization
}