package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/mpppk/everest/internal/option"

	"github.com/mpppk/everest/lib"
	"github.com/mpppk/everest/self"

	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

const cmdPkgPath = "github.com/mpppk/everest/cmd"
const executableName = "bin"

func newRebuildWithCmd(_fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "rebuild-with",
		Short: "rebuild everest with specified resources",
		Args:  cobra.ExactArgs(1),
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRebuildWithCmdConfigFromViper()
			if err != nil {
				return err
			}
			embeddedPath := args[0]
			configPath, hasConfig := getConfigFilePath(embeddedPath)

			var appConfig *lib.AppConfig
			if hasConfig {
				a, err := lib.ParseAppConfig(configPath)
				if err != nil {
					return fmt.Errorf("failed to parse everest appConfig from %s: %w", configPath, err)
				}
				appConfig = a
				lib.ApplyDefaultToAppConfig(appConfig, lib.DefaultAppConfig)
			} else if conf.App {
				appConfig = lib.DefaultAppConfig
			}

			dstPath, err := ioutil.TempDir(".", "everest-rebuild")
			if err != nil {
				return err
			}

			if buildLog, err := rebuildEverest(embeddedPath, dstPath, executableName, appConfig); err != nil {
				if err := lib.RemoveContents(dstPath); err != nil {
					return err
				}
				return err
			} else {
				cmd.Println(buildLog)
			}

			execPath := path.Join(dstPath, executableName)
			if appConfig != nil {
				switch runtime.GOOS {
				case "darwin":
					if macOsAppPath, err := lib.BuildMacOsApp(appConfig, execPath, "."); err != nil {
						return fmt.Errorf("failed to build MacOSApp: %w", err)
					} else {
						cmd.Printf("MacOS App is generated to %s\n", macOsAppPath)
					}
				default:
					cmd.Println("unknown OS:", runtime.GOOS)
				}
			} else {
				everestPath, err := os.Executable()
				if err != nil {
					return fmt.Errorf("failed to detect everest path: %w", err)
				}
				if err := os.Rename(execPath, everestPath); err != nil {
					return fmt.Errorf("failed to move new executable from %s to %s: %w", execPath, everestPath, err)
				}
			}

			cmd.Println("removing", dstPath)
			if err := lib.RemoveContents(dstPath); err != nil {
				return err
			}

			return nil
		},
	}

	appFlag := &option.BoolFlag{
		Flag: &option.Flag{
			Name:  "app",
			Usage: "enable app mode",
		},
		Value: false,
	}

	if err := option.RegisterBoolFlag(cmd, appFlag); err != nil {
		return nil, err
	}

	return cmd, nil
}

func getConfigFilePath(embeddedPath string) (string, bool) {
	configFilePath := path.Join(embeddedPath, "everest.yaml")
	if lib.IsExist(configFilePath) {
		return configFilePath, true
	}

	configFilePath = path.Join(embeddedPath, "everest.yml")
	if lib.IsExist(configFilePath) {
		return configFilePath, true
	}
	return "", false
}

func rebuildEverest(embeddedPath, workDir, execPath string, appConfig *lib.AppConfig) (buildLog string, err error) {
	if err := lib.GenerateEmbeddedPackage(embeddedPath, workDir); err != nil {
		return "", fmt.Errorf("failed to generate embedded package: %w", err)
	}

	if err := lib.WriteFs(self.Self, workDir); err != nil {
		return "", fmt.Errorf("failed to write self fs: %w", err)
	}

	if err := lib.GenerateSelfPackage(workDir, workDir); err != nil {
		return "", fmt.Errorf("failed to generate self package: %w", err)
	}

	buildOption := &lib.BuildOption{
		Option: lib.Option{
			Dir: workDir,
		},
		OutputPath: execPath,
		BuildPath:  ".",
	}
	if appConfig != nil {
		buildOption.LdFlags = append(buildOption.LdFlags,
			fmt.Sprintf("-X %s.appMode=true", cmdPkgPath),
			fmt.Sprintf("-X %s.width=%d", cmdPkgPath, appConfig.Width),
			fmt.Sprintf("-X %s.height=%d", cmdPkgPath, appConfig.Height),
		)
	}

	return lib.GoBuild(buildOption)
}

func init() {
	cmdGenerators = append(cmdGenerators, newRebuildWithCmd)
}
