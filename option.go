package jquants

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"
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
	var err error
	if err = json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to decode index option price error response: %w", err)
	}
	putCallDivision, err := strconv.ParseInt(raw.PutCallDivision, 10, 8)
	if err != nil {
		return fmt.Errorf("failed to decode index option price error response: %w", err)
	}
	iop.Date = *unmarshalTime(raw.Date)
	iop.Code = raw.Code
	iop.WholeDayOpen = unmarshalPrice(raw.WholeDayOpen)
	iop.WholeDayHigh = unmarshalPrice(raw.WholeDayHigh)
	iop.WholeDayLow = unmarshalPrice(raw.WholeDayLow)
	iop.WholeDayClose = unmarshalPrice(raw.WholeDayClose)
	iop.NightSessionOpen = unmarshalPrice(raw.NightSessionOpen)
	iop.NightSessionHigh = unmarshalPrice(raw.NightSessionHigh)
	iop.NightSessionLow = unmarshalPrice(raw.NightSessionLow)
	iop.NightSessionClose = unmarshalPrice(raw.NightSessionClose)
	iop.DaySessionOpen = unmarshalPrice(raw.DaySessionOpen)
	iop.DaySessionHigh = unmarshalPrice(raw.DaySessionHigh)
	iop.DaySessionLow = unmarshalPrice(raw.DaySessionLow)
	iop.DaySessionClose = unmarshalPrice(raw.DaySessionClose)
	iop.Volume = int64(raw.Volume)
	iop.OpenInterest = int64(raw.OpenInterest)
	iop.TurnoverValue = int64(raw.TurnoverValue)
	iop.ContractMonth = raw.ContractMonth
	iop.StrikePrice = int16(raw.StrikePrice)
	iop.VolumeOnlyAuction = unmarshalVolume(raw.VolumeOnlyAuction)
	iop.EmergencyMarginTriggerDivision = raw.EmergencyMarginTriggerDivision
	iop.PutCallDivision = int8(putCallDivision)
	iop.LastTradingDay = unmarshalTime(raw.LastTradingDay)
	iop.SpecialQuotationDay = unmarshalTime(raw.SpecialQuotationDay)
	iop.SettlementPrice = unmarshalPrice(raw.SettlementPrice)
	iop.TheoreticalPrice = unmarshalJSONNumber(raw.TheoreticalPrice)
	iop.BaseVolatility = unmarshalJSONNumber(raw.BaseVolatility)
	iop.UnderlyingPrice = unmarshalJSONNumber(raw.UnderlyingPrice)
	iop.ImpliedVolatility = unmarshalJSONNumber(raw.ImpliedVolatility)
	iop.InterestRate = unmarshalJSONNumber(raw.InterestRate)
	return nil
}

func unmarshalPrice(value interface{}) *int16 {
	switch v := value.(type) {
	case float64:
		f := float32(v)
		i := int16(f)
		return &i
	case string:
		return nil
	default:
		panic(fmt.Errorf("unknown type %T", value))
	}
}

func unmarshalVolume(value interface{}) *int64 {
	switch v := value.(type) {
	case float64:
		i := int64(v)
		return &i
	case string:
		return nil
	default:
		panic(fmt.Errorf("unknown type %T", value))
	}
}

func unmarshalJSONNumber(value interface{}) *json.Number {
	switch v := value.(type) {
	case float64:
		s := strconv.FormatFloat(v, 'f', -1, 64)
		n := json.Number(s)
		return &n
	case string:
		return nil
	default:
		panic(fmt.Errorf("unknown type %T\n", value))
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
	var data = make([]IndexOptionPrice, 0)
	var paginationKey *string
	ctx, cancel := context.WithTimeout(ctx, c.LoopTimeout)
	defer cancel()
	for {
		param := indexOptionPriceParameters{IndexOptionPriceRequest: req, PaginationKey: paginationKey}
		resp, err := c.sendIndexOptionPriceRequest(ctx, param)
		if err != nil {
			if errors.As(err, &InternalServerError{}) {
				slog.Warn("Retrying HTTP request", "error", err.Error())
				time.Sleep(c.RetryInterval)
				continue
			} else {
				return nil, err
			}
		}
		data = append(data, resp.Data...)
		paginationKey = resp.PaginationKey
		if resp.PaginationKey == nil {
			break
		}
	}
	return data, nil
}

func (c *Client) IndexOptionPriceWithChannel(ctx context.Context, req IndexOptionPriceRequest, ch chan<- IndexOptionPrice) error {
	var paginationKey *string
	ctx, cancel := context.WithTimeout(ctx, c.LoopTimeout)
	defer cancel()
	for {
		param := indexOptionPriceParameters{IndexOptionPriceRequest: req, PaginationKey: paginationKey}
		resp, err := c.sendIndexOptionPriceRequest(ctx, param)
		if err != nil {
			if errors.As(err, &InternalServerError{}) {
				slog.Warn("Retrying HTTP request", "error", err.Error())
				time.Sleep(c.RetryInterval)
				continue
			} else {
				return err
			}
		}
		for _, d := range resp.Data {
			ch <- d
		}
		paginationKey = resp.PaginationKey
		if resp.PaginationKey == nil {
			break
		}
	}
	close(ch)
	return nil
}
