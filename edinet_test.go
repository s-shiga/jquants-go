package jquants

import (
	"encoding/json"
	"testing"
)

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

func TestClient_MajorShareholders(t *testing.T) {
	code := "7203"
	client := setupClient(t)
	res, err := client.MajorShareholders(t.Context(), EdinetRequest{Code: &code})
	if err != nil {
		t.Errorf("Failed to get major shareholders: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty major shareholders")
	}
}

func TestClient_CrossShareholdings(t *testing.T) {
	code := "7203"
	client := setupClient(t)
	res, err := client.CrossShareholdings(t.Context(), EdinetRequest{Code: &code})
	if err != nil {
		t.Errorf("Failed to get cross shareholdings: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty cross shareholdings")
	}
}

func TestClient_CrossShareholdingsWithLargest(t *testing.T) {
	code := "8306"
	client := setupClient(t)
	res, err := client.CrossShareholdings(t.Context(), EdinetRequest{Code: &code})
	if err != nil {
		t.Errorf("Failed to get cross shareholdings: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty cross shareholdings")
	}
	foundLargest := false
	for _, r := range res {
		if r.Largest != nil {
			foundLargest = true
			break
		}
	}
	if !foundLargest {
		t.Error("Expected at least one record with a non-null Largest holder")
	}
}

func TestClient_LargeVolumeShareholders(t *testing.T) {
	code := "7203"
	client := setupClient(t)
	res, err := client.LargeVolumeShareholders(t.Context(), EdinetRequest{Code: &code})
	if err != nil {
		t.Errorf("Failed to get large volume shareholders: %s", err)
	}
	if len(res) == 0 {
		t.Error("Empty large volume shareholders")
	}
}
