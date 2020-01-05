package cmd

import (
	"net/http"

	"github.com/mpppk/everest/internal/option"

	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newServeCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve static files",
		Args:  cobra.ExactArgs(1),
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewServeCmdConfigFromViper()
			if err != nil {
				return err
			}
			rootPath := args[0]
			cmd.Println("Files are served on http://localhost:" + conf.Port)
			http.Handle("/", http.FileServer(http.Dir(rootPath)))
			return http.ListenAndServe(":"+conf.Port, nil)
		},
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newServeCmd)
}
