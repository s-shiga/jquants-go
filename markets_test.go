package jquants

import (
	"testing"
)

func TestClient_MarginTradingOutstanding(t *testing.T) {
	req := MarginTradingOutstandingRequest{Code: ptr("86970")}
	checkEndpoint(t, "/markets/margin-interest", "code=86970", `{"Code":"86970","IssType":"2","ShrtVol":1200,"LongVol":3400}`, true, MarginTradingOutstanding{Code: "86970", IssueType: 2, TotalShortBalance: 1200, TotalLongBalance: 3400}, func(c *Client) ([]MarginTradingOutstanding, error) {
		return c.MarginTradingOutstanding(t.Context(), req)
	})
}

func TestClient_ShortSellingValue(t *testing.T) {
	req := ShortSellingValueRequest{Sector33Code: ptr("0050"), Date: ptr("2026-07-17")}
	checkEndpoint(t, "/markets/short-ratio", "s33=0050&date=2026-07-17", `{"S33":"0050","SellExShortVa":100,"ShrtWithResVa":200,"ShrtNoResVa":300}`, true, ShortSellingValue{Sector33Code: "0050", LongSellingValue: 100, ShortSellingWithRestrictions: 200, ShortSellingWithoutRestrictions: 300}, func(c *Client) ([]ShortSellingValue, error) {
		return c.ShortSellingValue(t.Context(), req)
	})
}

func TestClient_OutstandingShortPosition(t *testing.T) {
	req := OutstandingShortPositionRequest{CalculationDate: ptr("2026-07-17")}
	checkEndpoint(t, "/markets/short-sale-report", "calc_date=2026-07-17", `{"Code":"86970","CalcDate":"2026-07-17"}`, true, OutstandingShortPosition{Code: "86970", CalculationDate: "2026-07-17"}, func(c *Client) ([]OutstandingShortPosition, error) {
		return c.OutstandingShortPosition(t.Context(), req)
	})
}

func TestClient_MarginAlert(t *testing.T) {
	req := MarginAlertRequest{Date: ptr("2026-07-17")}
	checkEndpoint(t, "/markets/margin-alert", "date=2026-07-17", `{"Code":"86970","PubReason":{"Restricted":"1"},"ShrtOutChg":"-"}`, true, MarginAlert{Code: "86970", PublicationReason: MarginAlertPublicationReason{Restricted: "1"}}, func(c *Client) ([]MarginAlert, error) {
		return c.MarginAlert(t.Context(), req)
	})
}

func TestClient_BreakdownTrading(t *testing.T) {
	req := BreakdownTradingRequest{Code: ptr("86970")}
	checkEndpoint(t, "/markets/breakdown", "code=86970", `{"Code":"86970","LongSellVa":1234}`, true, BreakdownTrading{Code: "86970", LongSellValue: 1234}, func(c *Client) ([]BreakdownTrading, error) {
		return c.BreakdownTrading(t.Context(), req)
	})
}

func TestClient_BreakdownTradingWithChannel(t *testing.T) {
	req := BreakdownTradingRequest{Code: ptr("86970")}
	checkEndpoint(t, "/markets/breakdown", "code=86970", `{"Code":"86970","LongSellVa":1234}`, true, BreakdownTrading{Code: "86970", LongSellValue: 1234}, func(c *Client) ([]BreakdownTrading, error) {
		return collectChannel(func(ch chan<- BreakdownTrading) error { return c.BreakdownTradingWithChannel(t.Context(), req, ch) })
	})
}

func TestClient_TradingCalendar(t *testing.T) {
	req := TradingCalendarRequest{HolidayDivision: ptr(int8(1))}
	checkEndpoint(t, "/markets/calendar", "hol_div=1", `{"Date":"2026-07-17","HolDiv":"1"}`, false, TradingCalendar{Date: "2026-07-17", DayType: 1}, func(c *Client) ([]TradingCalendar, error) {
		return c.TradingCalendar(t.Context(), req)
	})
}
