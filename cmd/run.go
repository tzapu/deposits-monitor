package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tzapu/deposits-monitor/monitor"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "monitor for deposits",
	Long:  "monitor for deposits",
	Run: func(cmd *cobra.Command, args []string) {
		monitor.Run()
	},
}

func init() {
	RootCmd.AddCommand(monitorCmd)
}
