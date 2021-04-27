package robinhood

import (
	"context"
	"strings"
)

// A Quote is a representation of the data returned by the Robinhood API for
// current stock quotes
type Quote struct {
	AdjustedPreviousClose       float64 `json:"adjusted_previous_close,string"`
	AskPrice                    float64 `json:"ask_price,string"`
	AskSize                     int     `json:"ask_size"`
	BidPrice                    float64 `json:"bid_price,string"`
	BidSize                     int     `json:"bid_size"`
	LastExtendedHoursTradePrice float64 `json:"last_extended_hours_trade_price,string"`
	LastTradePrice              float64 `json:"last_trade_price,string"`
	PreviousClose               float64 `json:"previous_close,string"`
	PreviousCloseDate           string  `json:"previous_close_date"`
	Symbol                      string  `json:"symbol"`
	TradingHalted               bool    `json:"trading_halted"`
	UpdatedAt                   string  `json:"updated_at"`
}

// CryptoQuote is a representation of data returned by robinhood api for cryto quote
type CryptoQuote struct {
	AskPrice  float64 `json:"ask_price,string"`
	BidPrice  float64 `json:"bid_price,string"`
	HighPrice float64 `json:"high_price,string"`
	ID        string  `json:"id"`
	LowPrice  float64 `json:"low_price,string"`
	MarkPrice float64 `json:"mark_price,string"`
	OpenPrice float64 `json:"open_price,string"`
	Symbol    string  `json:"symbol"`
	Volume    string  `json:"volume"`
}

// GetQuote returns all the latest stock quotes for the list of stocks provided
func (c *Client) GetQuote(ctx context.Context, stocks ...string) ([]Quote, error) {
	url := EPQuotes + "?symbols=" + strings.Join(stocks, ",")
	var r struct{ Results []Quote }
	err := c.GetAndDecode(ctx, url, &r)
	return r.Results, err
}

// Price returns the proper stock price even after hours
func (q Quote) Price() float64 {
	if IsRegularTradingTime() {
		return q.LastTradePrice
	}
	return q.LastExtendedHoursTradePrice
}

// GetCryptoQuote will return an array of current quotes
// these will change almost every second
func (c *Client) GetCryptoQuote(ctx context.Context, cryptoIds ...string) ([]CryptoQuote, error) {
	url := EPMarket + "forex/quotes/?ids=" + strings.Join(cryptoIds, ",")
	var r struct{ Results []CryptoQuote }
	err := c.GetAndDecode(ctx, url, &r)
	return r.Results, err
}
