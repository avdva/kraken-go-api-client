package kraken

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// The following consts match pair names in kraken api.
const (
	XETHXXBT = "XETHXXBT"
	XETHZCAD = "XETHZCAD"
	XETHZEUR = "XETHZEUR"
	XETHZGBP = "XETHZGBP"
	XETHZJPY = "XETHZJPY"
	XETHZUSD = "XETHZUSD"
	XLTCZCAD = "XLTCZCAD"
	XLTCZEUR = "XLTCZEUR"
	XLTCZUSD = "XLTCZUSD"
	XXBTXLTC = "XXBTXLTC"
	XXBTXNMC = "XXBTXNMC"
	XXBTXXDG = "XXBTXXDG"
	XXBTXXLM = "XXBTXXLM"
	XXBTXXRP = "XXBTXXRP"
	XXBTZCAD = "XXBTZCAD"
	XXBTZEUR = "XXBTZEUR"
	XXBTZGBP = "XXBTZGBP"
	XXBTZJPY = "XXBTZJPY"
	XXBTZUSD = "XXBTZUSD"
)

// Response wraps the Kraken API JSON response
type Response struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

// TimeResponse represents the server's time
type TimeResponse struct {
	// Unix timestamp
	Unixtime int64
	// RFC 1123 time format
	Rfc1123 string
}

// AssetPairsResponse includes asset pair informations
type AssetPairsResponse struct {
	Infos map[string]AssetPairInfo
}

// AssetPairInfo represents asset pair information
type AssetPairInfo struct {
	// Alternate pair name
	Altname string `json:"altname"`
	// Asset class of base component
	AssetClassBase string `json:"aclass_base"`
	// Asset id of base component
	Base string `json:"base"`
	// Asset class of quote component
	AssetClassQuote string `json:"aclass_quote"`
	// Asset id of quote component
	Quote string `json:"quote"`
	// Volume lot size
	Lot string `json:"lot"`
	// Scaling decimal places for pair
	PairDecimals int `json:"pair_decimals"`
	// Scaling decimal places for volume
	LotDecimals int `json:"lot_decimals"`
	// Amount to multiply lot volume by to get currency volume
	LotMultiplier int `json:"lot_multiplier"`
	// Array of leverage amounts available when buying
	LeverageBuy []float32 `json:"leverage_buy"`
	// Array of leverage amounts available when selling
	LeverageSell []float32 `json:"leverage_sell"`
	// Fee schedule array in [volume, percent fee] tuples
	Fees [][]float64 `json:"fees"`
	// // Maker fee schedule array in [volume, percent fee] tuples (if on maker/taker)
	FeesMaker [][]float64 `json:"fees_maker"`
	// // Volume discount currency
	FeeVolumeCurrency string `json:"fee_volume_currency"`
	// Margin call level
	MarginCall int `json:"margin_call"`
	// Stop-out/Liquidation margin level
	MarginStop int `json:"margin_stop"`
}

// AssetsResponse includes asset informations
type AssetsResponse struct {
	Infos map[string]AssetInfo
}

// AssetInfo represents an asset information
type AssetInfo struct {
	// Alternate name
	Altname string
	// Asset class
	AssetClass string
	// Scaling decimal places for record keeping
	Decimals int
	// Scaling decimal places for output display
	DisplayDecimals int `json:"display_decimals"`
}

// TickerResponse includes the requested ticker pairs
type TickerResponse struct {
	Infos map[string]PairTickerInfo
}

// PairTickerInfo represents ticker information for a pair
type PairTickerInfo struct {
	// Ask array(<price>, <whole lot volume>, <lot volume>)
	Ask []string `json:"a"`
	// Bid array(<price>, <whole lot volume>, <lot volume>)
	Bid []string `json:"b"`
	// Last trade closed array(<price>, <lot volume>)
	Close []string `json:"c"`
	// Volume array(<today>, <last 24 hours>)
	Volume []string `json:"v"`
	// Volume weighted average price array(<today>, <last 24 hours>)
	VolumeAveragePrice []string `json:"p"`
	// Number of trades array(<today>, <last 24 hours>)
	Trades []int `json:"t"`
	// Low array(<today>, <last 24 hours>)
	Low []string `json:"l"`
	// High array(<today>, <last 24 hours>)
	High []string `json:"h"`
	// Today's opening price
	OpeningPrice float32 `json:"o,string"`
}

// OrderBookItem is a piece of information about an order.
type OrderBookItem struct {
	Price  decimal.Decimal
	Amount decimal.Decimal
	Ts     int64
}

// UnmarshalJSON takes a json array from kraken and converts it into an OrderBookItem.
func (o *OrderBookItem) UnmarshalJSON(data []byte) error {
	var arr []interface{}
	var err error
	if err = json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 3 {
		return errors.Errorf("arr len %d != 3", len(arr))
	}
	if price, ok := arr[0].(string); ok {
		if o.Price, err = decimal.NewFromString(price); err != nil {
			return errors.Wrap(err, "failed to parse price")
		}
	} else {
		return errors.Errorf("price must be string, not %T", arr[0])
	}
	if amount, ok := arr[1].(string); ok {
		if o.Amount, err = decimal.NewFromString(amount); err != nil {
			return errors.Wrap(err, "failed to parse amount")
		}
	} else {
		return errors.Errorf("amount must be string, not %T", arr[1])
	}
	if ts, ok := arr[2].(int); ok {
		o.Ts = int64(ts)
	} else {
		return errors.Errorf("amount must be int, not %T", arr[2])
	}
	return nil
}

type OrderBookResponse map[string]OrderBook

type OrderBook struct {
	Asks []OrderBookItem
	Bids []OrderBookItem
}
