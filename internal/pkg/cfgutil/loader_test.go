package cfgutil_test

import (
	"fmt"
	"testing"

	"github.com/raphoester/space-trouble-api/internal/pkg/cfgutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ConfigExample struct {
	Key    string
	Nested struct {
		Key string
	}
}

func TestLoader(t *testing.T) {
	t.Run("should populate struct with yaml config", func(t *testing.T) {
		cfg, err := loadEnv(t, loadEnvParams{
			configPath: "/config/config.yaml",
			files: map[string]string{
				"/config/config.yaml": `---
key: value`,
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "value", cfg.Key)
	})

	t.Run("should populate content from env on top of the yaml file", func(t *testing.T) {
		cfg, err := loadEnv(t, loadEnvParams{
			configPath: "/config/config.yaml",
			env: map[string]string{
				"KEY": "fromEnv",
			},
			files: map[string]string{
				"/config/config.yaml": `---
key: value`,
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "fromEnv", cfg.Key)
	})

	t.Run("should find file with name that is not config.yaml", func(t *testing.T) {
		cfg, err := loadEnv(t, loadEnvParams{
			configPath: "/config/my-config.yaml",
			files: map[string]string{
				"/config/my-config.yaml": `---
key: value
nested:
  key: nestedValue`,
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "value", cfg.Key)
		assert.Equal(t, "nestedValue", cfg.Nested.Key)
	})
}

type loadEnvParams struct {
	env        map[string]string
	files      map[string]string
	configPath string
}

func loadEnv(t *testing.T, p loadEnvParams) (*ConfigExample, error) {
	loader := cfgutil.NewLoader(p.configPath)
	fs := afero.NewMemMapFs()

	for file, content := range p.files {
		err := afero.WriteFile(fs, file, []byte(content), 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}
	}

	loader.WithFS(fs)
	for key, value := range p.env {
		t.Setenv(key, value)
	}

	err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	rcv := ConfigExample{}
	err = loader.Unmarshal(&rcv)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &rcv, nil
}
