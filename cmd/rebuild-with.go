package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/rakyll/statik/fs"

	"github.com/spf13/afero"

	_ "github.com/mpppk/everest/statik"
	"github.com/spf13/cobra"
)

func newRebuildWithCmd(_fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "rebuild-with",
		Short: "rebuild everest",
		//Args:  cobra.ExactArgs(1),
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			statikFs, err := fs.New()
			if err != nil {
				fmt.Println("err1")
				return err
			}
			dstPath := os.TempDir()
			if err := writeFs(statikFs, dstPath); err != nil {
				return err
			}
			mainPath := filepath.Join(dstPath, "main.go")
			exeName := path.Base(os.Args[0])
			if err := exec.Command("go", "build", "-o", exeName, mainPath).Run(); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd, nil
}

func writeFs(fileSystem http.FileSystem, dst string) error {
	return fs.Walk(fileSystem, "/", func(path string, info os.FileInfo, err error) error {
		dstPath := filepath.Join(dst, path)
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
