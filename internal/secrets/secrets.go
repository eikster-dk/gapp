package secrets

import (
	"context"
	"github.com/eikc/gapp/internal/gh"
)

type CLI struct {
	reader           FileReader
	githubClient     ActionsClient
	encryptionWriter EncryptionWriter
	spinner          Spinner
}

type FileReader interface {
	IsFile(path string) (bool, error)
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]string, error)
}

type ActionsClient interface {
	GetPublicKey(ctx context.Context, owner, repo string) ([]byte, string, error)
	AddOrUpdateSecret(ctx context.Context, owner, repo string, secret gh.SecretParams) error
}

type EncryptionWriter interface {
	Encrypt(value string, pkey []byte) (string, error)
}

type Spinner interface {
	Start() error
	Message(msg string)
	Stop() error
	Fail() error
}

func NewSecretsCLI(ghActions ActionsClient, reader FileReader, writer EncryptionWriter, spinner Spinner) *CLI {
	return &CLI{
		reader:           reader,
		githubClient:     ghActions,
		encryptionWriter: writer,
		spinner:          spinner,
	}
}

