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
	Date                           string       `json:"Date"`
	Code                           string       `json:"Code"`
	WholeDayOpen                   *int16       `json:"WholeDayOpen"`
	WholeDayHigh                   *int16       `json:"WholeDayHigh"`
	WholeDayLow                    *int16       `json:"WholeDayLow"`
	WholeDayClose                  *int16       `json:"WholeDayClose"`
	NightSessionOpen               *int16       `json:"NightSessionOpen"`
	NightSessionHigh               *int16       `json:"NightSessionHigh"`
	NightSessionLow                *int16       `json:"NightSessionLow"`
	NightSessionClose              *int16       `json:"NightSessionClose"`
	DaySessionOpen                 *int16       `json:"DaySessionOpen"`
	DaySessionHigh                 *int16       `json:"DaySessionHigh"`
	DaySessionLow                  *int16       `json:"DaySessionLow"`
	DaySessionClose                *int16       `json:"DaySessionClose"`
	Volume                         int64        `json:"Volume"`
	OpenInterest                   int64        `json:"OpenInterest"`
	TurnoverValue                  int64        `json:"TurnoverValue"`
	ContractMonth                  string       `json:"ContractMonth"`
	StrikePrice                    int16        `json:"StrikePrice"`
	VolumeOnlyAuction              *int64       `json:"Volume(OnlyAuction)"`
	EmergencyMarginTriggerDivision string       `json:"EmergencyMarginTriggerDivision"`
	PutCallDivision                int8         `json:"PutCallDivision"`
	LastTradingDay                 *string      `json:"LastTradingDay"`
	SpecialQuotationDay            *string      `json:"SpecialQuotationDay"`
	SettlementPrice                *int16       `json:"SettlementPrice"`
	TheoreticalPrice               *json.Number `json:"TheoreticalPrice"`
	BaseVolatility                 *json.Number `json:"BaseVolatility"`
	UnderlyingPrice                *json.Number `json:"UnderlyingPrice"`
	ImpliedVolatility              *json.Number `json:"ImpliedVolatility"`
	InterestRate                   *json.Number `json:"InterestRate"`
}

func (iop *IndexOptionPrice) UnmarshalJSON(b []byte) error {
	var raw struct {
		Date                           string      `json:"Date"`
		Code                           string      `json:"Code"`
		WholeDayOpen                   interface{} `json:"WholeDayOpen"`
		WholeDayHigh                   interface{} `json:"WholeDayHigh"`
		WholeDayLow                    interface{} `json:"WholeDayLow"`
		WholeDayClose                  interface{} `json:"WholeDayClose"`
		NightSessionOpen               interface{} `json:"NightSessionOpen"`
		NightSessionHigh               interface{} `json:"NightSessionHigh"`
		NightSessionLow                interface{} `json:"NightSessionLow"`
		NightSessionClose              interface{} `json:"NightSessionClose"`
		DaySessionOpen                 interface{} `json:"DaySessionOpen"`
		DaySessionHigh                 interface{} `json:"DaySessionHigh"`
		DaySessionLow                  interface{} `json:"DaySessionLow"`
		DaySessionClose                interface{} `json:"DaySessionClose"`
		Volume                         float64     `json:"Volume"`
		OpenInterest                   float64     `json:"OpenInterest"`
		TurnoverValue                  float64     `json:"TurnoverValue"`
		ContractMonth                  string      `json:"ContractMonth"`
		StrikePrice                    float64     `json:"StrikePrice"`
		VolumeOnlyAuction              interface{} `json:"Volume(OnlyAuction)"`
		EmergencyMarginTriggerDivision string      `json:"EmergencyMarginTriggerDivision"`
		PutCallDivision                string      `json:"PutCallDivision"`
		LastTradingDay                 string      `json:"LastTradingDay"`
		SpecialQuotationDay            string      `json:"SpecialQuotationDay"`
		SettlementPrice                interface{} `json:"SettlementPrice"`
		TheoreticalPrice               interface{} `json:"TheoreticalPrice"`
		BaseVolatility                 interface{} `json:"BaseVolatility"`
		UnderlyingPrice                interface{} `json:"UnderlyingPrice"`
		ImpliedVolatility              interface{} `json:"ImpliedVolatility"`
		InterestRate                   interface{} `json:"InterestRate"`
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
		fmt.Printf("unknown type %T\n", value)
		return nil
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
		fmt.Printf("unknown type %T\n", value)
		return nil
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
	} else {
		return &value
	}
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
