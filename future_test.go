package jquants

import (
	"encoding/json"
	"testing"
)

func TestClient_FuturesPrice(t *testing.T) {
	req := FuturesPriceRequest{Date: "2026-07-17", Category: ptr("TOPIXF"), ContractFlag: ptr("1")}
	checkEndpoint(t, "/derivatives/bars/daily/futures", "date=2026-07-17&category=TOPIXF&contract_flag=1", `{"Code":"123","ProdCat":"TOPIXF","O":2700.5,"EO":"","OI":120,"LTD":"2026-09-10","CCMFlag":"1"}`, true, FuturesPrice{Code: "123", ProductCategory: "TOPIXF", WholeDayOpen: ptr(json.Number("2700.5")), OpenInterest: ptr(int64(120)), LastTradingDay: ptr("2026-09-10"), CentralContractMonthFlag: "1"}, func(c *Client) ([]FuturesPrice, error) {
		return c.FuturesPrice(t.Context(), req)
	})
}

func TestClient_FuturesPriceWithChannel(t *testing.T) {
	req := FuturesPriceRequest{Date: "2026-07-17", Category: ptr("TOPIXF"), ContractFlag: ptr("1")}
	checkEndpoint(t, "/derivatives/bars/daily/futures", "date=2026-07-17&category=TOPIXF&contract_flag=1", `{"Code":"123","ProdCat":"TOPIXF","O":2700.5,"EO":"","OI":120,"LTD":"2026-09-10","CCMFlag":"1"}`, true, FuturesPrice{Code: "123", ProductCategory: "TOPIXF", WholeDayOpen: ptr(json.Number("2700.5")), OpenInterest: ptr(int64(120)), LastTradingDay: ptr("2026-09-10"), CentralContractMonthFlag: "1"}, func(c *Client) ([]FuturesPrice, error) {
		return collectChannel(func(ch chan<- FuturesPrice) error { return c.FuturesPriceWithChannel(t.Context(), req, ch) })
	})
}
