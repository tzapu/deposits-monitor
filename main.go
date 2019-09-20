package main

import (
	"fmt"
	"os"

	"github.com/tzapu/deposits-monitor/cmd"
)

var (
	buildVersion string
)

func main() {
	cmd.RootCmd.Version = buildVersion

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
