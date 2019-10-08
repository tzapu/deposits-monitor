package importer

import (
	"context"
	"time"

	"github.com/tzapu/deposits-monitor/data"

	"github.com/tzapu/deposits-monitor/helper"

	"github.com/corpetty/go-alethio-api/alethio"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "monitor")

type Importer struct {
	address string
	api     *alethio.Client
	data    *data.Data
}

func (imp *Importer) Run() {
	// init
	ctx := context.Background()

	// check if we need to pre-fill
	synced := imp.Synced()

	// check if we need to pre-fill
	pollURL := imp.PollURL()
	if !synced {
		// start pre-filling db
		pollURL = imp.Backfill()
		log.Infof("backfill done")
	}

	// monitor for new stuff
	for {
		transfers, err := imp.api.EtherTransfers.Get(ctx, pollURL)
		helper.FatalIfError(err, "poll for transfers")

		//spew.Dump(transfers)
		_ = transfers
		time.Sleep(time.Second * 15)
	}

}

// Backfill the database
func (imp *Importer) Backfill() string {
	ctx := context.Background()
	pollURL := ""

	// is there a scrape in progress?
	scrapeURL := imp.ScrapeURL()

	// if we don't have a ScrapeURL then  this is the first run
	if scrapeURL == "" {
		// if we had no scrape url, then it's our first run
		// get initial transfers
		transfers, err := imp.api.Account.GetEtherTransfers(ctx, imp.address)
		helper.FatalIfError(err, "get transfers")

		_ = imp.processTransfers(transfers)
		helper.FatalIfError(err, "process transfers")

		// extract future poll urls from initial transfers
		pollURL = transfers.Links.Prev
		imp.SetPollURL(pollURL)

		// update scrape url  so we know where to start from if this fails
		scrapeURL = transfers.Links.Next
		imp.SetScrapedURL(scrapeURL)
	}

	for {
		transfers, err := imp.api.EtherTransfers.Get(ctx, scrapeURL)
		helper.FatalIfError(err, "traverse transfers", scrapeURL)
		//spew.Dump(transfers)
		done := imp.processTransfers(transfers)
		if done {
			break
		}

		// update scrape url for next page
		scrapeURL = transfers.Links.Next
		imp.SetScrapedURL(scrapeURL)
	}
	// mark as synced
	imp.SetSynced(true)

	return pollURL
}

func (imp *Importer) processTransfers(transfers *alethio.EtherTransfers) bool {
	if len(transfers.Data) == 0 {
		return true
	}

	return false
}

func New(address string, api *alethio.Client, data *data.Data) *Importer {
	return &Importer{
		address: address,
		api:     api,
		data:    data,
	}
}
