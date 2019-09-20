package cmd

import (
	formatter "github.com/kwix/logrus-module-formatter"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var log = logrus.WithField("module", "cmd")

var (
	verbose bool
	logging string

	RootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  "please use a command",
		Run: func(cmd *cobra.Command, args []string) {
			// fall back on default help if no command is passed
			cmd.HelpFunc()(cmd, args)
		},
	}
)

func initLogging() {
	if verbose && logging == "" {
		logging = "*=debug"
	}

	if logging == "" {
		logging = "*=info"
	}

	f, err := formatter.New(formatter.NewModulesMap(logging))
	if err != nil {
		panic(err)
	}

	logrus.SetFormatter(f)

	log.Debug("Debug mode")
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Display debug messages")
	RootCmd.PersistentFlags().StringVar(&logging, "logging", "", "Display debug messages")

	cobra.OnInitialize(initLogging)
}
