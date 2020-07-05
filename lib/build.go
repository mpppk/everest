package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	binName = "bin"
)

var DefaultAppConfig = &AppConfig{
	AppName:  "everest",
	IconPath: "./defaultembedded/src/everest.icns",
	Width:    720,
	Height:   480,
	MacOS: &MacOSAppConfig{
		Identifier: "com.github.mpppk.everest",
	},
}

func BuildMacOsApp(config *AppConfig, execPath, dstDir string) (string, error) {
	appName := config.AppName
	if !strings.Contains(appName, ".app") {
		appName += ".app"
	}
	appPath := path.Join(dstDir, appName)

	if IsExist(appPath) {
		if err := os.RemoveAll(appPath); err != nil {
			return "", fmt.Errorf("failed to remove app from %s: %w", appPath, err)
		}
	}

	resourcePath := path.Join(appPath, "Contents", "Resources")
	if err := os.MkdirAll(resourcePath, 0777); err != nil {
		return "", fmt.Errorf("failed to create app resource directory to %s: %w", appPath, err)
	}

	binaryDirPath := path.Join(appPath, "Contents", "MacOS")
	if err := os.MkdirAll(binaryDirPath, 0777); err != nil {
		return "", fmt.Errorf("failed to create app MacOS directory to %s: %w", appPath, err)
	}

	binaryPath := path.Join(binaryDirPath, binName)
	if err := copyFile(execPath, binaryPath); err != nil {
		return "", fmt.Errorf("failed to copy file from %s to %s: %w", execPath, binaryPath, err)
	}

	if err := os.Chmod(binaryPath, 0755); err != nil {
		return "", fmt.Errorf("failed to change permission of binary: %w", err)
	}

	iconName := path.Base(config.IconPath)
	iconPath := path.Join(resourcePath, iconName)

	infoPlistPath := path.Join(appPath, "Contents", "Info.plist")
	infoPlist := generateInfoPlist(binName, iconName, config.MacOS.Identifier)
	if err := ioutil.WriteFile(infoPlistPath, []byte(infoPlist), 0777); err != nil {
		return "", fmt.Errorf("failed to write info.plist to %s: %w", infoPlistPath, err)
	}

	if err := copyFile(config.IconPath, iconPath); err != nil {
		return "", fmt.Errorf("failed to copy file from %s to %s: %w", config.IconPath, iconPath, err)
	}
	return appPath, nil
}

func generateInfoPlist(executablePath, iconPath, identifier string) string {
	return fmt.Sprintf(infoPlistTemplate(), executablePath, iconPath, identifier)
}

func infoPlistTemplate() string {
	return `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>%s</string>
	<key>CFBundleIconFile</key>
	<string>%s</string>
	<key>CFBundleIdentifier</key>
	<string>%s</string>
</dict>
</plist>
`
}
