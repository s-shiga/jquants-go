package jquants

import (
	"net/http"
	"testing"
	"time"
)

func setup() (*Client, error) {
	httpClient := &http.Client{Timeout: time.Second * 10}
	return NewClient(httpClient)
}

func TestNewClient(t *testing.T) {
	_, err := setup()
	if err != nil {
		t.Error(err)
	}
}
