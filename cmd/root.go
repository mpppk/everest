package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rakyll/statik/fs"

	_ "github.com/mpppk/everest/embedded"
	"github.com/mpppk/everest/internal/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(aferoFs afero.Fs) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "everest",
		Short: "everest",
		RunE: func(cmd *cobra.Command, rags []string) error {
			conf, err := option.NewRootCmdConfigFromViper()
			if err != nil {
				return err
			}
			embeddedFs, err := fs.New()
			if err != nil {
				return err
			}
			cmd.Println("Embedded files are served on http://localhost:" + conf.Port)
			http.Handle("/", http.FileServer(embeddedFs))
			if err := http.ListenAndServe(":"+conf.Port, nil); err != nil {
				return err
			}
			return nil
		},
	}

	newPortFlag := func() *option.StringFlag {
		return &option.StringFlag{
			Flag: &option.Flag{
				Name:         "port",
				Usage:        "port",
				IsPersistent: true,
			},
			Value: "3000",
		}
	}
	if err := option.RegisterStringFlag(cmd, newPortFlag()); err != nil {
		return nil, err
	}

	var subCmds []*cobra.Command
	for _, cmdGen := range cmdGenerators {
		subCmd, err := cmdGen(aferoFs)
		if err != nil {
			return nil, err
		}
		subCmds = append(subCmds, subCmd)
	}
	cmd.AddCommand(subCmds...)

	return cmd, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd, err := NewRootCmd(afero.NewOsFs())
	if err != nil {
		panic(err)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
