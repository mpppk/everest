package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type MacOsAppConfig struct {
	AppName        string
	ExecutablePath string
	IconPath       string
	Identifier     string
}

const (
	binName = "bin"
)

// FIXME
const cmdPkgPath = "github.com/mpppk/everest/cmd"

func BuildMacOsApp(config *MacOsAppConfig, dst string) (string, error) {
	appPath := path.Join(dst, config.AppName)

	if isExist(appPath) {
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
	if err := copyFile(config.ExecutablePath, binaryPath); err != nil {
		return "", fmt.Errorf("failed to copy file from %s to %s: %w", config.ExecutablePath, binaryPath, err)
	}
	buildOpt := &BuildOption{
		Option:     Option{Dir: "."},
		OutputPath: binaryPath,
		BuildPath:  ".",
		LdFlags:    []string{fmt.Sprintf("-X %s.appMode=true", cmdPkgPath)},
	}
	if _, err := GoBuild(buildOpt); err != nil {
		return "", fmt.Errorf("failed to build go binary: %w", err)
	}

	iconName := path.Base(config.IconPath)
	iconPath := path.Join(resourcePath, iconName)

	infoPlistPath := path.Join(appPath, "Contents", "Info.plist")
	infoPlist := generateInfoPlist(binName, iconName, config.Identifier)
	if err := ioutil.WriteFile(infoPlistPath, []byte(infoPlist), 0777); err != nil {
		return "", fmt.Errorf("failed to write info.plist to %s: %w", infoPlistPath, err)
	}

	if err := copyFile(config.IconPath, iconPath); err != nil {
		return "", fmt.Errorf("failed to copy file from %s to %s: %w", config.IconPath, iconPath, err)
	}
	return appPath, nil
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
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
