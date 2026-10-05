package jquants

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// Populate every exported field so newly added fields are also checked for
// loss when callers cache API results with encoding/json.
func populateJSONRecord(v reflect.Value, path string, populated bool) {
	if v.Type() == reflect.TypeFor[json.Number]() {
		value := "0"
		if populated {
			value = "9007199254740993.1250"
		}
		v.SetString(value)
		return
	}
	if !populated && v.Kind() != reflect.Struct {
		return
	}
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			populateJSONRecord(v.Field(i), path+"."+v.Type().Field(i).Name, populated)
		}
	case reflect.Pointer:
		v.Set(reflect.New(v.Type().Elem()))
		populateJSONRecord(v.Elem(), path, populated)
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), 1, 1))
		populateJSONRecord(v.Index(0), path+"[0]", populated)
	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
		key, value := reflect.New(v.Type().Key()).Elem(), reflect.New(v.Type().Elem()).Elem()
		populateJSONRecord(key, path+".key", populated)
		populateJSONRecord(value, path+".value", populated)
		v.SetMapIndex(key, value)
	case reflect.String:
		v.SetString(path + ` "日本語"`)
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int64:
		v.SetInt(9007199254740993)
	case reflect.Int, reflect.Int8, reflect.Int32:
		v.SetInt(12)
	case reflect.Float64:
		v.SetFloat(12.375)
	default:
		panic(fmt.Sprintf("unhandled record field %s (%s)", path, v.Type()))
	}
}

func TestRecordsJSONRoundTrip(t *testing.T) {
	for _, record := range []any{
		IssueInformation{}, StockPrice{}, MinuteStockPrice{}, MorningSessionStockPrice{},
		EarningsCalendar{}, InvestorType{}, MarginTradingOutstanding{}, ShortSellingValue{},
		BreakdownTrading{}, OutstandingShortPosition{}, MarginAlert{}, TradingCalendar{},
		IndexPrice{}, TopixPrice{}, FuturesPrice{}, IndexOptionPrice{}, OptionPrice{},
		FinancialSummary{}, FinancialDetails{}, Dividend{}, BulkFile{}, TimelyDisclosure{},
		MajorShareholders{}, CrossShareholdings{}, LargeVolumeShareholders{},
		TimelyDisclosureFiles{}, TimelyDisclosureBulk{},
	} {
		recordType := reflect.TypeOf(record)
		for _, populated := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s/populated=%t", recordType.Name(), populated), func(t *testing.T) {
				want := reflect.New(recordType)
				populateJSONRecord(want.Elem(), recordType.Name(), populated)
				data, err := json.Marshal(want.Interface())
				if err != nil {
					t.Fatal(err)
				}
				got := reflect.New(recordType)
				if err := json.Unmarshal(data, got.Interface()); err != nil {
					t.Fatalf("reload serialized record: %v", err)
				}
				if !reflect.DeepEqual(got.Interface(), want.Interface()) {
					t.Fatalf("record changed after JSON round trip:\n got: %#v\nwant: %#v", got.Elem().Interface(), want.Elem().Interface())
				}
			})
		}
	}
}

func TestRecordsJSONRoundTripPreservesNullableValues(t *testing.T) {
	want := StockPrice{
		Code: "86970", AdjustmentFactor: json.Number("1.0"),
		MorningUpperLimit: ptr(false), AfternoonLowerLimit: ptr(true),
	}
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got StockPrice
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stock price = %#v, want %#v", got, want)
	}

	// A cached record with no revision must also clear a reused destination.
	wantDisclosure := TimelyDisclosure{DisclosureNumber: "20260717000001", Documents: []string{}}
	data, err = json.Marshal(wantDisclosure)
	if err != nil {
		t.Fatal(err)
	}
	gotDisclosure := TimelyDisclosure{RevisionNumber: "2", Documents: []string{"g"}}
	if err := json.Unmarshal(data, &gotDisclosure); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotDisclosure, wantDisclosure) {
		t.Fatalf("disclosure = %#v, want %#v", gotDisclosure, wantDisclosure)
	}
}
