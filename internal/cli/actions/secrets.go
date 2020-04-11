package actions

import (
	"github.com/eikc/gapp/internal/authentication"
	"github.com/eikc/gapp/internal/gh"
	"github.com/eikc/gapp/internal/secrets"
	"github.com/spf13/cobra"
)

func secretCmd() *cobra.Command {
	var secrets = &cobra.Command{
		Short: "Secrets will update one or more secrets within multiple repositories",
		Use:   "secrets",
	}
	secrets.AddCommand(manageSecrets())

	return secrets
}

func manageSecrets() *cobra.Command {
	var management = &cobra.Command{
		Short: "management will create or update secrets based on a yaml manifest file",
		Use:   "management",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			loc, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}

			user, err := authentication.GetUser()
			if err != nil {
				return err
			}

			client := gh.NewActionsClient(ctx, user)
			encryptionWriter := secrets.NewEncrypt()

			cli := secrets.NewSecretsCLI(client, nil, encryptionWriter )

			return cli.RunManagement(ctx, secrets.ManagementParams{
				File: loc,
			})
		},
	}

	management.Flags().StringP("file", "f", "", "location of the secrets.yaml file")
	management.MarkFlagRequired("file")

	return management
}
