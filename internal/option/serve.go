package option

type ServeCmdConfig struct {
	Port string
}

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
