package secrets

import (
	"context"
)

type Writer interface {
	updateSecrets(ctx context.Context, sortedSecrets map[string][]Secret) error
}

type Parser interface {
	Parse(path string) (map[string][]Secret, error)
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

	err = s.writer.updateSecrets(ctx, secrets)
	if err != nil {
		s.spinner.Fail()
		return err
	}

	return s.spinner.Stop()
}
