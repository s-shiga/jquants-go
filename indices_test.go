package jquants

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestClient_IndexPrice(t *testing.T) {
	req := IndexPriceRequest{Code: ptr("0000")}
	checkEndpoint(t, "/indices/bars/daily", "code=0000", `{"Code":"0000","O":100.1,"H":102.2,"L":99.3,"C":101.4}`, true, IndexPrice{Code: "0000", Open: ptr(json.Number("100.1")), High: ptr(json.Number("102.2")), Low: ptr(json.Number("99.3")), Close: json.Number("101.4")}, func(c *Client) ([]IndexPrice, error) {
		return c.IndexPrice(t.Context(), req)
	})
}

// Indices that publish only a closing value return null for O, H, and L.
func TestIndexPriceCloseOnly(t *testing.T) {
	var got IndexPrice
	if err := json.Unmarshal([]byte(`{"Date":"2026-07-17","Code":"0500","O":null,"H":null,"L":null,"C":2450.5}`), &got); err != nil {
		t.Fatal(err)
	}
	want := IndexPrice{Date: "2026-07-17", Code: "0500", Close: json.Number("2450.5")}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("index price = %#v, want %#v", got, want)
	}
}

func TestClient_TopixPrices(t *testing.T) {
	req := TopixPriceRequest{From: ptr("2026-07-01"), To: ptr("2026-07-17")}
	checkEndpoint(t, "/indices/bars/daily/topix", "from=2026-07-01&to=2026-07-17", `{"O":100.1,"H":102.2,"L":99.3,"C":101.4}`, true, TopixPrice{Open: json.Number("100.1"), High: json.Number("102.2"), Low: json.Number("99.3"), Close: json.Number("101.4")}, func(c *Client) ([]TopixPrice, error) {
		return c.TopixPrices(t.Context(), req)
	})
}
