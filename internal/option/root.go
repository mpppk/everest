package option

import (
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

type CmdConfig struct {
	Port string
}

func NewRootCmdConfigFromViper() (*CmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newCmdConfigFromRawConfig(rawConfig), err
}

func newCmdRawConfig() (*CmdRawConfig, error) {
	var conf CmdRawConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, xerrors.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, xerrors.Errorf("failed to create root cmd config: %w", err)
	}
	return &conf, nil
}

func newCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *CmdConfig {
	return &CmdConfig{
		Port: rawConfig.Port,
	}
}

type CmdRawConfig struct {
	SumCmdConfig `mapstructure:",squash"`
	CmdConfig    `mapstructure:",squash"`
}

func (c *CmdRawConfig) validate() error {
	return nil
}
