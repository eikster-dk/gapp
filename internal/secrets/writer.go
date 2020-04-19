package secrets

import (
	"context"
	"fmt"
	"github.com/eikc/gapp/internal/gh"
	"strings"
)

type ActionsClient interface {
	GetPublicKey(ctx context.Context, owner, repo string) ([]byte, string, error)
	AddOrUpdateSecret(ctx context.Context, owner, repo string, secret gh.SecretParams) error
}

type EncryptionWriter interface {
	Encrypt(value string, pkey []byte) (string, error)
}

type writer struct {
	client  ActionsClient
	writer  EncryptionWriter
	spinner Spinner
}

func NewWriter(client ActionsClient, encryptor EncryptionWriter, spinner Spinner) *writer {
	return &writer{
		client:  client,
		writer:  encryptor,
		spinner: spinner,
	}
}

func (g *writer) updateSecrets(ctx context.Context, sortedSecrets map[string][]Secret) error {
	for repo, secrets := range sortedSecrets {
		splitted := strings.Split(repo, "/")
		if len(splitted) < 2 {
			return fmt.Errorf("repository is not correctly formattted, use [owner]/[repository] pattern, got: %s",
				splitted[0])
		}

		owner := splitted[0]
		repo := splitted[1]
		pkey, pkeyId, err := g.client.GetPublicKey(ctx, owner, repo)
		if err != nil {
			return err
		}

		for _, secret := range secrets {
			g.spinner.Message(fmt.Sprintf("repo: %s secret: %s", repo, secret.Name))
			encoded, err := g.writer.Encrypt(secret.Value, pkey)
			if err != nil {
				return err
			}

			err = g.client.AddOrUpdateSecret(ctx, owner, repo, gh.SecretParams{
				Name:   secret.Name,
				PkeyId: pkeyId,
				Value:  encoded,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
