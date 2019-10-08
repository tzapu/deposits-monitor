package importer

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/boltdb/bolt"
	"github.com/corpetty/go-alethio-api/alethio"
	"github.com/tzapu/deposits-monitor/helper"
)

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

func (imp *Importer) TransfersCount() int {
	cnt, err := imp.data.Count(TransfersBucket)
	helper.FatalIfError(err, "get stats for transfers bucket")

	return cnt
}

func (imp *Importer) SaveTransfers(transfers *alethio.EtherTransfers) {
	err := imp.data.DB.Update(func(tx *bolt.Tx) error {
		transfersBucket := tx.Bucket([]byte(TransfersBucket))
		dailyBucket := tx.Bucket([]byte(DailyBucket))
		daily := make(map[string]big.Int)

		// save transfers and aggregate values
		for _, t := range transfers.Data {
			date := time.Unix(int64(t.Attributes.BlockCreationTime), 0)
			key := DateToByte(date)
			day := DateToDay(date)
			data, err := json.Marshal(t)
			if err != nil {
				return err
			}
			err = transfersBucket.Put(key, data)
			if err != nil {
				return err
			}

			total := daily[day]
			value := new(big.Int)
			value.SetString(t.Attributes.Value, 10)
			total.Add(value, &total)
			daily[day] = total
		}

		for k, v := range daily {
			key := []byte(k)
			t := dailyBucket.Get(key)
			total := new(big.Int)
			total.SetString(string(t), 10)
			total.Add(total, &v)
			err := dailyBucket.Put(key, []byte(total.String()))
			if err != nil {
				return err
			}
		}
		return nil
	})

	helper.FatalIfError(err, "save transfers")
}

func TimestampToRFC3339(timestamp int) string {
	dt := time.Unix(int64(timestamp), 0)
	return dt.UTC().Format(time.RFC3339)
}

func DateToByte(date time.Time) []byte {
	return []byte(date.UTC().Format(time.RFC3339))
}

func DateToDay(date time.Time) string {
	return date.UTC().Format("2006-01-02")
}
