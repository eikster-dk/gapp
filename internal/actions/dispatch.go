package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type DispatchParams struct {
	Event   string
	Repo    string
	Payload map[string]string
}

func (cli *CLI) Dispatch(ctx context.Context, params DispatchParams) error {
	splitted := strings.Split(params.Repo, "/")
	if len(splitted) < 2 {
		return fmt.Errorf("[owner/repository] is not in the correct format, got %s", splitted[0])
	}

	owner := splitted[0]
	repo := splitted[1]

	payload, err := json.Marshal(params.Payload)
	if err != nil {
		return err
	}

	if err = cli.ghClient.Dispatch(ctx, owner, repo, params.Event, payload); err != nil {
		return err
	}

	fmt.Fprintf(cli.stdOut, "%s dispatch initiated \n", params.Event)

	return nil
}
