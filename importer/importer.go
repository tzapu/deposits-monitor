package importer

import (
	"context"
	"fmt"
	"time"

	"github.com/tzapu/deposits-monitor/helper"

	"github.com/tzapu/deposits-monitor/data"

	"github.com/corpetty/go-alethio-api/alethio"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "monitor")

type Importer struct {
	api  *alethio.Client
	data *data.Data
}

func (imp *Importer) Run() {
	// let all else spew messages first
	time.Sleep(time.Second)

	markets := map[string]string{
		//"0x6c8c6b02e7b2be14d4fa6022dfd6d75921d90e4e": "cBAT",
		"0xf5dce57282a584d2746faf1593d3121fcac444dc": "cDAI",
		//"0x158079ee67fce2f58472a96584a73c7ab9ac95c1": "cREP",
		//"0xb3319f5d18bc0d84dd1b4825dcde5d5f7266d407": "cZRX",
		//"0x39aa39c021dfbae8fac545936693ac917d5e7563": "cUSDC",
		//"0xc11b1268c1a384e55c48c2391d8d480264a3a7f4": "cWBTC",
		//"0x4ddc2d193948926d02f9b1fe9e1daa0718270ed5": "cETH",
	}
	events := map[string]string{
		// "0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f": "Mint",
		// "0xe5b754fb1abb7f01b499791d0b820ae3b6af3424ac1c59768edb53f4ec31a929": "Redeem",
		// "0x13ed6866d4e1ee6da46f845c46d7e54120883d75c5ea9a2dacc1c4ca8984ab80": "Borrow",
		// "0x298637f684da70674f26509b10f07ec2fbc77a335ab1e7d6215a4b2484d8bb52": "LiquidateBorrow",
		"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": "Transfer",
	}

	for addr, market := range markets {
		for topic, event := range events {
			go imp.monitorEvents(addr, market, topic, event)
		}
	}

	select {}
}

func (imp *Importer) monitorEvents(address, market, topic, event string) {
	// init
	ctx := context.Background()

	// monitor for new stuff
	log.Infof("%s %s: monitoring for new events", market, event)
	pollURL := imp.PollURL(market, event)
	if pollURL == "" {
		pollURL = initialLogEntriesLink(address, topic)
		log.Infof("%s %s: intial scrape from %s", market, event, pollURL)
	}

	for {
		events, err := imp.api.LogEntries.Get(ctx, pollURL)
		helper.FatalIfError(err, "poll for log entries")

		if len(events.Data) > 0 {
			// aggregate the data
			imp.ProcessTransfers(address, market, events.Data)

			// save progress
			pollURL = events.Links.Prev
			imp.SetPollURL(market, event, pollURL)
			log.Debugf("%s %s, set poll url to %s", market, event, pollURL)
		} else {
			log.Debugf("%s %s empty log entries response, sleeping", market, event)
			time.Sleep(time.Second * 15)
		}

		time.Sleep(time.Second)
	}

}

/*
func (imp *Importer) processTransfers(transfers *alethio.EtherTransfers) bool {
	if len(transfers.Data) == 0 {
		return true
	}

	// TODO send just data, once it's a separate struct in the API
	imp.SaveTransfers(transfers)

	log.Debugf("processed %d records", len(transfers.Data))

	return false
}
*/
func initialLogEntriesLink(contract, topic string) string {
	return fmt.Sprintf("https://api.aleth.io/v1/log-entries?filter[loggedBy]=%s&filter[hasLogTopics.0]=%s&page[limit]=10&page[prev]=0x00000000000000000000000000000000", contract, topic)
}

func New(api *alethio.Client, data *data.Data) *Importer {
	return &Importer{
		api:  api,
		data: data,
	}
}
