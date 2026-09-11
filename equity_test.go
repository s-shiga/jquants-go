package jquants

import (
	"encoding/json"
	"testing"
)

func TestClient_IssueInformation(t *testing.T) {
	req := IssueInformationRequest{Code: ptr("86970")}
	checkEndpoint(t, "/equities/master", "code=86970", `{"Code":"86970","CoName":"日本取引所","S17":"16","Mrgn":"2","ProdCat":"011"}`, false, IssueInformation{Code: "86970", CompanyName: "日本取引所", Sector17Code: 16, MarginCode: ptr(int8(2)), ProductCategory: "011"}, func(c *Client) ([]IssueInformation, error) {
		return c.IssueInformation(t.Context(), req)
	})
}

func TestClient_StockPrice(t *testing.T) {
	req := StockPriceRequest{Code: ptr("86970"), From: ptr("2026-07-01"), To: ptr("2026-07-17")}
	checkEndpoint(t, "/equities/bars/daily", "code=86970&from=2026-07-01&to=2026-07-17", `{"Code":"86970","O":2047.5,"C":null,"UL":"0","LL":"1","Vo":1200.0,"AdjFactor":0.5,"AdjO":1023.75,"AdjVo":2400.0}`, true, StockPrice{Code: "86970", Open: ptr(json.Number("2047.5")), LowerLimit: true, Volume: ptr(int64(1200)), AdjustmentFactor: json.Number("0.5"), AdjustedOpen: ptr(json.Number("1023.75")), AdjustedVolume: ptr(int64(2400))}, func(c *Client) ([]StockPrice, error) {
		return c.StockPrice(t.Context(), req)
	})
}

func TestClient_StockPriceWithChannel(t *testing.T) {
	req := StockPriceRequest{Code: ptr("86970"), From: ptr("2026-07-01"), To: ptr("2026-07-17")}
	checkEndpoint(t, "/equities/bars/daily", "code=86970&from=2026-07-01&to=2026-07-17", `{"Code":"86970","O":2047.5,"C":null,"UL":"0","LL":"1","Vo":1200.0,"AdjFactor":0.5,"AdjO":1023.75,"AdjVo":2400.0}`, true, StockPrice{Code: "86970", Open: ptr(json.Number("2047.5")), LowerLimit: true, Volume: ptr(int64(1200)), AdjustmentFactor: json.Number("0.5"), AdjustedOpen: ptr(json.Number("1023.75")), AdjustedVolume: ptr(int64(2400))}, func(c *Client) ([]StockPrice, error) {
		return collectChannel(func(ch chan<- StockPrice) error { return c.StockPriceWithChannel(t.Context(), req, ch) })
	})
}

func TestClient_MinuteStockPrice(t *testing.T) {
	req := MinuteStockPriceRequest{Code: ptr("86970")}
	checkEndpoint(t, "/equities/bars/minute", "code=86970", `{"Code":"86970","Time":"09:00","O":2047.5,"C":2050,"Vo":12500,"Va":25625000}`, true, MinuteStockPrice{Code: "86970", Time: "09:00", Open: ptr(json.Number("2047.5")), Close: ptr(json.Number("2050")), Volume: ptr(int64(12500)), TurnoverValue: ptr(int64(25625000))}, func(c *Client) ([]MinuteStockPrice, error) {
		return c.MinuteStockPrice(t.Context(), req)
	})
}

func TestClient_MinuteStockPriceWithChannel(t *testing.T) {
	req := MinuteStockPriceRequest{Code: ptr("86970")}
	checkEndpoint(t, "/equities/bars/minute", "code=86970", `{"Code":"86970","Time":"09:00","O":2047.5,"C":2050,"Vo":12500,"Va":25625000}`, true, MinuteStockPrice{Code: "86970", Time: "09:00", Open: ptr(json.Number("2047.5")), Close: ptr(json.Number("2050")), Volume: ptr(int64(12500)), TurnoverValue: ptr(int64(25625000))}, func(c *Client) ([]MinuteStockPrice, error) {
		return collectChannel(func(ch chan<- MinuteStockPrice) error { return c.MinuteStockPriceWithChannel(t.Context(), req, ch) })
	})
}

func TestClient_MorningSessionStockPrice(t *testing.T) {
	req := MorningSessionStockPriceRequest{Code: ptr("86970")}
	checkEndpoint(t, "/equities/bars/daily/am", "code=86970", `{"Code":"86970","MO":100.5,"MH":null,"MVo":100}`, true, MorningSessionStockPrice{Code: "86970", Open: ptr(json.Number("100.5")), Volume: ptr(int64(100))}, func(c *Client) ([]MorningSessionStockPrice, error) {
		return c.MorningSessionStockPrice(t.Context(), req)
	})
}

func TestClient_EarningsCalendar(t *testing.T) {
	req := EarningsCalendarRequest{}
	checkEndpoint(t, "/equities/earnings-calendar", "", `{"Code":"86970","CoName":"日本取引所","FY":"3月","FQ":"第1四半期"}`, true, EarningsCalendar{Code: "86970", CompanyName: "日本取引所", FiscalYear: "3月", FiscalQuarter: "第1四半期"}, func(c *Client) ([]EarningsCalendar, error) {
		return c.EarningsCalendar(t.Context(), req)
	})
}

func TestClient_InvestorType(t *testing.T) {
	req := InvestorTypeRequest{Section: ptr("TSEPrime")}
	checkEndpoint(t, "/equities/investor-types", "section=TSEPrime", `{"Section":"TSEPrime","PropSell":10,"PropBuy":20,"PropTot":30,"PropBal":10}`, true, InvestorType{Section: "TSEPrime", Proprietary: TradingBalance{Sales: 10, Purchases: 20, Total: 30, Balance: 10}}, func(c *Client) ([]InvestorType, error) {
		return c.InvestorType(t.Context(), req)
	})
}
