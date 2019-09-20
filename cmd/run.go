package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tzapu/deposits-monitor/monitor"
	"github.com/tzapu/deposits-monitor/server"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run server",
	Long:  "run server and monitor for deposits",
	Run: func(cmd *cobra.Command, args []string) {
		monitor.Run()
		server.Serve()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
