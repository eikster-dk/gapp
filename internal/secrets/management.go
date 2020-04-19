package secrets

import (
	"context"
)

type ManagementParams struct {
	Path string
}

func (cli *CLI) RunManagement(ctx context.Context, params ManagementParams) error {
	cli.spinner.Start()

	sortedSecrets, err := cli.parser.Parse(params.Path)
	if err != nil {
		cli.spinner.Fail()
		return err
	}

	err = cli.client.updateSecrets(ctx, sortedSecrets)
	if err != nil {
		cli.spinner.Fail()
		return err
	}

	return cli.spinner.Stop()
}
