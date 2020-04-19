package secrets

import (
	"fmt"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/eikc/gapp/internal/secrets/mocks"
	"github.com/golang/mock/gomock"
	"gopkg.in/yaml.v2"
	"testing"
)

func Test_parser_Parse(t *testing.T) {
	tests := []struct {
		name      string
		params    string
		configure func(reader *mocks.MockFileReader)
	}{
		{
			name:   "path from params is a file",
			params: "testFile",
			configure: func(reader *mocks.MockFileReader) {
				testData := YamlDefinition{
					Secrets: []SecretDefinition{
						{
							Name:  "name",
							Value: "value",
							Repos: []string{"eikc/repo-one"},
						},
					},
				}
				bytes, _ := yaml.Marshal(&testData)

				reader.EXPECT().IsFile("testFile").Times(1).Return(true, nil)
				reader.EXPECT().ReadFile("testFile").Times(1).Return(bytes, nil)
			},
		},
		{
			name:   "provided path is a directory",
			params: "testFolder",
			configure: func(reader *mocks.MockFileReader) {
				testPaths := []string{
					"testFile.yaml",
					"testTwo.yaml",
				}

				fileOneContent := YamlDefinition{
					Secrets: []SecretDefinition{
						{
							Name:  "secret 1",
							Value: "123321",
							Repos: []string{"eikc/gapp"},
						},
					},
				}

				fileTwoContent := YamlDefinition{
					Secrets: []SecretDefinition{
						{
							Name:  "secret 2",
							Value: "secret 2",
							Repos: []string{"eikc/gapp"},
						},
					},
				}

				bytesOne, _ := yaml.Marshal(&fileOneContent)
				bytesTwo, _ := yaml.Marshal(&fileTwoContent)

				reader.EXPECT().IsFile("testFolder").Times(1).Return(false, nil)
				reader.EXPECT().ReadDir("testFolder").Times(1).Return(testPaths, nil)

				reader.EXPECT().ReadFile("testFile.yaml").Times(1).Return(bytesOne, nil)
				reader.EXPECT().ReadFile("testTwo.yaml").Times(1).Return(bytesTwo, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockFileReader := mocks.NewMockFileReader(ctrl)

			tt.configure(mockFileReader)

			p := parser{
				reader: mockFileReader,
			}

			got, err := p.Parse(tt.params)

			cupaloy.SnapshotT(t, got, err)
		})
	}
}

func Test_parser_parseSecret(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		configure func(reader *mocks.MockFileReader)
	}{
		{
			name: "Parse secrets will parse the readers content correctly",
			path: "./secrets.yaml",
			configure: func(reader *mocks.MockFileReader) {
				returned := YamlDefinition{
					Secrets: []SecretDefinition{
						{
							Name:  "github token",
							Value: "123332121312312312312",
							Repos: []string{"eikc/gapp"},
						},
					},
				}

				bytes, err := yaml.Marshal(&returned)
				if err != nil {
					t.Errorf("Configuring test failed: %s", err)
				}

				reader.EXPECT().ReadFile("./secrets.yaml").Times(1).Return(bytes, nil)
			},
		},
		{
			name: "Reader returns an error",
			path: "./secrets.yaml",
			configure: func(reader *mocks.MockFileReader) {
				reader.EXPECT().ReadFile("./secrets.yaml").Return(nil, fmt.Errorf("a fake error from io.ReadFile"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockReader := mocks.NewMockFileReader(ctrl)
			tt.configure(mockReader)

			p := parser{
				reader: mockReader,
			}

			got, err := p.parseSecret(tt.path)

			cupaloy.SnapshotT(t, got, err)
		})
	}
}

func Test_parser_sortSecrets(t *testing.T) {
	secrets := YamlDefinition{
		Secrets: []SecretDefinition{
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

	p := parser{}

	sorted := p.sortSecrets(secrets)

	cupaloy.SnapshotT(t, sorted)
}
