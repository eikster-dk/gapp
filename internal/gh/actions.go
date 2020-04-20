package gh

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eikc/gapp/internal/authentication"
	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

type ActionsClient struct {
	client *github.Client
	pkeys  map[string]*github.PublicKey
}

func (c *ActionsClient) Dispatch(ctx context.Context, owner, repo, event string, payload []byte) error {
	raw := json.RawMessage(payload)
	req := github.DispatchRequestOptions{
		EventType:     event,
		ClientPayload: &raw,
	}

	_, _, err := c.client.Repositories.Dispatch(ctx, owner, repo, req)

	return err
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
		pkeys:  map[string]*github.PublicKey{},
	}
}

func (c *ActionsClient) GetPublicKey(ctx context.Context, owner, repo string) ([]byte, string, error) {
	repoUrl := fmt.Sprintf("%s/%s", owner, repo)

	val, ok := c.pkeys[repoUrl]
	if !ok {
		key, _, err := c.client.Actions.GetPublicKey(ctx, owner, repo)
		if err != nil {
			return nil, "", err
		}

		c.pkeys[repoUrl] = key
		val = key
	}

	decoded, err := base64.StdEncoding.DecodeString(*val.Key)
	if err != nil {
		return nil, "", err
	}

	return decoded, *val.KeyID, nil
}

func (c *ActionsClient) AddOrUpdateSecret(ctx context.Context, owner, repo string, secret SecretParams) error {

	_, err := c.client.Actions.CreateOrUpdateSecret(ctx, owner, repo, &github.EncryptedSecret{
		Name:           secret.Name,
		KeyID:          secret.PkeyId,
		EncryptedValue: secret.Value,
	})

	return err
}
