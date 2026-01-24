package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type IndexOptionPrice struct {
	Date                           string
	Code                           string
	WholeDayOpen                   *int16
	WholeDayHigh                   *int16
	WholeDayLow                    *int16
	WholeDayClose                  *int16
	NightSessionOpen               *int16
	NightSessionHigh               *int16
	NightSessionLow                *int16
	NightSessionClose              *int16
	DaySessionOpen                 *int16
	DaySessionHigh                 *int16
	DaySessionLow                  *int16
	DaySessionClose                *int16
	Volume                         int64
	OpenInterest                   int64
	TurnoverValue                  int64
	ContractMonth                  string
	StrikePrice                    int16
	VolumeOnlyAuction              *int64
	EmergencyMarginTriggerDivision string
	PutCallDivision                int8
	LastTradingDay                 *string
	SpecialQuotationDay            *string
	SettlementPrice                *int16
	TheoreticalPrice               *json.Number
	BaseVolatility                 *json.Number
	UnderlyingPrice                *json.Number
	ImpliedVolatility              *json.Number
	InterestRate                   *json.Number
}

// unmarshaler accumulates errors during unmarshaling, allowing cleaner code flow.
type unmarshaler struct {
	err error
}

func (u *unmarshaler) price(v interface{}) *int16 {
	if u.err != nil {
		return nil
	}
	result, err := unmarshalPrice(v)
	u.err = err
	return result
}

func (u *unmarshaler) volume(v interface{}) *int64 {
	if u.err != nil {
		return nil
	}
	result, err := unmarshalVolume(v)
	u.err = err
	return result
}

func (u *unmarshaler) jsonNumber(v interface{}) *json.Number {
	if u.err != nil {
		return nil
	}
	result, err := unmarshalJSONNumber(v)
	u.err = err
	return result
}

func (iop *IndexOptionPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date                           string      `json:"Date"`
		Code                           string      `json:"Code"`
		WholeDayOpen                   interface{} `json:"O"`
		WholeDayHigh                   interface{} `json:"H"`
		WholeDayLow                    interface{} `json:"L"`
		WholeDayClose                  interface{} `json:"C"`
		NightSessionOpen               interface{} `json:"EO"`
		NightSessionHigh               interface{} `json:"EH"`
		NightSessionLow                interface{} `json:"EL"`
		NightSessionClose              interface{} `json:"EC"`
		DaySessionOpen                 interface{} `json:"AO"`
		DaySessionHigh                 interface{} `json:"AH"`
		DaySessionLow                  interface{} `json:"AL"`
		DaySessionClose                interface{} `json:"AC"`
		Volume                         float64     `json:"Vo"`
		OpenInterest                   float64     `json:"OI"`
		TurnoverValue                  float64     `json:"Va"`
		ContractMonth                  string      `json:"CM"`
		StrikePrice                    float64     `json:"Strike"`
		VolumeOnlyAuction              interface{} `json:"VoOA"`
		EmergencyMarginTriggerDivision string      `json:"EmMrgnTrgDiv"`
		PutCallDivision                string      `json:"PCDiv"`
		LastTradingDay                 string      `json:"LTD"`
		SpecialQuotationDay            string      `json:"SQD"`
		SettlementPrice                interface{} `json:"Settle"`
		TheoreticalPrice               interface{} `json:"Theo"`
		BaseVolatility                 interface{} `json:"BaseVol"`
		UnderlyingPrice                interface{} `json:"UnderPx"`
		ImpliedVolatility              interface{} `json:"IV"`
		InterestRate                   interface{} `json:"IR"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to decode index option price: %w", err)
	}
	putCallDivision, err := strconv.ParseInt(raw.PutCallDivision, 10, 8)
	if err != nil {
		return fmt.Errorf("failed to parse put/call division: %w", err)
	}

	u := &unmarshaler{}

	iop.Date = *unmarshalTime(raw.Date)
	iop.Code = raw.Code
	iop.WholeDayOpen = u.price(raw.WholeDayOpen)
	iop.WholeDayHigh = u.price(raw.WholeDayHigh)
	iop.WholeDayLow = u.price(raw.WholeDayLow)
	iop.WholeDayClose = u.price(raw.WholeDayClose)
	iop.NightSessionOpen = u.price(raw.NightSessionOpen)
	iop.NightSessionHigh = u.price(raw.NightSessionHigh)
	iop.NightSessionLow = u.price(raw.NightSessionLow)
	iop.NightSessionClose = u.price(raw.NightSessionClose)
	iop.DaySessionOpen = u.price(raw.DaySessionOpen)
	iop.DaySessionHigh = u.price(raw.DaySessionHigh)
	iop.DaySessionLow = u.price(raw.DaySessionLow)
	iop.DaySessionClose = u.price(raw.DaySessionClose)
	iop.Volume = int64(raw.Volume)
	iop.OpenInterest = int64(raw.OpenInterest)
	iop.TurnoverValue = int64(raw.TurnoverValue)
	iop.ContractMonth = raw.ContractMonth
	iop.StrikePrice = int16(raw.StrikePrice)
	iop.VolumeOnlyAuction = u.volume(raw.VolumeOnlyAuction)
	iop.EmergencyMarginTriggerDivision = raw.EmergencyMarginTriggerDivision
	iop.PutCallDivision = int8(putCallDivision)
	iop.LastTradingDay = unmarshalTime(raw.LastTradingDay)
	iop.SpecialQuotationDay = unmarshalTime(raw.SpecialQuotationDay)
	iop.SettlementPrice = u.price(raw.SettlementPrice)
	iop.TheoreticalPrice = u.jsonNumber(raw.TheoreticalPrice)
	iop.BaseVolatility = u.jsonNumber(raw.BaseVolatility)
	iop.UnderlyingPrice = u.jsonNumber(raw.UnderlyingPrice)
	iop.ImpliedVolatility = u.jsonNumber(raw.ImpliedVolatility)
	iop.InterestRate = u.jsonNumber(raw.InterestRate)

	return u.err
}

func unmarshalPrice(value interface{}) (*int16, error) {
	switch v := value.(type) {
	case float64:
		f := float32(v)
		i := int16(f)
		return &i, nil
	case string:
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unmarshalPrice: unknown type %T", v)
	}
}

func unmarshalVolume(value interface{}) (*int64, error) {
	switch v := value.(type) {
	case float64:
		i := int64(v)
		return &i, nil
	case string:
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unmarshalVolume: unknown type %T", v)
	}
}

func unmarshalJSONNumber(value interface{}) (*json.Number, error) {
	switch v := value.(type) {
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		n := json.Number(s)
		return &n, nil
	case string:
		return nil, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unmarshalJSONNumber: unknown type %T", v)
	}
}

func unmarshalTime(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

type IndexOptionPriceRequest struct {
	Date string
}

type indexOptionPriceParameters struct {
	IndexOptionPriceRequest
	PaginationKey *string
}

func (p indexOptionPriceParameters) values() (url.Values, error) {
	v := url.Values{}
	v.Add("date", p.Date)
	if p.PaginationKey != nil {
		v.Add("pagination_key", *p.PaginationKey)
	}
	return v, nil
}

type indexOptionPriceResponse struct {
	Data          []IndexOptionPrice `json:"data"`
	PaginationKey *string            `json:"pagination_key"`
}

func (r indexOptionPriceResponse) getData() []IndexOptionPrice { return r.Data }
func (r indexOptionPriceResponse) getPaginationKey() *string   { return r.PaginationKey }

func (c *Client) sendIndexOptionPriceRequest(ctx context.Context, params indexOptionPriceParameters) (indexOptionPriceResponse, error) {
	var r indexOptionPriceResponse
	resp, err := c.sendRequest(ctx, "/derivatives/bars/daily/options/225", params)
	if err != nil {
		return r, fmt.Errorf("failed to send GET request: %w", err)
	}
	if resp.StatusCode != 200 {
		return r, handleErrorResponse(resp)
	}
	if err = decodeResponse(resp, &r); err != nil {
		return r, fmt.Errorf("failed to decode HTTP response: %w", err)
	}
	return r, nil
}

func (c *Client) IndexOptionPrice(ctx context.Context, req IndexOptionPriceRequest) ([]IndexOptionPrice, error) {
	return fetchAllPages(ctx, c, func(ctx context.Context, paginationKey *string) (indexOptionPriceResponse, error) {
		params := indexOptionPriceParameters{IndexOptionPriceRequest: req, PaginationKey: paginationKey}
		return c.sendIndexOptionPriceRequest(ctx, params)
	})
}

func (c *Client) IndexOptionPriceWithChannel(ctx context.Context, req IndexOptionPriceRequest, ch chan<- IndexOptionPrice) error {
	return fetchAllPagesWithChannel(ctx, c, ch, func(ctx context.Context, paginationKey *string) (indexOptionPriceResponse, error) {
		params := indexOptionPriceParameters{IndexOptionPriceRequest: req, PaginationKey: paginationKey}
		return c.sendIndexOptionPriceRequest(ctx, params)
	})
}
