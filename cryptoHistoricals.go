package robinhood

import "fmt"

// HistoricalData represents class ohlc data i.e open, high, low, close
type HistoricalData struct {
	DataPoints []Historical `json:"data_points"`
	Bounds     string       `json:"bounds"`
	Interval   string       `json:"interval"`
	Symbol     string       `json:"symbol"`
	ID         string       `json:"id"`
}

// Historical data represents class ohlc data i.e open, high, low, close
type Historical struct {
	BeginsAt     string  `json:"begins_at"`
	OpenPrice    float64 `json:"open_price,string"`
	ClosePrice   float64 `json:"close_price,string"`
	HighPrice    float64 `json:"high_price,string"`
	LowPrice     float64 `json:"low_price,string"`
	Volume       int     `json:"volume"`
	Session      string  `json:"session"`
	Interpolated bool    `json:"interpolated"`
}

// GetCryptoHistoricals will give the high low open, close data fro the given symbol and span
func (c *Client) GetCryptoHistoricals(cryptoID string, interval string, span string, bounds string) (HistoricalData, error) {
	var r = HistoricalData{}

	// Set defaults
	if interval == "" { // Options: "15second", "5minute", "10minute", "hour", "day", "week"
		interval = "hour"
	}
	if span == "" { // Options: "hour", "day", "week", "month", "3month", "year", "5year"
		span = "week"
	}
	if bounds == "" { // Options: "Regular" is 6 hours a day, "trading" is 9 hours a day, "extended" is 16 hours a day, "24_7" is 24 hours a day.
		bounds = "24_7"
	}

	// Build the URL
	url := EPMarket + fmt.Sprintf("forex/historicals/%s", cryptoID) + "/?bounds=" + bounds + "&interval=" + interval + "&span=" + span

	// Get the data
	err := c.GetAndDecode(url, &r)
	return r, err
}

// GetDailyHistoricals will give daily high, low, ope, close data for the given symbol
func (c *Client) GetDailyHistoricals(cryptoID string) (HistoricalData, error) {
	url := EPMarket + fmt.Sprintf("forex/historicals/%s", cryptoID) + "/?bounds=24_7&interval=day&span=week"
	var r = HistoricalData{}
	err := c.GetAndDecode(url, &r)
	return r, err
}
