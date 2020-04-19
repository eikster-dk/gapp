package secrets

import (
	"context"
)

type Spinner interface {
	Start() error
	Message(msg string)
	Stop() error
	Fail() error
}

type client interface {
	updateSecrets(ctx context.Context, sortedSecrets map[string][]Secret) error
}

type parser interface {
	Parse(path string) (map[string][]Secret, error)
}

type Secret struct {
	Name  string
	Value string
}

type CLI struct {
	parser  parser
	client  client
	spinner Spinner
}

func NewSecretsCLI(parser parser, client client, spinner Spinner) *CLI {
	return &CLI{
		parser:  parser,
		client:  client,
		spinner: spinner,
	}
}
