package robinhood

import (
	"bytes"
	"strings"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"encoding/json"

	"net/http"

	"github.com/shopspring/decimal"
)

const ethPrecision = int32(6)
const defaultPrecision = int32(8)

// CryptoOrder is the payload to create a crypto currency order
type CryptoOrder struct {
	AccountID      string  `json:"account_id,omitempty"`
	CurrencyPairID string  `json:"currency_pair_id,omitempty"`
	Price          float64 `json:"price,omitempty"`
	RefID          string  `json:"ref_id,omitempty"`
	Side           string  `json:"side,omitempty"`
	TimeInForce    string  `json:"time_in_force,omitempty"`
	Quantity       float64 `json:"quantity,omitempty"`
	Type           string  `json:"type,omitempty"`
}

type Execution struct {
	EffectivePrice float64 `json:"effective_Price,string"`
	ID             string  `json:"id"`
	Quantity       float64 `json:"quantity,string"`
	Timestamp      string  `json:"timestamp"`
}

// CryptoOrderOutput holds the response from api
type CryptoOrderOutput struct {
	Meta
	Account            string      `json:"account_id"`
	AveragePrice       float64     `json:"average_price,string"`
	CancelURL          string      `json:"cancel_url"`
	CumulativeQuantity string      `json:"cumulative_quantity"`
	CurrencyPairID     string      `json:"currency_pair_id"`
	Executions         []Execution `json:"executions"`
	ID                 string      `json:"id"`
	LastTransactionAt  string      `json:"last_transaction_at"`
	Price              float64     `json:"price,string"`
	Quantity           string      `json:"quantity"`
	RejectReason       string      `json:"reject_reason"`
	Side               string      `json:"side"`
	State              string      `json:"state"`
	StopPrice          float64     `json:"stop_price,string"`
	TimeInForce        string      `json:"time_in_force"`
	Type               string      `json:"type"`

	client *Client
}

// CryptoOrderOpts encapsulates differences between order types
type CryptoOrderOpts struct {
	Side            OrderSide
	Type            OrderType
	AmountInDollars float64
	Quantity        float64
	Price           float64
	TimeInForce     TimeInForce
	ExtendedHours   bool
	Stop, Force     bool
}

// CryptoOrder will actually place the order
func (c *Client) CryptoOrder(ctx context.Context, cryptoPair CryptoCurrencyPair, o CryptoOrderOpts) (*CryptoOrderOutput, error) {
	var amountInDollars = decimal.NewFromFloat32(float32(o.AmountInDollars))
	var price = decimal.NewFromFloat32(float32(o.Price))
	var precision = defaultPrecision
	if cryptoPair.CyrptoAssetCurrency.Code == "ETH" {
		precision = ethPrecision
	}
	var quantity = amountInDollars.DivRound(price, precision)
	exactQuantity, _ := quantity.Float64()
	a := CryptoOrder{
		AccountID:      c.CryptoAccount.ID,
		CurrencyPairID: cryptoPair.ID,
		Quantity:       exactQuantity,
		Price:          o.Price,
		RefID:          uuid.New().String(),
		Side:           strings.ToLower(o.Side.String()),
		TimeInForce:    strings.ToLower(o.TimeInForce.String()),
		Type:           strings.ToLower(o.Type.String()),
	}

	payload, err := json.Marshal(a)

	if err != nil {
		return nil, err
	}

	post, err := http.NewRequest("POST", EPCryptoOrders, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("could not create Crypto http.Request: %v", err)
	}

	post.Header.Add("Content-Type", "application/json")

	var out CryptoOrderOutput
	err = c.DoAndDecode(ctx, post, &out)
	out.client = c
	return &out, err
}

// Cancel will cancel the order.
func (o CryptoOrderOutput) Cancel(ctx context.Context) error {
	post, err := http.NewRequest("POST", o.CancelURL, nil)
	if err != nil {
		return err
	}

	var output CryptoOrderOutput
	err = o.client.DoAndDecode(ctx, post, &output)

	if err != nil {
		return errors.Wrap(err, "could not decode response")
	}

	if output.RejectReason != "" {
		return errors.New(output.RejectReason)
	}

	return nil
}

// GetCryptoOrder will get the order info from robinhood
func (c *Client) GetCryptoOrder(ctx context.Context, orderID string) (*CryptoOrderOutput, error) {
	url := EPCryptoOrders + orderID
	get, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var output CryptoOrderOutput
	err = c.DoAndDecode(ctx, get, &output)

	if err != nil {
		return nil, err
	}

	output.client = c
	return &output, nil
}
