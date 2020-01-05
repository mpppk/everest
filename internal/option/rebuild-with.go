package option

type RebuildWithCmdConfig struct {
	Port string
}

func NewRebuildWithCmdConfigFromViper() (*RebuildWithCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newRebuildWithCmdConfigFromRawConfig(rawConfig), err
}

func newRebuildWithCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *RebuildWithCmdConfig {
	return &RebuildWithCmdConfig{
		Port: rawConfig.Port,
	}
}

func (c *RebuildWithCmdConfig) validate() error {
	return nil
}
