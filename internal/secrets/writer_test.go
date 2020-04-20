package secrets

import (
	"context"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/eikc/gapp/internal"
	"github.com/eikc/gapp/internal/gh"
	"github.com/eikc/gapp/internal/secrets/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_github_updateSecrets(t *testing.T) {
	type args struct {
		owner  string
		repo   string
		secret internal.Secret
	}
	tests := []struct {
		name                      string
		params                    args
		configureClient           func(gh *mocks.MockActionsClient)
		configureEncryptionWriter func(writer *mocks.MockEncryptionWriter)
	}{
		{
			name: "When All secrets are correctly formatted, updateSecrets will succeed",
			params: args{
				owner: "eikc",
				repo:  "gapp",
				secret: internal.Secret{
					Name:  "top secret",
					Value: "top secret VALUE!",
				},
			},
			configureClient: func(client *mocks.MockActionsClient) {
				client.EXPECT().
					GetPublicKey(gomock.Any(), "eikc", "gapp").
					Times(1).
					Return([]byte("testing"), "1", nil)

				client.EXPECT().
					AddOrUpdateSecret(gomock.Any(), "eikc", "gapp", gh.SecretParams{
						Name:   "top secret",
						Value:  "encrypted",
						PkeyId: "1",
					}).
					Times(1).
					Return(nil)
			},
			configureEncryptionWriter: func(writer *mocks.MockEncryptionWriter) {
				writer.EXPECT().Encrypt("top secret VALUE!", []byte("testing")).Times(1).Return("encrypted", nil)
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

			cli := &writer{
				client: mockClient,
				writer: mockEncryption,
			}

			err := cli.updateSecret(ctx, tt.params.owner, tt.params.repo, tt.params.secret)

			cupaloy.SnapshotT(t, err)
		})
	}

}
