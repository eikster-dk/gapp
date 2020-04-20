package secrets

import (
	"context"
	"github.com/eikc/gapp/internal"
	"github.com/eikc/gapp/internal/gh"
)

type ActionsClient interface {
	GetPublicKey(ctx context.Context, owner, repo string) ([]byte, string, error)
	AddOrUpdateSecret(ctx context.Context, owner, repo string, secret gh.SecretParams) error
}

type EncryptionWriter interface {
	Encrypt(value string, pkey []byte) (string, error)
}

type writer struct {
	client ActionsClient
	writer EncryptionWriter
}

func (w *writer) UpdateSecret(ctx context.Context, owner, repo string, secret internal.Secret) error {
	pkey, pkeyId, err := w.client.GetPublicKey(ctx, owner, repo)
	if err != nil {
		return err
	}

	encoded, err := w.writer.Encrypt(secret.Value, pkey)
	if err != nil {
		return err
	}

	err = w.client.AddOrUpdateSecret(ctx, owner, repo, gh.SecretParams{
		Name:   secret.Name,
		PkeyId: pkeyId,
		Value:  encoded,
	})

	if err != nil {
		return err
	}

	return nil
}

func NewWriter(client ActionsClient, encryptor EncryptionWriter) *writer {
	return &writer{
		client: client,
		writer: encryptor,
	}
}
