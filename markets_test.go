package jquants

import (
	"encoding/json"
	"reflect"
	"testing"
)

const marginTradingOutstandingDailyJSON = `{
	"PubDate":"2026-09-28","Date":"2026-09-25","Code":"86970","IssType":"2",
	"ShrtVol":1200.0,"LongVol":3400.0,"ShrtNegVol":1000.0,"LongNegVol":3000.0,
	"ShrtStdVol":200.0,"LongStdVol":400.0,
	"ShrtVal":120000.5,"LongVal":340000.0,"ShrtNegVal":100000.25,"LongNegVal":300000.0,
	"ShrtStdVal":20000.25,"LongStdVal":40000.0
}`

func TestClient_MarginTradingOutstanding(t *testing.T) {
	want := MarginTradingOutstanding{
		PublicationDate: ptr("2026-09-28"), Date: "2026-09-25", Code: "86970", IssueType: 2,
		TotalShortBalance: 1200, TotalLongBalance: 3400,
		ShortNegotiableBalance: 1000, LongNegotiableBalance: 3000,
		ShortStandardizedBalance: 200, LongStandardizedBalance: 400,
		TotalShortValue: ptr(120000.5), TotalLongValue: ptr(340000.0),
		ShortNegotiableValue: ptr(100000.25), LongNegotiableValue: ptr(300000.0),
		ShortStandardizedValue: ptr(20000.25), LongStandardizedValue: ptr(40000.0),
	}
	for _, tc := range []struct {
		name  string
		req   MarginTradingOutstandingRequest
		query string
	}{
		{"code", MarginTradingOutstandingRequest{Code: ptr("86970")}, "code=86970"},
		{"code and date", MarginTradingOutstandingRequest{Code: ptr("86970"), Date: ptr("2026-09-25")}, "code=86970&date=2026-09-25"},
		{"code and range", MarginTradingOutstandingRequest{Code: ptr("86970"), From: ptr("2026-09-25"), To: ptr("2026-09-30")}, "code=86970&from=2026-09-25&to=2026-09-30"},
		{"date", MarginTradingOutstandingRequest{Date: ptr("2026-09-25")}, "date=2026-09-25"},
		{"published date", MarginTradingOutstandingRequest{PublishedDate: ptr("2026-09-28")}, "published_date=2026-09-28"},
		{"code and published date", MarginTradingOutstandingRequest{Code: ptr("86970"), PublishedDate: ptr("20260928")}, "code=86970&published_date=20260928"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checkEndpoint(t, "/markets/margin-interest", tc.query, marginTradingOutstandingDailyJSON, true, want, func(c *Client) ([]MarginTradingOutstanding, error) {
				return c.MarginTradingOutstanding(t.Context(), tc.req)
			})
		})
	}
}

func TestClient_MarginTradingOutstandingRejectsInvalidFilters(t *testing.T) {
	for _, tc := range []struct {
		name string
		req  MarginTradingOutstandingRequest
	}{
		{"missing selector", MarginTradingOutstandingRequest{}},
		{"range without selector", MarginTradingOutstandingRequest{From: ptr("2026-09-25"), To: ptr("2026-09-30")}},
		{"published date and date", MarginTradingOutstandingRequest{PublishedDate: ptr("2026-09-28"), Date: ptr("2026-09-25")}},
		{"published date and from", MarginTradingOutstandingRequest{PublishedDate: ptr("2026-09-28"), From: ptr("2026-09-25")}},
		{"published date and to", MarginTradingOutstandingRequest{PublishedDate: ptr("2026-09-28"), To: ptr("2026-09-30")}},
		{"code and conflicting dates", MarginTradingOutstandingRequest{Code: ptr("86970"), PublishedDate: ptr("2026-09-28"), Date: ptr("2026-09-25")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// No fixture responses: invalid filters must fail before any HTTP request.
			_, err := fixtureClient(t).MarginTradingOutstanding(t.Context(), tc.req)
			if err == nil {
				t.Fatal("expected invalid filter error")
			}
		})
	}
}

func TestMarginTradingOutstandingHistoricalRecords(t *testing.T) {
	for _, tc := range []struct {
		name string
		data string
	}{
		{"null fields", `{"Date":"2026-09-18","Code":"86970","IssType":"2","ShrtVol":1200,"LongVol":3400,"PubDate":null,"ShrtVal":null,"LongVal":null,"ShrtNegVal":null,"LongNegVal":null,"ShrtStdVal":null,"LongStdVal":null}`},
		{"omitted fields", `{"Date":"2026-09-18","Code":"86970","IssType":"2","ShrtVol":1200,"LongVol":3400}`},
		{"placeholder fields", `{"Date":"2026-09-18","Code":"86970","IssType":"2","ShrtVol":1200,"LongVol":3400,"ShrtVal":"","LongVal":"-","ShrtNegVal":"*","LongNegVal":"","ShrtStdVal":"-","LongStdVal":"*"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var value MarginTradingOutstanding
			if err := json.Unmarshal([]byte(marginTradingOutstandingDailyJSON), &value); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tc.data), &value); err != nil {
				t.Fatal(err)
			}
			want := MarginTradingOutstanding{Date: "2026-09-18", Code: "86970", IssueType: 2, TotalShortBalance: 1200, TotalLongBalance: 3400}
			if !reflect.DeepEqual(value, want) {
				t.Fatalf("historical record = %#v, want %#v", value, want)
			}
		})
	}
}

func TestMarginTradingOutstandingZeroValues(t *testing.T) {
	var value MarginTradingOutstanding
	if err := json.Unmarshal([]byte(`{"IssType":"3","ShrtVal":0,"LongVal":0,"ShrtNegVal":0,"LongNegVal":0,"ShrtStdVal":0,"LongStdVal":0}`), &value); err != nil {
		t.Fatal(err)
	}
	want := MarginTradingOutstanding{
		IssueType: 3, TotalShortValue: ptr(0.0), TotalLongValue: ptr(0.0),
		ShortNegotiableValue: ptr(0.0), LongNegotiableValue: ptr(0.0),
		ShortStandardizedValue: ptr(0.0), LongStandardizedValue: ptr(0.0),
	}
	if !reflect.DeepEqual(value, want) {
		t.Fatalf("zero values = %#v, want %#v", value, want)
	}
}

func TestMarginTradingOutstandingStringValues(t *testing.T) {
	var value MarginTradingOutstanding
	if err := json.Unmarshal([]byte(`{"IssType":"1","ShrtVal":"514800000.0","LongVal":"450000000","ShrtNegVal":"485600000","LongNegVal":"163800000","ShrtStdVal":"29200000","LongStdVal":"286200000"}`), &value); err != nil {
		t.Fatal(err)
	}
	want := MarginTradingOutstanding{
		IssueType: 1, TotalShortValue: ptr(514800000.0), TotalLongValue: ptr(450000000.0),
		ShortNegotiableValue: ptr(485600000.0), LongNegotiableValue: ptr(163800000.0),
		ShortStandardizedValue: ptr(29200000.0), LongStandardizedValue: ptr(286200000.0),
	}
	if !reflect.DeepEqual(value, want) {
		t.Fatalf("string values = %#v, want %#v", value, want)
	}
}

func TestMarginTradingOutstandingRejectsInvalidValue(t *testing.T) {
	var value MarginTradingOutstanding
	if err := json.Unmarshal([]byte(`{"IssType":"1","ShrtVal":"abc"}`), &value); err == nil {
		t.Fatal("expected invalid value error")
	}
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
	checkEndpoint(t, "/markets/calendar", "hol_div=1", `{"Date":"2026-07-17","HolDiv":"1"}`, true, TradingCalendar{Date: "2026-07-17", DayType: 1}, func(c *Client) ([]TradingCalendar, error) {
		return c.TradingCalendar(t.Context(), req)
	})
}

func TestMarketIntegerFieldsPreservePrecision(t *testing.T) {
	const exact = int64(9007199254740993)
	for _, tc := range []struct {
		name string
		data string
		got  func() (int64, error)
	}{
		{
			name: "margin interest",
			data: `{"IssType":"1","ShrtVol":9007199254740993}`,
			got: func() (int64, error) {
				var value MarginTradingOutstanding
				err := json.Unmarshal([]byte(`{"IssType":"1","ShrtVol":9007199254740993}`), &value)
				return value.TotalShortBalance, err
			},
		},
		{
			name: "short selling",
			data: `{"SellExShortVa":9007199254740993}`,
			got: func() (int64, error) {
				var value ShortSellingValue
				err := json.Unmarshal([]byte(`{"SellExShortVa":9007199254740993}`), &value)
				return value.LongSellingValue, err
			},
		},
		{
			name: "breakdown volume",
			data: `{"LongSellVo":9007199254740993}`,
			got: func() (int64, error) {
				var value BreakdownTrading
				err := json.Unmarshal([]byte(`{"LongSellVo":9007199254740993}`), &value)
				return value.LongSellVolume, err
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.got()
			if err != nil || got != exact {
				t.Fatalf("decoded integer = %d, %v; want %d", got, err, exact)
			}
		})
	}
}

func TestMarketIntegerFieldsRejectFractions(t *testing.T) {
	for _, target := range []any{&MarginTradingOutstanding{}, &ShortSellingValue{}, &BreakdownTrading{}} {
		if err := json.Unmarshal([]byte(`{"IssType":"1","ShrtVol":1.5,"SellExShortVa":1.5,"LongSellVo":1.5}`), target); err == nil {
			t.Fatalf("json.Unmarshal into %T accepted a fractional integer", target)
		}
	}
}

func TestShortSellingValueParameters(t *testing.T) {
	p := func(req ShortSellingValueRequest) parameters {
		return shortSellingValueParameters{ShortSellingValueRequest: req}
	}
	checkValues(t, []valuesCase{
		{p(ShortSellingValueRequest{Date: ptr("2026-07-17")}), "date=2026-07-17"},
		{p(ShortSellingValueRequest{Sector33Code: ptr("0050")}), "s33=0050"},
		{p(ShortSellingValueRequest{Sector33Code: ptr("0050"), Date: ptr("2026-07-17")}), "date=2026-07-17&s33=0050"},
		{p(ShortSellingValueRequest{Sector33Code: ptr("0050"), From: ptr("2026-07-01"), To: ptr("2026-07-17")}), "from=2026-07-01&s33=0050&to=2026-07-17"},
		{p(ShortSellingValueRequest{}), rejected},
		{p(ShortSellingValueRequest{From: ptr("2026-07-01"), To: ptr("2026-07-17")}), rejected},
		{p(ShortSellingValueRequest{Sector33Code: ptr("0050"), Date: ptr("2026-07-17"), From: ptr("2026-07-01")}), rejected},
	})
}

func TestOutstandingShortPositionParameters(t *testing.T) {
	p := func(req OutstandingShortPositionRequest) parameters {
		return outstandingShortPositionParameters{OutstandingShortPositionRequest: req}
	}
	code, day := ptr("86970"), ptr("2026-07-17")
	checkValues(t, []valuesCase{
		// The six combinations the spec allows.
		{p(OutstandingShortPositionRequest{Code: code}), "code=86970"},
		{p(OutstandingShortPositionRequest{Code: code, DisclosureDate: day}), "code=86970&disc_date=2026-07-17"},
		{p(OutstandingShortPositionRequest{Code: code, DisclosureDateFrom: ptr("2026-07-01"), DisclosureDateTo: day}), "code=86970&disc_date_from=2026-07-01&disc_date_to=2026-07-17"},
		{p(OutstandingShortPositionRequest{Code: code, CalculationDate: day}), "calc_date=2026-07-17&code=86970"},
		{p(OutstandingShortPositionRequest{DisclosureDate: day}), "disc_date=2026-07-17"},
		{p(OutstandingShortPositionRequest{CalculationDate: day}), "calc_date=2026-07-17"},
		{p(OutstandingShortPositionRequest{}), rejected},
		{p(OutstandingShortPositionRequest{DisclosureDateFrom: ptr("2026-07-01"), DisclosureDateTo: day}), rejected},
		{p(OutstandingShortPositionRequest{DisclosureDate: day, CalculationDate: day}), rejected},
		{p(OutstandingShortPositionRequest{Code: code, DisclosureDate: day, DisclosureDateTo: day}), rejected},
		{p(OutstandingShortPositionRequest{Code: code, DisclosureDateFrom: day, CalculationDate: day}), rejected},
	})
}

func TestMarginAlertParameters(t *testing.T) {
	p := func(req MarginAlertRequest) parameters { return marginAlertParameters{MarginAlertRequest: req} }
	checkValues(t, []valuesCase{
		{p(MarginAlertRequest{Code: ptr("86970")}), "code=86970"},
		{p(MarginAlertRequest{Code: ptr("86970"), Date: ptr("2026-07-17")}), "code=86970&date=2026-07-17"},
		{p(MarginAlertRequest{Code: ptr("86970"), From: ptr("2026-07-01"), To: ptr("2026-07-17")}), "code=86970&from=2026-07-01&to=2026-07-17"},
		{p(MarginAlertRequest{Date: ptr("2026-07-17")}), "date=2026-07-17"},
		{p(MarginAlertRequest{}), rejected},
		{p(MarginAlertRequest{From: ptr("2026-07-01"), To: ptr("2026-07-17")}), rejected},
		{p(MarginAlertRequest{Date: ptr("2026-07-17"), From: ptr("2026-07-01")}), rejected},
	})
}
