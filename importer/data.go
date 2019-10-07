package importer

import "github.com/tzapu/deposits-monitor/helper"

// Buckets
const (
	SettingsBucket  = "settings"
	TransfersBucket = "transfers"
	DailyBucket     = "daily"
)

var Buckets = []string{SettingsBucket, TransfersBucket, DailyBucket}

// Keys
const (
	SyncedKey    = "synced"
	ScrapeURLKey = "scrapeURL"
	PollURLKey   = "pollURL"
)

func (imp *Importer) Synced() bool {
	synced, err := imp.data.Bool(SettingsBucket, SyncedKey)
	helper.FatalIfError(err, "get poll url")

	return synced
}

func (imp *Importer) SetSynced(synced bool) {
	err := imp.data.PutBool(SettingsBucket, SyncedKey, synced)
	helper.FatalIfError(err, "set poll url")
}

func (imp *Importer) ScrapeURL() string {
	u, err := imp.data.String(SettingsBucket, ScrapeURLKey)
	helper.FatalIfError(err, "get scraped url")

	return u
}

func (imp *Importer) SetScrapedURL(u string) {
	err := imp.data.PutString(SettingsBucket, ScrapeURLKey, u)
	helper.FatalIfError(err, "set poll url")
}

func (imp *Importer) PollURL() string {
	u, err := imp.data.String(SettingsBucket, PollURLKey)
	helper.FatalIfError(err, "get scraped url")

	return u
}

func (imp *Importer) SetPollURL(u string) {
	err := imp.data.PutString(SettingsBucket, PollURLKey, u)
	helper.FatalIfError(err, "set poll url")
}
