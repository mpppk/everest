package cmd

import (
	"net/http"

	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

func newServeCmd(fs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve static files",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			http.Handle("/", http.FileServer(http.Dir(".")))
			return http.ListenAndServe(":3000", nil)
		},
	}
	return cmd, nil
}

func init() {
	cmdGenerators = append(cmdGenerators, newServeCmd)
}
