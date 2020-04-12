package actions

import (
	"context"
	"io"
)

type GithubClient interface {
	Dispatch(ctx context.Context, owner, repo, event string, payload []byte) error
}

type CLI struct {
	stdOut   io.Writer
	ghClient GithubClient
}

func NewCLI(stdOut io.Writer, client GithubClient) *CLI {
	return &CLI{
		stdOut:   stdOut,
		ghClient: client,
	}
}
