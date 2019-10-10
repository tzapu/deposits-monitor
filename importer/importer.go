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
	if !synced {
		// start pre-filling db
		imp.Backfill()
		// mark as synced
		imp.SetSynced(true)

		cnt := imp.TransfersCount()

		log.
			WithField("transfers count", cnt).
			Infof("backfill done")
	}
	pollURL := imp.PollURL()

	// monitor for new stuff
	log.Infof("monitoring for new transfers")
	for {
		transfers, err := imp.api.EtherTransfers.Get(ctx, pollURL)
		helper.FatalIfError(err, "poll for transfers")

		imp.processTransfers(transfers)

		pollURL = transfers.Links.Prev
		imp.SetPollURL(pollURL)
		log.Debugf("set poll url to %s", pollURL)

		time.Sleep(time.Second * 15)
	}

}

// Backfill the database
func (imp *Importer) Backfill() {
	ctx := context.Background()

	// is there a scrape in progress?
	scrapeURL := imp.ScrapeURL()

	// if we don't have a ScrapeURL then  this is the first run
	if scrapeURL == "" {
		log.Debugf("starting backfill")
		// if we had no scrape url, then it's our first run
		// get initial transfers
		transfers, err := imp.api.Account.GetEtherTransfers(ctx, imp.address)
		helper.FatalIfError(err, "get transfers")

		_ = imp.processTransfers(transfers)
		helper.FatalIfError(err, "process transfers")

		// extract future poll urls from initial transfers
		imp.SetPollURL(transfers.Links.Prev)

		// update scrape url  so we know where to start from if this fails
		scrapeURL = transfers.Links.Next
		imp.SetScrapedURL(scrapeURL)
	} else {
		log.Debugf("continuing backfill")
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
}

func (imp *Importer) processTransfers(transfers *alethio.EtherTransfers) bool {
	if len(transfers.Data) == 0 {
		return true
	}

	// TODO send just data, once it's a separate struct in the API
	imp.SaveTransfers(transfers)

	log.Debugf("processed %d records", len(transfers.Data))

	return false
}

func New(address string, api *alethio.Client, data *data.Data) *Importer {
	return &Importer{
		address: address,
		api:     api,
		data:    data,
	}
}
