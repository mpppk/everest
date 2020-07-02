package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mpppk/everest/internal/option"

	"github.com/mpppk/everest/lib"

	"github.com/mpppk/everest/self"

	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

const cmdPkgPath = "github.com/mpppk/everest/cmd"

func newRebuildWithCmd(_fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "rebuild-with",
		Short: "rebuild everest",
		Args:  cobra.ExactArgs(1),
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRebuildWithCmdConfigFromViper()
			if err != nil {
				return err
			}
			embeddedPath := args[0]
			dstPath, err := ioutil.TempDir(".", "everest-rebuild")
			if err != nil {
				return err
			}

			if err := rebuild(cmd, embeddedPath, dstPath, conf.App); err != nil {
				if err := lib.RemoveContents(dstPath); err != nil {
					return err
				}
				return err
			}
			fmt.Println("removing", dstPath)
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

func rebuild(cmd *cobra.Command, embeddedPath, dstPath string, appMode bool) error {
	if err := lib.GenerateEmbeddedPackage(embeddedPath, dstPath); err != nil {
		return fmt.Errorf("failed to generate embedded package: %w", err)
	}

	if err := lib.WriteFs(self.Self, dstPath); err != nil {
		return fmt.Errorf("failed to write self fs: %w", err)
	}

	if err := lib.GenerateSelfPackage(dstPath, dstPath); err != nil {
		return fmt.Errorf("failed to generate self package: %w", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	buildOption := &lib.BuildOption{
		Option: lib.Option{
			Dir: dstPath,
		},
		OutputPath: exePath,
		BuildPath:  ".",
	}
	if appMode {
		buildOption.LdFlags = append(buildOption.LdFlags,
			fmt.Sprintf("-X %s.appMode=true", cmdPkgPath),
		)
	}

	stdout, err := lib.GoBuild(buildOption)

	if err != nil {
		return err
	}
	if stdout != "" {
		cmd.Println(stdout)
	}
	return nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newRebuildWithCmd)
}
