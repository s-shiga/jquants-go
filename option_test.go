package jquants

import (
	"encoding/json"
	"testing"
)

func TestClient_IndexOptionPrice(t *testing.T) {
	req := IndexOptionPriceRequest{Date: "2026-07-17"}
	checkEndpoint(t, "/derivatives/bars/daily/options/225", "date=2026-07-17", `{"Code":"123","PCDiv":"1","O":980,"EO":"","Strike":20000,"Theo":974.641}`, true, IndexOptionPrice{Code: "123", PutCallDivision: 1, WholeDayOpen: ptr(int32(980)), StrikePrice: 20000, TheoreticalPrice: ptr(json.Number("974.641"))}, func(c *Client) ([]IndexOptionPrice, error) {
		return c.IndexOptionPrice(t.Context(), req)
	})
}

func TestClient_IndexOptionPriceWithChannel(t *testing.T) {
	req := IndexOptionPriceRequest{Date: "2026-07-17"}
	checkEndpoint(t, "/derivatives/bars/daily/options/225", "date=2026-07-17", `{"Code":"123","PCDiv":"1","O":980,"EO":"","Strike":20000,"Theo":974.641}`, true, IndexOptionPrice{Code: "123", PutCallDivision: 1, WholeDayOpen: ptr(int32(980)), StrikePrice: 20000, TheoreticalPrice: ptr(json.Number("974.641"))}, func(c *Client) ([]IndexOptionPrice, error) {
		return collectChannel(func(ch chan<- IndexOptionPrice) error { return c.IndexOptionPriceWithChannel(t.Context(), req, ch) })
	})
}

func TestClient_OptionPrice(t *testing.T) {
	req := OptionPriceRequest{Date: "2026-07-17", Category: ptr("TOPIXE"), ContractFlag: ptr("1")}
	checkEndpoint(t, "/derivatives/bars/daily/options", "date=2026-07-17&category=TOPIXE&contract_flag=1", `{"Code":"123","PCDiv":"2","O":980.5,"EO":"","Strike":20000.5}`, true, OptionPrice{Code: "123", PutCallDivision: 2, WholeDayOpen: ptr(json.Number("980.5")), StrikePrice: 20000.5}, func(c *Client) ([]OptionPrice, error) {
		return c.OptionPrice(t.Context(), req)
	})
}

func TestClient_OptionPriceWithChannel(t *testing.T) {
	req := OptionPriceRequest{Date: "2026-07-17", Category: ptr("TOPIXE"), ContractFlag: ptr("1")}
	checkEndpoint(t, "/derivatives/bars/daily/options", "date=2026-07-17&category=TOPIXE&contract_flag=1", `{"Code":"123","PCDiv":"2","O":980.5,"EO":"","Strike":20000.5}`, true, OptionPrice{Code: "123", PutCallDivision: 2, WholeDayOpen: ptr(json.Number("980.5")), StrikePrice: 20000.5}, func(c *Client) ([]OptionPrice, error) {
		return collectChannel(func(ch chan<- OptionPrice) error { return c.OptionPriceWithChannel(t.Context(), req, ch) })
	})
}
