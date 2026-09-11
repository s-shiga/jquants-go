package jquants

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestCodeDateRangeValues(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		code, date, from, to *string
		want                 string
		invalid              bool
	}{
		{name: "missing filters", invalid: true},
		{name: "code", code: ptr("86970"), want: "code=86970"},
		{name: "date", date: ptr("2026-07-17"), want: "date=2026-07-17"},
		{name: "code and date", code: ptr("86970"), date: ptr("2026-07-17"), want: "code=86970&date=2026-07-17"},
		{name: "range", code: ptr("86970"), from: ptr("2026-07-01"), to: ptr("2026-07-17"), want: "code=86970&from=2026-07-01&to=2026-07-17"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := codeDateRangeValues(tc.code, tc.date, tc.from, tc.to, ptr("next+/="))
			if tc.invalid {
				if err == nil {
					t.Fatal("expected validation error")
				}
				return
			}
			want, parseErr := url.ParseQuery(tc.want + "&pagination_key=next%2B%2F%3D")
			if parseErr != nil {
				t.Fatal(parseErr)
			}
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("query = %v, error = %v; want %v", got, err, want)
			}
		})
	}
}

func TestStockPrice_CodeAndDate(t *testing.T) {
	for _, stream := range []bool{false, true} {
		t.Run(fmt.Sprint(stream), func(t *testing.T) {
			req := StockPriceRequest{Code: ptr("86970"), Date: ptr("2026-07-17")}
			checkEndpoint(t, "/equities/bars/daily", "code=86970&date=2026-07-17", `{"Code":"86970","UL":"0","LL":"0"}`, true, StockPrice{Code: "86970"}, func(c *Client) ([]StockPrice, error) {
				if stream {
					return collectChannel(func(ch chan<- StockPrice) error { return c.StockPriceWithChannel(t.Context(), req, ch) })
				}
				return c.StockPrice(t.Context(), req)
			})
		})
	}
}
