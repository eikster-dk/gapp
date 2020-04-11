package actions

import (
	"encoding/json"
	"fmt"
	"github.com/eikc/gapp/internal/authentication"
	"github.com/google/go-github/v30/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"strings"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "actions",
		Short: "Actions commands will... ",
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
		Use:     "dispatch [owner/repository] [branch] [action]",
		Short:   "Dispatches a command to start a workflow",
		Example: "",
		Args:    cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			branch := args[1]
			action := args[2]

			splitted := strings.Split(path, "/")
			if len(splitted) != 2 {
				return fmt.Errorf("[owner/repository] is not in the correct format")
			}

			owner := splitted[0]
			repo := splitted[1]

			auth, err := authentication.GetUser()
			if err != nil {
				return err
			}

			ctx := cmd.Context()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: auth.Token},
			)

			tc := oauth2.NewClient(ctx, ts)
			client := github.NewClient(tc)
			payload := DispatchPayload{Branch: branch}

			p, _ := json.Marshal(&payload)
			raw := json.RawMessage(p)
			req := github.DispatchRequestOptions{
				EventType:     action,
				ClientPayload: &raw,
			}

			_, _, err = client.Repositories.Dispatch(ctx, owner, repo, req)
			if err != nil {
				return err
			}

			fmt.Println("command dispatch initiated")

			return nil
		},
	}

	return dispatchCmd
}
