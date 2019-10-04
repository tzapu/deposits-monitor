package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/tzapu/deposits-monitor/data"

	"github.com/corpetty/go-alethio-api/alethio"
	"github.com/spf13/cobra"
	"github.com/tzapu/deposits-monitor/monitor"
	"github.com/tzapu/deposits-monitor/server"
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

		// make sure we catch all defers
		var err error
		defer func() {
			if err != nil {
				log.Fatalf("main: %s", err)
			}
		}()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)

		go func() {
			<-signals
			log.Println("[signals] We're in here!")
			os.Exit(0)
		}()

		dbFile := fmt.Sprintf("db/%s.bolt", address)
		data, err := data.New(dbFile)
		if err != nil {
			return
		}
		defer func() {
			err := data.Close()
			if err != nil {
				log.Errorf("failed to close db: %s", err)
				return
			}
			log.Infof("db closed")
		}()

		client, err := alethio.NewClient(
			alethio.Opts.URL(apiEndpoint),
			alethio.Opts.APIKey(apiKey),
		)
		if err != nil {
			return
		}

		monitor.Run(client, address)
		server.Serve()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}
