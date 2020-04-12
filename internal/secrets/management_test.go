package secrets

import (
	"context"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/eikc/gapp/internal/secrets/mocks"
	"github.com/golang/mock/gomock"
	"gopkg.in/yaml.v2"
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

func TestRunManagement(t *testing.T) {
	tests := []struct {
		Name             string
		Secrets          Secrets
		FileConfig       func(reader *mocks.MockFileReader, secrets Secrets)
		ghMockConfig     func(gh *mocks.MockActionsClient, secrets Secrets)
		encryptionConfig func(writer *mocks.MockEncryptionWriter, secrets Secrets)
		ExpectedError    bool
	}{
		{
			Name: "A successful run of secret management",
			Secrets: Secrets{
				Secrets: []Secret{
					{
						Name:  "test 1",
						Value: "test 1",
						Repos: []string{"eikc/actions-playground", "eikc/gapp"},
					},
					{
						Name:  "testing 2",
						Value: "test 2",
						Repos: []string{"eikc/gapp"},
					},
				},
			},
			FileConfig: func(reader *mocks.MockFileReader, secrets Secrets) {
				bytes, err := yaml.Marshal(&secrets)
				if err != nil {
					t.Fatalf("yaml marshalling failed: %s", err)
				}

				reader.
					EXPECT().
					ReadFile(gomock.Eq("./test.yaml")).
					Return(bytes, nil)
			},
			ghMockConfig: func(gh *mocks.MockActionsClient, secrets Secrets) {
				gh.
					EXPECT().
					GetPublicKey(gomock.Any(), "eikc", "actions-playground").
					Times(1).
					Return([]byte("some pkey 2"), "2", nil)

				gh.
					EXPECT().
					GetPublicKey(gomock.Any(), "eikc", "gapp").
					Times(1).
					Return([]byte("some pkey 3"), "3", nil)

				gh.
					EXPECT().
					AddOrUpdateSecret(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(3)
			},
			encryptionConfig: func(writer *mocks.MockEncryptionWriter, secrets Secrets) {
				writer.
					EXPECT().
					Encrypt(gomock.Eq("test 1"), gomock.Eq([]byte("some pkey 2"))).
					Times(1).
					Return("some encrypted secret 1", nil)

				writer.
					EXPECT().
					Encrypt(gomock.Eq("test 2"), gomock.Eq([]byte("some pkey 3"))).
					Times(1).
					Return("some encrypted secret 2", nil)

				writer.
					EXPECT().
					Encrypt(gomock.Eq("test 1"), gomock.Eq([]byte("some pkey 3"))).
					Times(1).
					Return("some encrypted secret 2", nil)
			},
			ExpectedError: false,
		},
		{
			Name: "Wrongly formatted repository will return an error",
			Secrets: Secrets{
				Secrets: []Secret{
					{
						Name:  "eikc",
						Value: "test",
						Repos: []string{"eikc"},
					},
				},
			},
			FileConfig: func(reader *mocks.MockFileReader, secrets Secrets) {
				bytes, err := yaml.Marshal(&secrets)
				if err != nil {
					t.Fatalf("yaml marshalling failed: %s", err)
				}

				reader.
					EXPECT().
					ReadFile(gomock.Eq("./test.yaml")).
					Return(bytes, nil)
			},
			ghMockConfig: func(gh *mocks.MockActionsClient, secrets Secrets) {

			},
			encryptionConfig: func(writer *mocks.MockEncryptionWriter, secrets Secrets) {

			},
			ExpectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			ghMock := mocks.NewMockActionsClient(ctrl)
			fileReader := mocks.NewMockFileReader(ctrl)
			encryptionMock := mocks.NewMockEncryptionWriter(ctrl)

			test.FileConfig(fileReader, test.Secrets)
			test.ghMockConfig(ghMock, test.Secrets)
			test.encryptionConfig(encryptionMock, test.Secrets)

			cli := NewSecretsCLI(ghMock, fileReader, encryptionMock)

			ctx := context.Background()

			err := cli.RunManagement(ctx, ManagementParams{
				File: "./test.yaml",
			})

			if err != nil {
				if test.ExpectedError == false {
					t.Fatalf("Test failed: %s", err)
				}

				cupaloy.SnapshotT(t, err)
			}

		})
	}
}
