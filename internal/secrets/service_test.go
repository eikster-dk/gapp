package secrets

import (
	"context"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/eikc/gapp/internal"
	"github.com/eikc/gapp/internal/secrets/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestService_RunManagement(t *testing.T) {
	tests := []struct {
		name            string
		args            ManagementParams
		configuerWriter func(w *mocks.MockWriter)
		configureParser func(p *mocks.MockParser)
	}{
		{
			name: "Given multiple correct secrets, it will update them all in github",
			args: ManagementParams{
				Path: "some path",
			},
			configuerWriter: func(w *mocks.MockWriter) {
				w.EXPECT().UpdateSecret(gomock.Any(), "eikc", "gapp", internal.Secret{
					Name:  "test secret",
					Value: "test 123",
				})
			},
			configureParser: func(p *mocks.MockParser) {
				returned := map[string][]internal.Secret{
					"eikc/gapp": {
						{
							Name:  "test secret",
							Value: "test 123",
						},
					},
				}

				p.EXPECT().Parse("some path").Times(1).Return(returned, nil)
			},
		},
		{
			name: "Given an incorrectly formatted repo, it will return an error",
			args: ManagementParams{
				Path: "some path",
			},
			configuerWriter: func(w *mocks.MockWriter) {

			},
			configureParser: func(p *mocks.MockParser) {
				returned := map[string][]internal.Secret{
					"eikc-gapp": {
						{
							Name:  "test secret",
							Value: "test 123",
						},
					},
				}

				p.EXPECT().Parse("some path").Times(1).Return(returned, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			writer := mocks.NewMockWriter(ctrl)
			parser := mocks.NewMockParser(ctrl)
			spinner := mocks.NoopSpinner{}

			tt.configuerWriter(writer)
			tt.configureParser(parser)

			s := NewService(writer, parser, spinner)

			err := s.RunManagement(ctx, tt.args)

			cupaloy.SnapshotT(t, err)
		})
	}
}
