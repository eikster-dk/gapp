package secrets

import (
	"context"
	"github.com/eikc/gapp/internal/gh"
)

type CLI struct {
	reader           FileReader
	githubClient     ActionsClient
	encryptionWriter EncryptionWriter
}

type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

type ActionsClient interface {
	GetPublicKey(ctx context.Context, owner, repo string) ([]byte, string, error)
	AddOrUpdateSecret(ctx context.Context, owner, repo string, secret gh.SecretParams) error
}

type EncryptionWriter interface {
	Encrypt(value string, pkey []byte) (string, error)
}

func NewSecretsCLI(ghActions ActionsClient, reader FileReader, writer EncryptionWriter) *CLI {
	return &CLI{
		reader:           reader,
		githubClient:     ghActions,
		encryptionWriter: writer,
	}
}
