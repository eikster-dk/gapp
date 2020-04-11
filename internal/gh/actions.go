package gh

import (
	"context"
	"encoding/base64"
	"github.com/eikc/gapp/internal/authentication"
	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

type ActionsClient struct {
	client *github.Client
}

type SecretParams struct {
	Name   string
	Value  string
	PkeyId string
}

func NewActionsClient(ctx context.Context, user authentication.User) *ActionsClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.Token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &ActionsClient{
		client: client,
	}
}

func (c *ActionsClient) GetPublicKey(ctx context.Context, owner, repo string) ([]byte, string, error) {
	pkey, _, err := c.client.Actions.GetPublicKey(ctx, owner, repo)
	if err != nil {
		return nil, nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(*pkey.Key)
	if err != nil {
		return nil, nil, err
	}

	return decoded, *pkey.KeyID, nil
}

func (c *ActionsClient) AddOrUpdateSecret(ctx context.Context, owner, repo string, secret SecretParams) error {

	_, err := c.client.Actions.CreateOrUpdateSecret(ctx, owner, repo, &github.EncryptedSecret{
		Name:           secret.Name,
		KeyID:          secret.PkeyId,
		EncryptedValue: secret.Value,
	})

	return err
}
