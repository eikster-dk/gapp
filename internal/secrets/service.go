package secrets

import (
	"context"
	"fmt"
	"github.com/eikc/gapp/internal"
	"strings"
)

type Writer interface {
	UpdateSecret(ctx context.Context, owner, repo string, secret internal.Secret) error
}

type Parser interface {
	Parse(path string) (map[string][]internal.Secret, error)
}

type Spinner interface {
	Start() error
	Message(msg string)
	Stop() error
	Fail() error
}

type Service struct {
	writer  Writer
	parser  Parser
	spinner Spinner
}

func NewService(writer Writer, parser Parser, spinner Spinner) *Service {
	return &Service{
		writer:  writer,
		parser:  parser,
		spinner: spinner,
	}
}

type ManagementParams struct {
	Path string
}

func (s *Service) RunManagement(ctx context.Context, params ManagementParams) error {
	s.spinner.Start()

	secrets, err := s.parser.Parse(params.Path)
	if err != nil {
		s.spinner.Fail()
		return err
	}

	for repoIdentifier, ss := range secrets {
		splitted := strings.Split(repoIdentifier, "/")
		if len(splitted) < 2 {
			return fmt.Errorf("repository is not correctly formattted, use [owner]/[repository] pattern, got: %s",
				splitted[0])
		}

		owner := splitted[0]
		repo := splitted[1]

		for _, secret := range ss {
			s.spinner.Message(fmt.Sprintf("repo: %s secret: %s", repoIdentifier, secret.Name))
			err = s.writer.UpdateSecret(ctx, owner, repo, secret)
			if err != nil {
				s.spinner.Fail()
				return err
			}
		}
	}

	return s.spinner.Stop()
}
