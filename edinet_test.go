package jquants

import (
	"encoding/json"
	"testing"
)

func TestClient_MajorShareholders(t *testing.T) {
	req := EdinetRequest{Code: ptr("7203")}
	checkEndpoint(t, "/edinet/major-shareholders", "code=7203", `{"Code":"7203","Hldrs":[{"Rank":1,"ShsHeld":140365.3}]}`, true, MajorShareholders{Code: "7203", Holders: []MajorShareholder{{Rank: 1, SharesHeld: 140365.3}}}, func(c *Client) ([]MajorShareholders, error) {
		return c.MajorShareholders(t.Context(), req)
	})
}

func TestClient_CrossShareholdings(t *testing.T) {
	req := EdinetRequest{EdinetCode: ptr("E00001"), Date: ptr("2026-07-17")}
	checkEndpoint(t, "/edinet/cross-shareholdings", "edinet_code=E00001&date=2026-07-17", `{"DocId":"doc-1","Largest":{"HldrName":"holder","Spec":[{"IsrName":"issuer","CurBookVal":12.5}]},"SecondLargest":null}`, true, CrossShareholdings{DocumentID: "doc-1", Largest: &CrossShareholdingEntry{HolderName: "holder", Specified: []CrossShareholdingIssue{{IssuerName: "issuer", CurrentBookValue: ptr(12.5)}}}}, func(c *Client) ([]CrossShareholdings, error) {
		return c.CrossShareholdings(t.Context(), req)
	})
}

func TestClient_LargeVolumeShareholders(t *testing.T) {
	req := EdinetRequest{Code: ptr("7203")}
	checkEndpoint(t, "/edinet/large-volume-shareholders", "code=7203", `{"Code":"7203","Hldrs":[{"AcqDisp":[{"Price":718.33}]}]}`, true, LargeVolumeShareholders{Code: "7203", Holders: []LargeVolumeHolder{{AcquisitionsDisposals: []LargeVolumeAcquisitionDisposal{{Price: 718.33}}}}}, func(c *Client) ([]LargeVolumeShareholders, error) {
		return c.LargeVolumeShareholders(t.Context(), req)
	})
}

// TestLargeVolumeAcquisitionDisposal_DecimalPrice pins the regression where the
// live EDINET API returned a decimal transaction price (e.g. 718.33). The field
// was previously int64, which aborted the entire fetch with a JSON unmarshal
// error. Money-denominated EDINET fields are now float64.
func TestLargeVolumeAcquisitionDisposal_DecimalPrice(t *testing.T) {
	const raw = `{"Date":"2024-01-15","SecType":"株式","Shs":1000,"Ratio":0.05,"Price":718.33}`
	var got LargeVolumeAcquisitionDisposal
	if err := json.Unmarshal([]byte(raw), &got); err != nil {
		t.Fatalf("failed to unmarshal decimal price: %s", err)
	}
	if got.Price != 718.33 {
		t.Errorf("Price = %v, want 718.33", got.Price)
	}
	if got.Shares != 1000 {
		t.Errorf("Shares = %v, want 1000", got.Shares)
	}
}

// TestMajorShareholder_DecimalSharesHeld pins the regression where the live
// EDINET API returned a fractional shares-held value (ShsHeld=140365.3, filings
// express holdings in thousands of shares). The field was previously int64,
// which aborted the entire fetch with a JSON unmarshal error. All EDINET
// numerics are now float64 by policy.
func TestMajorShareholder_DecimalSharesHeld(t *testing.T) {
	const raw = `{"Rank":1,"HldrName":"テスト","HldrAddr":"東京","ShsHeld":140365.3,"ShsRatio":0.12}`
	var got MajorShareholder
	if err := json.Unmarshal([]byte(raw), &got); err != nil {
		t.Fatalf("failed to unmarshal decimal shares held: %s", err)
	}
	if got.SharesHeld != 140365.3 {
		t.Errorf("SharesHeld = %v, want 140365.3", got.SharesHeld)
	}
}
