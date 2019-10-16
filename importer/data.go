package importer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"time"

	"github.com/alethio/web3-go/ethconv"
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

func (imp *Importer) TransfersList() []Transfer {
	var transfers []Transfer
	ts, err := imp.data.Last(TransfersBucket, 12)
	helper.FatalIfError(err, "get last transfers")

	for i := range ts {
		var at APITransfer
		err = json.Unmarshal(ts[i].Value, &at)
		helper.FatalIfError(err, "unmarshal  transfer", ts[i].Key)

		// TODO make network aware
		url := "https://goerli.aleth.io/"
		switch at.Attributes.TransferType {
		case "TransactionTransfer":
			url = fmt.Sprintf("%stx/%s", url, at.Relationships.Transaction.Data.ID)
		case "ContractMessageTransfer":
			// TODO implement contract messages links
			/*
				 Data: (map[string]interface {}) (len=2) {
				    (string) (len=2) "id": (string) (len=72) "msg:0xf212d20e70d4e2c6e135f5bf392a4c346a7e3b52b4ceb3161b9564c8947b9f39:1",
				    (string) (len=4) "type": (string) (len=15) "ContractMessage"
					}
			*/
			url = fmt.Sprintf("%stx/%s", url, at.Relationships.Transaction.Data.ID)
		}

		ev, _ := ethconv.FromWei(at.Attributes.Value, ethconv.Eth, 2)
		t := Transfer{
			Hash:              at.Relationships.Transaction.Data.ID,
			BlockCreationTime: TimestampToTime(at.Attributes.BlockCreationTime),
			TransferType:      at.Attributes.TransferType,
			Value:             at.Attributes.Value,
			ETHValue:          ev,
			URL:               url,
		}
		transfers = append(transfers, t)
	}

	for i := len(transfers)/2 - 1; i >= 0; i-- {
		opp := len(transfers) - 1 - i
		transfers[i], transfers[opp] = transfers[opp], transfers[i]
	}

	return transfers
}

func (imp *Importer) DailyList() []Daily {
	var daily []Daily
	ds, err := imp.data.Last(DailyBucket, 10365)
	helper.FatalIfError(err, "get last daily")

	acc := new(big.Int)
	for i := range ds {
		date, err := time.Parse("2006-01-02", ds[i].Key)
		helper.FatalIfError(err, "parse date from key")
		value := StringToBigInt(string(ds[i].Value))
		value.Div(value, big.NewInt(1000000000000000000))
		acc.Add(acc, value)
		daily = append(daily, []int64{
			date.UTC().Unix() * 1000,
			value.Int64(),
		})
	}

	return daily
}

func StringToBigInt(s string) *big.Int {
	bi := new(big.Int)
	bi.SetString(s, 10)
	return bi
}

func TimestampToTime(timestamp int) time.Time {
	dt := time.Unix(int64(timestamp), 0)
	return dt.UTC()
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

func FormatDate(d time.Time) string {
	return d.Format(time.Stamp)
}

func FormatStart(s string) string {
	return s[:6]
}

func FormatEnd(s string) string {
	return s[len(s)-6:]
}

func FormatMiddle(s string) string {
	return s[6 : len(s)-6]
}

func FormatJSON(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

type Transfer struct {
	Hash              string
	BlockCreationTime time.Time
	TransferType      string
	Value             string
	ETHValue          string
	URL               string
}

type Daily []int64

// TODO migrate to api transfer when available, or  own internal struct
type APITransfer struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes struct {
		BlockCreationTime int    `json:"blockCreationTime"`
		Cursor            string `json:"cursor"`
		Fee               string `json:"fee"`
		GlobalRank        []int  `json:"globalRank"`
		Total             string `json:"total"`
		TransferType      string `json:"transferType"`
		Value             string `json:"value"`
	} `json:"attributes"`
	Relationships struct {
		Block struct {
			Data struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
			Links struct {
				Related string `json:"related"`
			} `json:"links"`
		} `json:"block"`
		ContractMessage struct {
			Data interface{} `json:"data"`
		} `json:"contractMessage"`
		FeeRecipient struct {
			Data struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
			Links struct {
				Related string `json:"related"`
			} `json:"links"`
		} `json:"feeRecipient"`
		From struct {
			Data interface{} `json:"data"`
		} `json:"from"`
		To struct {
			Data struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
			Links struct {
				Related string `json:"related"`
			} `json:"links"`
		} `json:"to"`
		Transaction struct {
			Data struct {
				Type string `json:"type"`
				ID   string `json:"id"`
			} `json:"data"`
		} `json:"transaction"`
	} `json:"relationships"`
}
