package monitor

import (
	"context"
	"time"

	"github.com/tzapu/deposits-monitor/helper"

	"github.com/corpetty/go-alethio-api/alethio"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "monitor")

func Run(client *alethio.Client, address string) {
	// init
	ctx := context.Background()
	pollURL := ""

	// check if we need to pre-fill
	if true {
		// TODO remove when we have pre-fill
		transfers, err := client.Account.GetEtherTransfers(ctx, address)
		helper.FatalIfError(err, "get transfers")
		pollURL = transfers.Links.Prev

		if len(transfers.Data) > 0 {
			for {
				transfers, err := client.EtherTransfers.GetNext(ctx, transfers)
				helper.FatalIfError(err, "traverse transfers")
				log.Print(len(transfers.Data))
				if len(transfers.Data) == 0 {
					// no more transfers
					break
				}
			}
		}
	}

	// monitor for new stuff
	for {
		transfers, err := client.EtherTransfers.Get(ctx, pollURL)
		helper.FatalIfError(err, "poll for transfers")

		spew.Dump(transfers)
		time.Sleep(time.Second * 15)

	}

}
