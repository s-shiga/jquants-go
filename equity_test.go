package jquants

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestClient_IssueInformation(t *testing.T) {
	req := IssueInformationRequest{Code: ptr("86970")}
	checkEndpoint(t, "/equities/master", "code=86970", `{"Code":"86970","CoName":"日本取引所","S17":"16","Mrgn":"2","ProdCat":"011"}`, true, IssueInformation{Code: "86970", CompanyName: "日本取引所", Sector17Code: 16, MarginCode: ptr(int8(2)), ProductCategory: "011"}, func(c *Client) ([]IssueInformation, error) {
		return c.IssueInformation(t.Context(), req)
	})
}

func TestIssueInformationClearsMissingMarginCode(t *testing.T) {
	value := IssueInformation{MarginCode: ptr(int8(2))}
	if err := json.Unmarshal([]byte(`{"S17":"16"}`), &value); err != nil {
		t.Fatal(err)
	}
	if value.MarginCode != nil {
		t.Fatalf("MarginCode = %v, want nil", *value.MarginCode)
	}
}

func TestClient_StockPrice(t *testing.T) {
	req := StockPriceRequest{Code: ptr("86970"), From: ptr("2026-07-01"), To: ptr("2026-07-17")}
	checkEndpoint(t, "/equities/bars/daily", "code=86970&from=2026-07-01&to=2026-07-17", `{"Code":"86970","O":2047.5,"C":null,"UL":"0","LL":"1","Vo":1200.0,"AdjFactor":0.5,"AdjO":1023.75,"AdjVo":2400.5}`, true, StockPrice{Code: "86970", Open: ptr(json.Number("2047.5")), LowerLimit: true, Volume: ptr(int64(1200)), AdjustmentFactor: json.Number("0.5"), AdjustedOpen: ptr(json.Number("1023.75")), AdjustedVolume: ptr(json.Number("2400.5"))}, func(c *Client) ([]StockPrice, error) {
		return c.StockPrice(t.Context(), req)
	})
}

func TestClient_StockPriceWithChannel(t *testing.T) {
	req := StockPriceRequest{Code: ptr("86970"), From: ptr("2026-07-01"), To: ptr("2026-07-17")}
	checkEndpoint(t, "/equities/bars/daily", "code=86970&from=2026-07-01&to=2026-07-17", `{"Code":"86970","O":2047.5,"C":null,"UL":"0","LL":"1","Vo":1200.0,"AdjFactor":0.5,"AdjO":1023.75,"AdjVo":2400.5}`, true, StockPrice{Code: "86970", Open: ptr(json.Number("2047.5")), LowerLimit: true, Volume: ptr(int64(1200)), AdjustmentFactor: json.Number("0.5"), AdjustedOpen: ptr(json.Number("1023.75")), AdjustedVolume: ptr(json.Number("2400.5"))}, func(c *Client) ([]StockPrice, error) {
		return collectChannel(func(ch chan<- StockPrice) error { return c.StockPriceWithChannel(t.Context(), req, ch) })
	})
}

func TestStockPrice_SessionFields(t *testing.T) {
	wire := `{"UL":"0","LL":"0","MO":101.1,"MH":102.2,"ML":99.9,"MC":100.5,"MUL":"1","MLL":"0","MVo":111,"MVa":222,"MAdjO":50.55,"MAdjH":51.1,"MAdjL":49.95,"MAdjC":50.25,"MAdjVo":333.3,"AO":103.3,"AH":104.4,"AL":98.8,"AC":102.5,"AUL":"0","ALL":"1","AVo":444,"AVa":555,"AAdjO":51.65,"AAdjH":52.2,"AAdjL":49.4,"AAdjC":51.25,"AAdjVo":666.7}`
	want := StockPrice{
		MorningOpen:             ptr(json.Number("101.1")),
		MorningHigh:             ptr(json.Number("102.2")),
		MorningLow:              ptr(json.Number("99.9")),
		MorningClose:            ptr(json.Number("100.5")),
		MorningUpperLimit:       ptr(true),
		MorningLowerLimit:       ptr(false),
		MorningVolume:           ptr(int64(111)),
		MorningTurnoverValue:    ptr(int64(222)),
		MorningAdjustedOpen:     ptr(json.Number("50.55")),
		MorningAdjustedHigh:     ptr(json.Number("51.1")),
		MorningAdjustedLow:      ptr(json.Number("49.95")),
		MorningAdjustedClose:    ptr(json.Number("50.25")),
		MorningAdjustedVolume:   ptr(json.Number("333.3")),
		AfternoonOpen:           ptr(json.Number("103.3")),
		AfternoonHigh:           ptr(json.Number("104.4")),
		AfternoonLow:            ptr(json.Number("98.8")),
		AfternoonClose:          ptr(json.Number("102.5")),
		AfternoonUpperLimit:     ptr(false),
		AfternoonLowerLimit:     ptr(true),
		AfternoonVolume:         ptr(int64(444)),
		AfternoonTurnoverValue:  ptr(int64(555)),
		AfternoonAdjustedOpen:   ptr(json.Number("51.65")),
		AfternoonAdjustedHigh:   ptr(json.Number("52.2")),
		AfternoonAdjustedLow:    ptr(json.Number("49.4")),
		AfternoonAdjustedClose:  ptr(json.Number("51.25")),
		AfternoonAdjustedVolume: ptr(json.Number("666.7")),
	}
	var got StockPrice
	if err := json.Unmarshal([]byte(wire), &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stock price = %#v, want %#v", got, want)
	}
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

func TestInvestorTypePreservesIntegerPrecision(t *testing.T) {
	var value InvestorType
	if err := json.Unmarshal([]byte(`{"PropSell":9007199254740993}`), &value); err != nil {
		t.Fatal(err)
	}
	if value.Proprietary.Sales != 9007199254740993 {
		t.Fatalf("Proprietary.Sales = %d", value.Proprietary.Sales)
	}
}

func TestInvestorTypeRejectsFractionalInteger(t *testing.T) {
	var value InvestorType
	if err := json.Unmarshal([]byte(`{"PropSell":1.5}`), &value); err == nil {
		t.Fatal("json.Unmarshal accepted a fractional trading balance")
	}
}

func TestStockPrice_ExRightsType(t *testing.T) {
	for _, tc := range []struct {
		name, data string
		want       *string
	}{
		{"rights issue", `{"UL":"0","LL":"0","ExRT":"3"}`, ptr("3")},
		{"no corporate action", `{"UL":"0","LL":"0","ExRT":null}`, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got StockPrice
			if err := json.Unmarshal([]byte(tc.data), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got.ExRightsType, tc.want) {
				t.Fatalf("ExRightsType = %v, want %v", got.ExRightsType, tc.want)
			}
		})
	}
}
