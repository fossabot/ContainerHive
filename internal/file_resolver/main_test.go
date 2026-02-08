package file_resolver

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetFileCandidates(t *testing.T) {
	tests := map[string]struct {
		baseName   string
		extensions []string
		expected   []string
	}{
		"without extensions with Dockerfile": {
			baseName:   "Dockerfile",
			extensions: nil,
			expected:   []string{"Dockerfile.gotpl"},
		},
		"with yaml and yml extension": {
			baseName:   "test",
			extensions: []string{"yaml", "yml"},
			expected:   []string{"test.yaml.gotpl", "test.yml.gotpl"},
		},
		"with only yaml extension": {
			baseName:   "config",
			extensions: []string{"yaml"},
			expected:   []string{"config.yaml.gotpl"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetFileCandidates(tc.baseName, tc.extensions...)

			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("getFileCandidates() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}
