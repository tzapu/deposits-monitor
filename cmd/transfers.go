package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/tzapu/deposits-monitor/server"

	"github.com/corpetty/go-alethio-api/alethio"
	"github.com/spf13/cobra"
	"github.com/tzapu/deposits-monitor/data"
	"github.com/tzapu/deposits-monitor/helper"
	"github.com/tzapu/deposits-monitor/importer"
)

var runCmd = &cobra.Command{
	Use:   "transfers",
	Short: "transfers monitor",
	Long:  "monitors ether transfers to an address and provides a visualisation interface",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting up server")
		apiEndpoint := "https://api.aleth.io/v1"
		apiKey := ""
		address := "0x0000000000000000000000000000000000000000"
		//address := "0x3378eeaf39dffb316a95f31f17910cbb21ace6bb" // eth2 goerli deposit contract

		// API client
		client, err := alethio.NewClient(
			alethio.Opts.URL(apiEndpoint),
			alethio.Opts.APIKey(apiKey),
		)
		helper.FatalIfError(err)

		// DB
		dbFile := fmt.Sprintf("db/%s.bolt", address)
		log.Infof("opening db %s", dbFile)
		data, err := data.New(dbFile, importer.Buckets)
		helper.FatalIfError(err, "db open")
		defer func() {
			err := data.Close()
			if err != nil {
				log.Fatalf("failed to close db: %s", err)
			}
			log.Infof("db closed")
		}()

		// catch interrupts
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Kill, os.Interrupt)

		// importer service
		imp := importer.New(address, client, data)

		// Work
		go imp.Run()
		go server.Serve()

		// wait on interrupt
		<-interrupt
		log.Info("got interrupt")
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
