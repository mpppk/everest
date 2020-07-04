package option

type RebuildWithCmdConfig struct {
	App bool
}

func NewRebuildWithCmdConfigFromViper() (*RebuildWithCmdConfig, error) {
	rawConfig, err := newCmdRawConfig()
	return newRebuildWithCmdConfigFromRawConfig(rawConfig), err
}

func newRebuildWithCmdConfigFromRawConfig(rawConfig *CmdRawConfig) *RebuildWithCmdConfig {
	return &RebuildWithCmdConfig{
		App: rawConfig.App,
	}
}

func (c *RebuildWithCmdConfig) validate() error {
	return nil
}
