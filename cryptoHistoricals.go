package robihood

import "fmt"

type Historical struct {
	BeginsAt     time.Time `json:"begins_at"`
	OpenPrice    float64    `json:"open_price,string"`
	ClosePrice   float64    `json:"close_price,string"`
	HighPrice    float64    `json:"high_price,string"`
	LowPrice     float64    `json:"low_price,string"`
	Volume       int       `json:"volume"`
	Session      string    `json:"session"`
	Interpolated bool      `json:"interpolated"`
}

func (c *Client) GetDailyHistoricals(cryptoID string) ([]Historical, error) {
	url := EPMarket + fmt.Sprintf("forex/historicals/${0}", cryptoID) + "/?bound=24_7&interval=day&span=week"
	var r struct{ Results []Historical }
	err := c.GetAndDecode(url, &r)
	return r.Results, err
}