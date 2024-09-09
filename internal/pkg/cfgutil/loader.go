package cfgutil

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type Loader struct {
	v *viper.Viper
}

func NewLoader(filePath string) *Loader {
	v := viper.New()
	v.SetFs(afero.NewOsFs())
	v.AddConfigPath(filepath.Dir(filePath))
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	// pass config file name without extension as config name
	b := filepath.Base(filePath)
	v.SetConfigName(strings.Replace(b, filepath.Ext(b), "", 1))
	return &Loader{v}
}

func (l *Loader) WithFS(fs afero.Fs) {
	l.v.SetFs(fs)
}

func (l *Loader) Unmarshal(cfg any) error {

	if err := l.v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	if err := l.v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
