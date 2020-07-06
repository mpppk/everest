package option

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

type CmdConfig struct {
	Port    string
	App     bool
	Server  bool
	Verbose bool
}

func NewRootCmdConfigFromViper() (*CmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}
	return newRootCmdConfigFromRawConfig(rawConfig), err
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

func newRootCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *CmdConfig {
	fmt.Printf("new root cmd %#v\n", rawConfig)
	return &CmdConfig{
		Port:    rawConfig.Port,
		App:     rawConfig.App,
		Verbose: rawConfig.Verbose,
		Server:  rawConfig.Server,
	}
}

type CmdRawConfig struct {
	CmdConfig `mapstructure:",squash"`
}

func (c *CmdRawConfig) validate() error {
	if c.App && c.Server {
		return errors.New("only one of --app and --server can be given")
	}
	return nil
}
