package option

import (
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

type PersistCmdConfig struct {
	Verbose bool
}

func NewPersistCmdConfigFromViper() (*PersistCmdConfig, error) {
	var conf PersistCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, xerrors.Errorf("failed to unmarshal config from viper: %w", err)
	}
	return &conf, nil
}
