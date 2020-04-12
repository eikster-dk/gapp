package actions

import (
	"github.com/eikc/gapp/internal/actions"
	"github.com/eikc/gapp/internal/authentication"
	"github.com/eikc/gapp/internal/gh"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actions",
		Short: "Commands for github actions",
	}

	cmd.AddCommand(dispatchCmd())
	cmd.AddCommand(secretCmd())

	return cmd
}

type DispatchPayload struct {
	Branch string `json:"branch"`
}

func dispatchCmd() *cobra.Command {
	var dispatchCmd = &cobra.Command{
		Use:     "dispatch [owner/repository]",
		Short:   "Dispatches a command to start a workflow",
		Example: "",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			event, err := cmd.Flags().GetString("event")
			if err != nil {
				return err
			}

			payload, err := cmd.Flags().GetStringToString("payload")
			if err != nil {
				return err
			}

			branch, err := cmd.Flags().GetString("branch")
			if err != nil {
				return err
			}

			if branch != "" {
				payload["branch"] = branch
			}

			repo := args[0]

			user, err := authentication.GetUser()
			if err != nil {
				return err
			}

			client := gh.NewActionsClient(ctx, user)
			cli := actions.NewCLI(cmd.OutOrStdout(), client)

			return cli.Dispatch(ctx, actions.DispatchParams{
				Event:   event,
				Repo:    repo,
				Payload: payload,
			})
		},
	}

	dispatchCmd.Flags().StringToStringP("payload", "p", nil, "payload is used to add a additional json payload to the event")
	dispatchCmd.Flags().StringP("branch", "b", "", "Branch is a shortcut to add a branch value to the json payload")
	dispatchCmd.Flags().StringP("event", "e", "", "Event is the event type to dispatch to github")

	dispatchCmd.MarkFlagRequired("event")

	return dispatchCmd
}
