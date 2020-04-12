package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/eikc/gapp/internal/actions/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestCLI_Dispatch(t *testing.T) {
	tests := []struct {
		name          string
		params        DispatchParams
		configGh      func(gh *mocks.MockGithubClient, params DispatchParams)
		expectedError bool
	}{
		{
			name: "Dispatch will succesfully dispatch correct event",
			params: DispatchParams{
				Event: "deploy-test",
				Repo:  "eikc/actions-playground",
				Payload: map[string]string{
					"branch": "testing",
				},
			},
			configGh: func(gh *mocks.MockGithubClient, params DispatchParams) {
				payload, _ := json.Marshal(params.Payload)

				gh.
					EXPECT().
					Dispatch(gomock.Any(), "eikc", "actions-playground", "deploy-test", gomock.Eq(payload)).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Wrongly formatted repo will result in error",
			params: DispatchParams{
				Event:   "test-event",
				Repo:    "eikc",
				Payload: nil,
			},
			configGh: func(gh *mocks.MockGithubClient, params DispatchParams) {

			},
			expectedError: true,
		},
		{
			name: "gh returns an error",
			params: DispatchParams{
				Event:   "test-event",
				Repo:    "eikc/actions-playground",
				Payload: nil,
			},
			configGh: func(gh *mocks.MockGithubClient, params DispatchParams) {
				gh.
					EXPECT().
					Dispatch(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("repository does not exist"))
			},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			var stdOut bytes.Buffer

			ghClient := mocks.NewMockGithubClient(ctrl)
			tt.configGh(ghClient, tt.params)

			cli := NewCLI(&stdOut, ghClient)
			ctx := context.Background()

			if err := cli.Dispatch(ctx, tt.params); err != nil {
				if tt.expectedError == false {
					t.Errorf("Dispatch() error = %v", err)
				}

				cupaloy.SnapshotT(t, err.Error())
				return
			}

			cupaloy.SnapshotT(t, stdOut.String())
		})
	}
}
