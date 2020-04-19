package secrets

import (
	"context"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/eikc/gapp/internal/gh"
	"github.com/eikc/gapp/internal/secrets/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_github_updateSecrets(t *testing.T) {
	tests := []struct {
		name                      string
		params                    map[string][]Secret
		configureClient           func(gh *mocks.MockActionsClient)
		configureEncryptionWriter func(writer *mocks.MockEncryptionWriter)
	}{
		{
			name: "When All secrets are correctly formatted, updateSecrets will succeed",
			params: map[string][]Secret{
				"eikc/gapp": {
					{
						Name:  "test secret",
						Value: "test value",
					},
					{
						Name:  "test secret 2",
						Value: "test value 2",
					},
				},
			},
			configureClient: func(client *mocks.MockActionsClient) {
				client.EXPECT().
					GetPublicKey(gomock.Any(), "eikc", "gapp").
					Times(1).
					Return([]byte("testing"), "1", nil)

				client.EXPECT().
					AddOrUpdateSecret(gomock.Any(), "eikc", "gapp", gh.SecretParams{
						Name:   "test secret",
						Value:  "1",
						PkeyId: "1",
					}).
					Times(1).
					Return(nil)

				client.EXPECT().
					AddOrUpdateSecret(gomock.Any(), "eikc", "gapp", gh.SecretParams{
						Name:   "test secret 2",
						Value:  "2",
						PkeyId: "1",
					}).
					Times(1).
					Return(nil)
			},
			configureEncryptionWriter: func(writer *mocks.MockEncryptionWriter) {
				writer.EXPECT().Encrypt("test value", []byte("testing")).Times(1).Return("1", nil)
				writer.EXPECT().Encrypt("test value 2", []byte("testing")).Times(1).Return("2", nil)
			},
		},
		{
			name: "Will return an error when repository is incorrectly formatted",
			params: map[string][]Secret{
				"incorrect-format": {
					{
						Name:  "test",
						Value: "test 2",
					},
				},
			},
			configureClient: func(gh *mocks.MockActionsClient) {

			},
			configureEncryptionWriter: func(writer *mocks.MockEncryptionWriter) {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			mockClient := mocks.NewMockActionsClient(ctrl)
			mockEncryption := mocks.NewMockEncryptionWriter(ctrl)

			tt.configureClient(mockClient)
			tt.configureEncryptionWriter(mockEncryption)

			cli := &Client{
				client:  mockClient,
				writer:  mockEncryption,
				spinner: mocks.NoopSpinner{},
			}

			err := cli.updateSecrets(ctx, tt.params)

			cupaloy.SnapshotT(t, err)
		})
	}

}
