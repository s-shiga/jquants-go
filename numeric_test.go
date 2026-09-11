package jquants

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestNullableNumber(t *testing.T) {
	for _, tc := range []struct {
		name string
		wire string
		want *json.Number
	}{
		{name: "null", wire: "null"},
		{name: "empty", wire: `""`},
		{name: "dash", wire: `"-"`},
		{name: "asterisk", wire: `"*"`},
		{name: "number", wire: "12.50", want: ptr(json.Number("12.50"))},
		{name: "numeric string", wire: `"12.50"`, want: ptr(json.Number("12.50"))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got nullableNumber
			if err := json.Unmarshal([]byte(tc.wire), &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got.jsonNumber(), tc.want) {
				t.Fatalf("number = %v, want %v", got.jsonNumber(), tc.want)
			}
		})
	}
}

func TestNullableNumberRejectsInvalidValues(t *testing.T) {
	for _, wire := range []string{"\"not available\"", "\"NaN\"", "\"Inf\"", "true", "{}"} {
		t.Run(wire, func(t *testing.T) {
			var got nullableNumber
			if err := json.Unmarshal([]byte(wire), &got); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestNullableNumberIntegerConversion(t *testing.T) {
	for _, tc := range []struct {
		name string
		wire string
		want int64
	}{
		{name: "integer", wire: "9007199254740993", want: 9007199254740993},
		{name: "decimal notation", wire: "12500.0", want: 12500},
		{name: "exponent notation", wire: "1e3", want: 1000},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var number nullableNumber
			if err := json.Unmarshal([]byte(tc.wire), &number); err != nil {
				t.Fatal(err)
			}
			got, err := number.int64()
			if err != nil || got == nil || *got != tc.want {
				t.Fatalf("integer = %v, error = %v; want %d", got, err, tc.want)
			}
		})
	}

	for _, wire := range []string{"1.5", "9223372036854775808"} {
		var number nullableNumber
		if err := json.Unmarshal([]byte(wire), &number); err != nil {
			t.Fatal(err)
		}
		if _, err := number.int64(); err == nil {
			t.Fatalf("expected %s to be rejected as an int64", wire)
		}
	}
}

func TestNullableNumberInt32Range(t *testing.T) {
	var number nullableNumber
	if err := json.Unmarshal([]byte("2147483648"), &number); err != nil {
		t.Fatal(err)
	}
	if _, err := number.int32(); err == nil {
		t.Fatal("expected int32 overflow error")
	}
}

func TestStockPricePreservesLargeInteger(t *testing.T) {
	var got StockPrice
	err := json.Unmarshal([]byte(`{"UL":"0","LL":"0","Vo":9007199254740993,"Va":9007199254740995,"AdjFactor":1,"AdjVo":9007199254740997}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	if got.Volume == nil || *got.Volume != 9007199254740993 || got.TurnoverValue == nil || *got.TurnoverValue != 9007199254740995 || got.AdjustedVolume == nil || *got.AdjustedVolume != 9007199254740997 {
		t.Fatalf("large integer fields lost precision: %#v", got)
	}
}

func TestNumericPlaceholderErrorSurfacesFromRecord(t *testing.T) {
	var got Dividend
	err := json.Unmarshal([]byte(`{"DivRate":"not available"}`), &got)
	if err == nil {
		t.Fatal("expected invalid placeholder to fail decoding")
	}
	if !strings.Contains(err.Error(), "invalid numeric value") {
		t.Fatalf("expected numeric validation error, got %v", err)
	}
}
