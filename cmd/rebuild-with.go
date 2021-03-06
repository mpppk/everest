package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"

	"github.com/mpppk/everest/internal/option"

	"github.com/mpppk/everest/lib"
	"github.com/mpppk/everest/self"

	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

const cmdPkgPath = "github.com/mpppk/everest/cmd"
const executableName = "bin"

func newRebuildWithCmd(_fs afero.Fs) (*cobra.Command, error) {
	pRunE := func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlag("app", cmd.Flags().Lookup("app")); err != nil {
			return err
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:     "rebuild-with",
		Short:   "rebuild everest with specified resources",
		Args:    cobra.ExactArgs(1),
		Long:    ``,
		PreRunE: pRunE,
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
			} else if conf.App {
				appConfig = lib.DefaultAppConfig
			}

			dstPath, err := ioutil.TempDir(".", "everest-rebuild")
			if err != nil {
				return err
			}
			log.Printf("debug: temp dir is created as %s\n", dstPath)

			defer func() {
				log.Printf("debug: removing temp dir from %s\n", dstPath)
				if err := lib.RemoveContents(dstPath); err != nil {
					log.Println("error: failed to remove temp dir from", dstPath)
					panic(err)
				}
			}()

			if buildLog, err := rebuildEverest(embeddedPath, dstPath, executableName, appConfig); err != nil {
				return fmt.Errorf("failed to rebuild everest: %w", err)
			} else if buildLog != "" {
				log.Println(buildLog)
			}

			execPath := path.Join(dstPath, executableName)

			if runtime.GOOS == "windows" {
				appDstPath := "./everest.exe"
				if appConfig != nil {
					appDstPath = fmt.Sprintf("./%s.exe",appConfig.AppName)
				}
				if err := os.Rename(execPath, appDstPath); err != nil {
					return fmt.Errorf("failed to move new executable from %s to %s: %w", execPath, appDstPath, err)
				}
				log.Printf("debug: move new binary from %s to %s\n", execPath, appDstPath)
				return nil
			}

			if appConfig != nil {
				switch runtime.GOOS {
				case "darwin":
					if macOsAppPath, err := lib.BuildMacOsApp(appConfig, execPath, "."); err != nil {
						return fmt.Errorf("failed to build MacOSApp: %w", err)
					} else {
						log.Printf("info: MacOS App is generated to %s\n", macOsAppPath)
					}
				default:
					log.Println("error: unknown OS:", runtime.GOOS)
				}
			} else {
				everestPath, err := os.Executable()
				if err != nil {
					return fmt.Errorf("failed to detect everest path: %w", err)
				}
				if err := os.Rename(execPath, everestPath); err != nil {
					return fmt.Errorf("failed to move new executable from %s to %s: %w", execPath, everestPath, err)
				}
				log.Printf("debug: move new binary from %s to %s\n", execPath, everestPath)
			}

			return nil
		},
	}

	registerFlags := func(cmd *cobra.Command) error {
		appFlag := &option.BoolFlag{
			Flag: &option.Flag{
				Name:  "app",
				Usage: "enable app mode",
			},
			Value: false,
		}

		if err := option.RegisterBoolFlag(cmd, appFlag); err != nil {
			return fmt.Errorf("failed  to register app flag: %w", err)
		}
		return nil
	}

	if err := registerFlags(cmd); err != nil {
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
	log.Printf("debug: target files (%s) are converted to Go sources and outputted to %s\n", embeddedPath, filepath.Join(workDir, "embedded"))

	if err := lib.WriteFs(self.Self, workDir); err != nil {
		return "", fmt.Errorf("failed to write self fs: %w", err)
	}
	log.Printf("debug: everest sources are outputted to %s\n", filepath.Join(workDir))

	if err := lib.GenerateSelfPackage(workDir, workDir); err != nil {
		return "", fmt.Errorf("failed to generate self package: %w", err)
	}
	log.Printf("debug: everest sources are converted to Go sources and outputted to %s\n", filepath.Join(workDir, "self"))

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

		if runtime.GOOS == "windows" {
			buildOption.LdFlags = append(buildOption.LdFlags, "-H windowsgui")
		}
	}

	return lib.GoBuild(buildOption)
}

func init() {
	cmdGenerators = append(cmdGenerators, newRebuildWithCmd)
}
