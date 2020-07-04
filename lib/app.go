package lib

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type MacOSAppConfig struct {
	Identifier string
}

type AppConfig struct {
	AppName  string
	IconPath string
	Width    int
	Height   int
	MacOS    *MacOSAppConfig
}

func ParseAppConfig(configPath string) (*AppConfig, error) {
	var config AppConfig
	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load everest config file from %s: %w", configPath, err)
	}
	if err := yaml.Unmarshal(contents, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config YAML from %s: %w", configPath, err)
	}
	return &config, nil
}

func ApplyDefaultToAppConfig(base, defaultConf *AppConfig) {
	if base.AppName == "" {
		base.AppName = defaultConf.AppName
	}
	if base.IconPath == "" {
		base.IconPath = defaultConf.IconPath
	}
	if base.Width == 0 {
		base.Width = defaultConf.Width
	}
	if base.Height == 0 {
		base.Height = defaultConf.Height
	}

	if base.MacOS == nil {
		base.MacOS = &MacOSAppConfig{}
	}

	if base.MacOS.Identifier == "" {
		base.MacOS.Identifier = defaultConf.MacOS.Identifier
	}
}
