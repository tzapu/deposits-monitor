package importer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"time"

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

func (imp *Importer) PollURL(market, event string) string {
	key := fmt.Sprintf("%s:%s:%s", PollURLKey, market, event)
	u, err := imp.data.String(SettingsBucket, key)
	helper.FatalIfError(err, "get scraped url")

	return u
}

func (imp *Importer) SetPollURL(market, event, pollURL string) {
	key := fmt.Sprintf("%s:%s:%s", PollURLKey, market, event)
	err := imp.data.PutString(SettingsBucket, key, pollURL)
	helper.FatalIfError(err, "set poll url")
}

/*
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
			/ *
				 Data: (map[string]interface {}) (len=2) {
				    (string) (len=2) "id": (string) (len=72) "msg:0xf212d20e70d4e2c6e135f5bf392a4c346a7e3b52b4ceb3161b9564c8947b9f39:1",
				    (string) (len=4) "type": (string) (len=15) "ContractMessage"
					}
			* /
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
*/
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

type Item struct {
	Address string //address

	// cToken market contract addresses
	// 0x6c8c6b02e7b2be14d4fa6022dfd6d75921d90e4e - cBAT
	// 0xf5dce57282a584d2746faf1593d3121fcac444dc - cDAI
	// 0x158079ee67fce2f58472a96584a73c7ab9ac95c1 - cREP
	// 0xb3319f5d18bc0d84dd1b4825dcde5d5f7266d407 - cZRX
	// 0x39aa39c021dfbae8fac545936693ac917d5e7563 - cUSDC
	// 0xc11b1268c1a384e55c48c2391d8d480264a3a7f4 - cWBTC
	// 0x4ddc2d193948926d02f9b1fe9e1daa0718270ed5 - cETH
	CToken string //cToken

	// "Mint (address minter, uint256 mintAmount, uint256 mintTokens):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f'
	NumberOfMints int64 //number_of_mints
	TotalMinted   int64 //total_minted

	// "Redeem (address redeemer, uint256 redeemAmount, uint256 redeemTokens):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0xe5b754fb1abb7f01b499791d0b820ae3b6af3424ac1c59768edb53f4ec31a929'"
	NumberOfRedeems int64 //number_of_redeems
	TotalRedeemed   int64 //total_redeemed
	SupplyBalance   int64 //supply_balance

	// Borrow (address borrower, uint256 borrowAmount, uint256 accountBorrows, uint256 totalBorrows):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0x13ed6866d4e1ee6da46f845c46d7e54120883d75c5ea9a2dacc1c4ca8984ab80'
	NumberOfBorrows    int64 //number_of_borrows
	TotalBorrowed      int64 //total_borrowed
	NumberOfRepayments int64 //number_of_repayments
	TotalRepaid        int64 //total_repaid

	// LiquidateBorrow (address liquidator, address borrower, uint256 repayAmount, address cTokenCollateral, uint256 seizeTokens):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0x298637f684da70674f26509b10f07ec2fbc77a335ab1e7d6215a4b2484d8bb52'
	//
	// when user = liquidator
	NumberOfLiquidations int64 //number_of_liquidations
	TotalLiquidated      int64 //total_liquidated
	// when user = borrower
	NumberOfLiquidated int64 //number_of_liquidated
	TotalLiquidatedFor int64 //total_liquidated_for

	// Transfer (index_topic_1 address from, index_topic_2 address to, uint256 amount):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef'"
	//
	// received by user
	NumberOfTransfersIn int64 //number_of_transfers_in
	TotalTransfersIn    int64 //total_transfers_in
	// sent by user
	NumberOfTransfersOut int64 //number_of_transfers_out
	TotalTransfersOut    int64 //total_transfers_out

	BorrowBalance     int64 //borrow_balance
	PercentRedeemable int64 //percent_redeemable
	NumActions        int64 //num_actions
	FirstAction       int64 //first_action
	LastAction        int64 //last_action
	DaysActive        int64 //days_active
	DaysSinceFirst    int64 //days_since_first
	DaysSinceLast     int64 //days_since_last
	EthBorrowingPower int64 //ETH_borrowing_power
	SupplyBalanceEth  int64 //supply_balance_ETH
	BorrowBalanceEth  int64 //borrow_balance_ETH
	CollateralRatio   int64 //collateral_ratio
}
