package jquants

import "testing"

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
