package lib

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type MacOSAppConfig struct {
	Identifier  string `yaml:"Identifier"`
	IconPath    string `yaml:"IconPath"`
	AbsIconPath string
}

type AppConfig struct {
	BaseDir string
	AppName string          `yaml:"AppName"`
	Width   int             `yaml:"Width"`
	Height  int             `yaml:"Height"`
	MacOS   *MacOSAppConfig `yaml:"MacOS"`
}

func ParseAppConfig(configPath string) (*AppConfig, error) {
	var config AppConfig
	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load everest config file from %s: %w", configPath, err)
	}
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config YAML from %s: %w", configPath, err)
	}

	isDefaultMacOSIcon := true
	if config.MacOS != nil && config.MacOS.IconPath != "" {
		isDefaultMacOSIcon = false
	}

	ApplyDefaultToAppConfig(&config, DefaultAppConfig)

	config.BaseDir = filepath.Dir(configPath)

	if !isDefaultMacOSIcon && !filepath.IsAbs(config.MacOS.IconPath) {
		relPath := filepath.Join(config.BaseDir, config.MacOS.IconPath)
		absPath, err := filepath.Abs(relPath)
		if err != nil {
			return nil, fmt.Errorf("failed to convert icon path(%s) to abs path: %w", relPath, err)
		}
		config.MacOS.AbsIconPath = absPath
	}

	return &config, nil
}

func ApplyDefaultToAppConfig(base, defaultConf *AppConfig) {
	if base.AppName == "" {
		base.AppName = defaultConf.AppName
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
