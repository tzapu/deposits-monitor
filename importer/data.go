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
	SyncedKey = "synced"
)

func (imp *Importer) Synced() bool {
	synced, err := imp.data.Bool(SettingsBucket, SyncedKey)
	helper.FatalIfError(err, "get poll url")

	return synced
}

func (imp *Importer) SetSynced(synced bool) {
	err := imp.data.PutBool(SettingsBucket, SyncedKey, synced)
	helper.FatalIfError(err, "get poll url")
}
