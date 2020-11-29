package robinhood

import "fmt"

// HistoricalData represents class ohlc data i.e open, high, low, close
type HistoricalData struct {
	Historicals []Historical `json:"data_points"`
	Bounds      string       `json:"bounds"`
	Interval    string       `json:"interval"`
	Symbol      string       `json:"symbol"`
	ID          string       `json:"id"`
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

// GetDailyHistoricals will give daily high, low, ope, close data for the given symbol
func (c *Client) GetDailyHistoricals(cryptoID string) (HistoricalData, error) {
	url := EPMarket + fmt.Sprintf("forex/historicals/%s", cryptoID) + "/?bounds=24_7&interval=day&span=week"
	var r struct{ Results HistoricalData }
	err := c.GetAndDecode(url, &r)
	return r.Results, err
}
