package cli

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/eikc/gapp/internal/authentication"
	"github.com/google/go-github/v30/github"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

func secretCmd() *cobra.Command {
	var secrets = &cobra.Command{
		Short: "Secrets will update one or more secrets within multiple repositories",
		Use:   "secrets",
		RunE:  secrets,
	}

	secrets.Flags().StringP("file", "f", "", "location of the secrets.yaml file")
	secrets.MarkFlagRequired("file")

	return secrets
}

func secrets(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	user, err := authentication.GetUser()
	if err != nil {
		return err
	}

	loc, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	type Secret struct {
		Name  string   `yaml:"name"`
		Value string   `yaml:"value"`
		Repos []string `yaml:"repos"`
	}
	type Secrets struct {
		Secrets []Secret `yaml:"secrets"`
	}

	content, err := ioutil.ReadFile(loc)
	if err != nil {
		return err
	}

	var secrets Secrets
	err = yaml.Unmarshal(content, &secrets)
	if err != nil {
		return err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: user.Token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for _, secret := range secrets.Secrets {
		if len(secret.Repos) == 0 {
			fmt.Printf("secret %s has no repos\n", secret.Name)
			continue
		}

		for _, r := range secret.Repos {
			splitted := strings.Split(r, "/")
			owner := splitted[0]
			repo := splitted[1]

			pkey, _, err := client.Actions.GetPublicKey(ctx, owner, repo)
			if err != nil {
				return err
			}
			decoded, err := base64.StdEncoding.DecodeString(*pkey.Key)
			if err != nil {
				return err
			}

			msg := []byte(secret.Value)
			var key [32]byte
			copy(key[:], decoded)

			var out []byte

			encrypted, err := box.SealAnonymous(out, msg, &key, rand.Reader)
			if err != nil {
				return err
			}

			encoded := base64.StdEncoding.EncodeToString(encrypted)

			_, err = client.Actions.CreateOrUpdateSecret(ctx, owner, repo, &github.EncryptedSecret{
				Name:           secret.Name,
				KeyID:          *pkey.KeyID,
				EncryptedValue: encoded,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
