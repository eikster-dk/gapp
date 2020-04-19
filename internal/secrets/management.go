package secrets

import (
	"context"
	"fmt"
	"github.com/eikc/gapp/internal/files"
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
	Path string
}

func (cli *CLI) RunManagement(ctx context.Context, params ManagementParams) error {
	cli.spinner.Start()

	sortedSecrets, err := cli.readAndParse(params)
	if err != nil {
		cli.spinner.Fail()
		return err
	}

	err = cli.updateSecrets(ctx, sortedSecrets)
	if err != nil {
		cli.spinner.Fail()
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
			cli.spinner.Message(fmt.Sprintf("repo: %s secret: %s", repo, secret.Name))
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

func (cli *CLI) readAndParse(params ManagementParams) (map[string][]Secret, error) {
	isFile, err := cli.reader.IsFile(params.Path)
	if err != nil {
		return nil, err
	}

	if isFile {
		secrets, err := parseSecrets(params.Path, cli.reader)
		if err != nil {
			return nil, err
		}

		sortedSecrets := sortSecrets(secrets)
		return sortedSecrets, nil
	} else {
		paths, err := cli.reader.ReadDir(params.Path)
		if err != nil {
			return nil, err
		}

		var secrets []Secret
		for _, p := range paths {
			ss, err := parseSecrets(p, cli.reader)
			if err != nil {
				return nil, err
			}

			secrets = append(secrets, ss.Secrets...)
		}

		sortedSecrets := sortSecrets(Secrets{secrets})

		return sortedSecrets, nil
	}
}

func parseSecrets(path string, reader files.ReadFile) (Secrets, error) {
	content, err := reader.ReadFile(path)
	if err != nil {
		return Secrets{}, err
	}

	var secrets Secrets
	err = yaml.Unmarshal(content, &secrets)
	if err != nil {
		return Secrets{}, err
	}

	return secrets, nil
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
