package robinhood

import "net/url"

type Position struct {
	Meta
	Account                 string  `json:"account"`
	AverageBuyPrice         float64 `json:"average_buy_price,string"`
	Instrument              string  `json:"instrument"`
	IntradayAverageBuyPrice float64 `json:"intraday_average_buy_price,string"`
	IntradayQuantity        float64 `json:"intraday_quantity,string"`
	Quantity                float64 `json:"quantity,string"`
	SharesHeldForBuys       float64 `json:"shares_held_for_buys,string"`
	SharesHeldForSells      float64 `json:"shares_held_for_sells,string"`
}

// CryptoCurrency represents a sub object listed in CryptoPosition
type CryptoCurrency struct {
	BrandColor string  `json:"brand_color"`
	Code       string  `json:"code"`
	ID         string  `json:"id"`
	Increment  float64 `json:"increment,string"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
}

// CostBases represents the actual cost the robinhood user paid for asset
type CostBases struct {
	CurrencyID      string  `json:"currency_id"`
	DirectCostBasis float64 `json:"direct_cost_basis,string"`
	DirectQuantity  float64 `json:"direct_quantity,string"`
	ID              string  `json:"id"`
}

// CryptoPosition returns all crypto position associated with an account
type CryptoPosition struct {
	Meta
	AccountID           string         `json:"account_id"`
	ID                  string         `json:"id"`
	Currency            CryptoCurrency `json:"currency"`
	Cost                CostBases      `json:"cost_bases"`
	Quantity            float64        `json:"quantity,string"`
	QuantityAvailable   float64        `json:"quantity_available,string"`
	QuantityHeldForBuy  float64        `json:"quantity_held_for_buy,string"`
	QuantityHeldForSell float64        `json:"quantity_held_for_sell,string"`
}

// GetPositions returns all the positions associated with an account.
func (c *Client) GetPositions(a Account) ([]Position, error) {
	return c.GetPositionsParams(a, PositionParams{})
}

// PositionParams encapsulates parameters known to the RobinHood positions API
// endpoint.
type PositionParams struct {
	NonZero bool
}

// Encode returns the query string associated with the requested parameters
func (p PositionParams) encode() string {
	v := url.Values{}
	if p.NonZero {
		v.Set("nonzero", "true")
	}
	return v.Encode()
}

// GetPositionsParams returns all the positions associated with a count, but
// passes the encoded PositionsParams object along to the RobinHood API as part
// of the query string.
func (c *Client) GetPositionsParams(a Account, p PositionParams) ([]Position, error) {
	u, err := url.Parse(a.Positions)
	if err != nil {
		return nil, err
	}
	u.RawQuery = p.encode()

	var r struct{ Results []Position }
	return r.Results, c.GetAndDecode(u.String(), &r)
}

// GetCryptoPositions returns all positions associated with the account
func (c *Client) GetCryptoPositions() ([]CryptoPosition, error) {
	var r struct{ Results []CryptoPosition }
	err := c.GetAndDecode(EPCryptoHoldings, &r)
	if err != nil {
		return nil, err
	}

	return r.Results, err
}
