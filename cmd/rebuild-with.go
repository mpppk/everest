package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mpppk/everest/lib"

	"github.com/mpppk/everest/self"

	"github.com/mpppk/everest/command"

	"github.com/rakyll/statik/fs"

	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newRebuildWithCmd(_fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "rebuild-with",
		Short: "rebuild everest",
		Args:  cobra.ExactArgs(1),
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			embeddedPath := args[0]
			dstPath := os.TempDir()

			if err := lib.GenerateEmbeddedPackage(embeddedPath, dstPath); err != nil {
				return err
			}

			if err := writeFs(self.Self, dstPath); err != nil {
				return err
			}

			mainPath := filepath.Join(dstPath, "main.go")
			exePath, err := os.Executable()
			if err != nil {
				return err
			}

			stdout, err := command.GoBuild(&command.BuildOption{
				OutputPath: exePath,
				BuildPath:  mainPath,
			})
			if err != nil {
				return err
			}
			if stdout != "" {
				cmd.Println(stdout)
			}
			return nil
		},
	}
	return cmd, nil
}

func writeFs(fileSystem http.FileSystem, dst string) error {
	return fs.Walk(fileSystem, "/", func(path string, info os.FileInfo, err error) error {
		dstPath := filepath.Join(dst, path)
		fmt.Println(dstPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0777)
		}
		file, err := fileSystem.Open(path)
		if err != nil {
			return err
		}
		// FIXME use io.Pipe
		contents, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(dstPath, contents, 0777)
	})
}

func init() {
	cmdGenerators = append(cmdGenerators, newRebuildWithCmd)
}
