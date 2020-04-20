package files

import (
	"github.com/bradleyjkemp/cupaloy"
	"testing"
)

func TestReader_ReadDir(t *testing.T) {
	reader := Reader{}

	paths, err := reader.ReadDir("./testData")
	if err != nil {
		t.Errorf("Read Dir failed: %s", err)
	}

	cupaloy.SnapshotT(t, paths)
}

func TestReader_IsFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			path:     "./testData",
			expected: false,
		},
		{
			path:     "./testData/fileOne.txt",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := Reader{}

			result, err := reader.IsFile(test.path)
			if err != nil {
				t.Errorf("test: %s, failed: %s", test.name, err)
			}

			if result != test.expected {
				t.Errorf("It read a directory as a file")
			}

		})
	}
}
