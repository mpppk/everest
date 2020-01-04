package option

// ServeCmdConfig is config for sum command
type ServeCmdConfig struct {
	Port string
}

// NewServeCmdConfigFromViper generate config for sum command from viper
func NewServeCmdConfigFromViper() (*ServeCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newServeCmdConfigFromRawConfig(rawConfig), err
}

func newServeCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *ServeCmdConfig {
	return &ServeCmdConfig{
		Port: rawConfig.Port,
	}
}

func (c *ServeCmdConfig) validate() error {
	return nil
}
