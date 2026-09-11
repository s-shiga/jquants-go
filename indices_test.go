package jquants

import (
	"encoding/json"
	"testing"
)

func TestClient_IndexPrice(t *testing.T) {
	req := IndexPriceRequest{Code: ptr("0000")}
	checkEndpoint(t, "/indices/bars/daily", "code=0000", `{"Code":"0000","O":100.1,"H":102.2,"L":99.3,"C":101.4}`, true, IndexPrice{Code: "0000", Open: json.Number("100.1"), High: json.Number("102.2"), Low: json.Number("99.3"), Close: json.Number("101.4")}, func(c *Client) ([]IndexPrice, error) {
		return c.IndexPrice(t.Context(), req)
	})
}

func TestClient_TopixPrices(t *testing.T) {
	req := TopixPriceRequest{From: ptr("2026-07-01"), To: ptr("2026-07-17")}
	checkEndpoint(t, "/indices/bars/daily/topix", "from=2026-07-01&to=2026-07-17", `{"O":100.1,"H":102.2,"L":99.3,"C":101.4}`, true, TopixPrice{Open: json.Number("100.1"), High: json.Number("102.2"), Low: json.Number("99.3"), Close: json.Number("101.4")}, func(c *Client) ([]TopixPrice, error) {
		return c.TopixPrices(t.Context(), req)
	})
}
