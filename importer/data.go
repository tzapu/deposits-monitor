package importer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"time"

	"github.com/boltdb/bolt"
	"github.com/corpetty/go-alethio-api/alethio"

	"github.com/tzapu/deposits-monitor/helper"
)

// Buckets
const (
	SettingsBucket = "settings"
	AccountsBucket = "accounts"
)

var Buckets = []string{SettingsBucket, AccountsBucket}

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

func (imp *Importer) ProcessTransfers(address, market string, events []alethio.LogEntry) {
	err := imp.data.DB.Batch(func(tx *bolt.Tx) error {
		accounts := tx.Bucket([]byte(AccountsBucket))
		for _, event := range events {
			fromAddress := event.Attributes.EventDecoded.Inputs[0].Value
			toAddress := event.Attributes.EventDecoded.Inputs[0].Value
			v := event.Attributes.EventDecoded.Inputs[2].Value
			value := new(big.Int)
			value.SetString(v, 10)

			// from
			fromAcc := getAccount(accounts, market, fromAddress)
			fromAcc.NumberOfTransfersOut++
			fromAcc.TotalTransfersOut.Add(&fromAcc.TotalTransfersOut, value)
			putAccount(accounts, &fromAcc)

			// to
			toAcc := getAccount(accounts, market, toAddress)
			toAcc.NumberOfTransfersIn++
			toAcc.TotalTransfersIn.Add(&toAcc.TotalTransfersIn, value)
			putAccount(accounts, &toAcc)
		}
		return nil
	})

	helper.FatalIfError(err, "batch update from events")
}

func (imp *Importer) GetAccounts() []Account {
	var accounts []Account
	err := imp.data.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(AccountsBucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var a Account
			err := json.Unmarshal(v, &a)
			helper.FatalIfError(err, "get accounts unmarshal")
			accounts = append(accounts, a)
		}

		return nil
	})
	helper.FatalIfError(err, "read accounts")

	return accounts
}

func getAccount(bucket *bolt.Bucket, market string, address string) Account {
	key := []byte(fmt.Sprintf("%s:%s", market, address))
	ba := bucket.Get(key)
	var account Account
	if len(ba) > 0 {
		err := json.Unmarshal(ba, &account)
		helper.FatalIfError(err, "unmarshall account", key)
	} else {
		account.CToken = market
		account.Address = address
	}
	return account
}

func putAccount(bucket *bolt.Bucket, account *Account) {
	key := []byte(fmt.Sprintf("%s:%s", account.CToken, account.Address))
	ba, err := json.Marshal(account)
	helper.FatalIfError(err, "encode account", key)
	err = bucket.Put(key, ba)
	helper.FatalIfError(err, "put account", key)
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

type Account struct {
	Address string `csv:"address"`

	// cToken market contract addresses
	// 0x6c8c6b02e7b2be14d4fa6022dfd6d75921d90e4e - cBAT
	// 0xf5dce57282a584d2746faf1593d3121fcac444dc - cDAI
	// 0x158079ee67fce2f58472a96584a73c7ab9ac95c1 - cREP
	// 0xb3319f5d18bc0d84dd1b4825dcde5d5f7266d407 - cZRX
	// 0x39aa39c021dfbae8fac545936693ac917d5e7563 - cUSDC
	// 0xc11b1268c1a384e55c48c2391d8d480264a3a7f4 - cWBTC
	// 0x4ddc2d193948926d02f9b1fe9e1daa0718270ed5 - cETH
	CToken string `csv:"cToken"`

	// "Mint (address minter, uint256 mintAmount, uint256 mintTokens):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f'
	NumberOfMints int64 `csv:"number_of_mints"`
	TotalMinted   int64 `csv:"total_minted"`

	// "Redeem (address redeemer, uint256 redeemAmount, uint256 redeemTokens):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0xe5b754fb1abb7f01b499791d0b820ae3b6af3424ac1c59768edb53f4ec31a929'"
	NumberOfRedeems int64 `csv:"number_of_redeems"`
	TotalRedeemed   int64 `csv:"total_redeemed"`
	SupplyBalance   int64 `csv:"supply_balance"`

	// Borrow (address borrower, uint256 borrowAmount, uint256 accountBorrows, uint256 totalBorrows):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0x13ed6866d4e1ee6da46f845c46d7e54120883d75c5ea9a2dacc1c4ca8984ab80'
	NumberOfBorrows    int64 `csv:"number_of_borrows"`
	TotalBorrowed      int64 `csv:"total_borrowed"`
	NumberOfRepayments int64 `csv:"number_of_repayments"`
	TotalRepaid        int64 `csv:"total_repaid"`

	// LiquidateBorrow (address liquidator, address borrower, uint256 repayAmount, address cTokenCollateral, uint256 seizeTokens):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0x298637f684da70674f26509b10f07ec2fbc77a335ab1e7d6215a4b2484d8bb52'
	//
	// when user = liquidator
	NumberOfLiquidations int64 `csv:"number_of_liquidations"`
	TotalLiquidated      int64 `csv:"total_liquidated"`
	// when user = borrower
	NumberOfLiquidated int64 `csv:"number_of_liquidated"`
	TotalLiquidatedFor int64 `csv:"total_liquidated_for"`

	// Transfer (index_topic_1 address from, index_topic_2 address to, uint256 amount):
	// - contractAdr in list of cToken addresses
	// - topics[0] = '0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef'"
	//
	// received by user
	NumberOfTransfersIn int64   `csv:"number_of_transfers_in"`
	TotalTransfersIn    big.Int `csv:"total_transfers_in"`
	// sent by user
	NumberOfTransfersOut int64   `csv:"number_of_transfers_out"`
	TotalTransfersOut    big.Int `csv:"total_transfers_out"`

	BorrowBalance     int64 `csv:"borrow_balance"`
	PercentRedeemable int64 `csv:"percent_redeemable"`
	NumActions        int64 `csv:"num_actions"`
	FirstAction       int64 `csv:"first_action"`
	LastAction        int64 `csv:"last_action"`
	DaysActive        int64 `csv:"days_active"`
	DaysSinceFirst    int64 `csv:"days_since_first"`
	DaysSinceLast     int64 `csv:"days_since_last"`
	EthBorrowingPower int64 `csv:"ETH_borrowing_power"`
	SupplyBalanceEth  int64 `csv:"supply_balance_ETH"`
	BorrowBalanceEth  int64 `csv:"borrow_balance_ETH"`
	CollateralRatio   int64 `csv:"collateral_ratio"`
}
