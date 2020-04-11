package secrets

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

func TestSortSecrets(t *testing.T) {
	secrets := Secrets{
		Secrets: []Secret{
			{
				Name:  "Test secret 1",
				Value: "123321",
				Repos: []string{"eikc/masscommerce", "eikc/gapp"},
			},
			{
				Name:  "Test secret 2",
				Value: "123321",
				Repos: []string{"eikc/masscommerce", "eikc/gapp"},
			},
			{
				Name:  "test secret 3",
				Value: "123321",
				Repos: []string{"eikc/gapp"},
			},
		},
	}

	sorted := sortSecrets(secrets)

	cupaloy.SnapshotT(t, sorted)
}