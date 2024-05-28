package pipeline

import (
	"reflect"
	"testing"
)

// Тест для функции parsePipeline
func TestParsePipeline(t *testing.T) {
	tests := []struct {
		name     string
		yamlData []byte
		expected Pipeline
		wantErr  bool
	}{
		{
			name: "Valid YAML",
			yamlData: []byte(`
pipeline:
  steps:
  - name: Build
    image: docker
    branch: master
    commands:
      - "docker build . -t myapp"
      - "docker push myapp"
`),
			expected: Pipeline{
				Steps: []Step{
					{
						Name:   "Build",
						Image:  "docker",
						Branch: "master",
						Commands: []string{
							"docker build . -t myapp",
							"docker push myapp",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "Invalid YAML",
			yamlData: []byte("invalid YAML data"),
			expected: Pipeline{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePipeline(tt.yamlData)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parsePipeline() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
