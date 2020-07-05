package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mpppk/everest/lib"

	"github.com/mpppk/everest/embedded"

	"github.com/mpppk/everest/internal/option"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// use string for "-ldflags -X"
var appMode = "false"
var width, height = "720", "480"

func NewRootCmd(aferoFs afero.Fs) (*cobra.Command, error) {
	pPreRunE := func(cmd *cobra.Command, args []string) error {
		conf, err := option.NewRootCmdConfigFromViper()
		if err != nil {
			return err
		}
		lib.InitializeLog(conf.Verbose)
		return nil
	}
	cmd := &cobra.Command{
		Use:               "everest",
		Short:             "everest",
		Args:              cobra.MaximumNArgs(1),
		PersistentPreRunE: pPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRootCmdConfigFromViper()
			if err != nil {
				return err
			}

			s := lib.New(conf.Port)

			if len(args) == 0 {
				cmd.Println("Embedded files are served on http://localhost:" + conf.Port)
				if err := s.AddFileHandlerFromFs(embedded.Embedded); err != nil {
					return err
				}
			} else {
				rootPath := args[0]
				cmd.Println("Files are served on http://localhost:" + conf.Port)
				if err := s.AddFileHandlerFromPath(rootPath); err != nil {
					return err
				}
			}
			if err := s.AddApiHandler(); err != nil {
				return err
			}

			if appMode == "true" {
				w, h, err := parseWidthAndHeight(width, height)
				if err != nil {
					return err
				}
				return s.StartWithApp(w, h)
			} else {
				return s.Start()
			}
		},
	}

	newPortFlag := func() *option.StringFlag {
		return &option.StringFlag{
			Flag: &option.Flag{
				Name:         "port",
				Usage:        "port",
				IsPersistent: false,
			},
			Value: "3000",
		}
	}
	if err := option.RegisterStringFlag(cmd, newPortFlag()); err != nil {
		return nil, err
	}

	verboseFlag := &option.BoolFlag{
		Flag: &option.Flag{
			Name:         "verbose",
			Usage:        "show details logs",
			IsPersistent: true,
		},
		Value: false,
	}
	if err := option.RegisterBoolFlag(cmd, verboseFlag); err != nil {
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

func parseWidthAndHeight(width, height string) (int, int, error) {
	w, err := strconv.Atoi(width)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse width from %s: %w", width, err)
	}

	h, err := strconv.Atoi(height)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse height from %s: %w", height, err)
	}
	return w, h, nil
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
