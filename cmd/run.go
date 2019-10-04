package cmd

import (
	"fmt"

	"github.com/tzapu/deposits-monitor/data"
	"github.com/tzapu/deposits-monitor/monitor"

	"github.com/corpetty/go-alethio-api/alethio"
	"github.com/spf13/cobra"
	"github.com/tzapu/deposits-monitor/helper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run server",
	Long:  "run server and monitor for deposits",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting up server")
		apiEndpoint := "https://api.aleth.io/v1"
		apiKey := ""
		address := "0x0000000000000000000000000000000000000000"
		//address := "0x3378eeaf39dffb316a95f31f17910cbb21ace6bb" // eth2 goerli deposit contract

		dbFile := fmt.Sprintf("db/%s.bolt", address)
		log.Infof("opening db %s", dbFile)
		data, err := data.New(dbFile)
		helper.FatalIfError(err, "db open")
		defer func() {
			err := data.Close()
			if err != nil {
				log.Fatalf("failed to close db: %s", err)
			}
			log.Infof("db closed")
		}()

		client, err := alethio.NewClient(
			alethio.Opts.URL(apiEndpoint),
			alethio.Opts.APIKey(apiKey),
		)
		helper.FatalIfError(err)

		monitor.Run(client, address)
		//server.Serve()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
