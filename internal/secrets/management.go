package secrets

import (
	"context"
	"fmt"
	"github.com/eikc/gapp/internal/gh"
	"gopkg.in/yaml.v2"
	"strings"
)

type Secrets struct {
	Secrets []Secret `yaml:"secrets"`
}

type Secret struct {
	Name  string   `yaml:"name"`
	Value string   `yaml:"value"`
	Repos []string `yaml:"repos"`
}

type ManagementParams struct {
	File string
}

func (cli *CLI) RunManagement(ctx context.Context, params ManagementParams) error {
	sortedSecrets, err := cli.parseAndSort(params)
	if err != nil {
		return err
	}

	err = cli.updateSecrets(ctx, sortedSecrets)
	if err != nil {
		return err
	}

	return nil
}

func (cli *CLI) updateSecrets(ctx context.Context, sortedSecrets map[string][]Secret) error {
	for repo, secrets := range sortedSecrets {
		splitted := strings.Split(repo, "/")
		if len(splitted) < 2 {
			return fmt.Errorf("repository is not correctly formattted, use [owner]/[repository] pattern, got: %s",
				splitted[0])
		}

		owner := splitted[0]
		repo := splitted[1]
		pkey, pkeyId, err := cli.githubClient.GetPublicKey(ctx, owner, repo)
		if err != nil {
			return err
		}

		for _, secret := range secrets {
			encoded, err := cli.encryptionWriter.Encrypt(secret.Value, pkey)
			if err != nil {
				return err
			}

			err = cli.githubClient.AddOrUpdateSecret(ctx, owner, repo, gh.SecretParams{
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

func (cli *CLI) parseAndSort(params ManagementParams) (map[string][]Secret, error) {
	content, err := cli.reader.ReadFile(params.File)
	if err != nil {
		return nil, err
	}

	var secrets Secrets
	err = yaml.Unmarshal(content, &secrets)
	if err != nil {
		return nil, err
	}

	sortedSecrets := sortSecrets(secrets)
	return sortedSecrets, nil
}

func sortSecrets(secrets Secrets) map[string][]Secret {
	mapped := make(map[string][]Secret)
	for _, secret := range secrets.Secrets {
		for _, r := range secret.Repos {
			rr := mapped[r]
			rr = append(rr, secret)

			mapped[r] = rr
		}
	}
	return mapped
}
