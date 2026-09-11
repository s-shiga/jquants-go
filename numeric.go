package jquants

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

// nullableNumber decodes the API's numeric wire format without first routing
// integers through float64. Null and the documented empty placeholders mean no
// value; malformed and unexpected strings remain errors.
type nullableNumber struct {
	value *json.Number
}

func (n *nullableNumber) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		n.value = nil
		return nil
	}
	if len(data) > 0 && data[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		return n.set(value)
	}
	return n.set(string(data))
}

func (n *nullableNumber) set(value string) error {
	switch value {
	case "", "-", "*":
		n.value = nil
		return nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(parsed) || math.IsInf(parsed, 0) {
		if err == nil {
			err = fmt.Errorf("value is not finite")
		}
		return fmt.Errorf("invalid numeric value %q: %w", value, err)
	}
	number := json.Number(value)
	n.value = &number
	return nil
}

func (n nullableNumber) jsonNumber() *json.Number {
	if n.value == nil {
		return nil
	}
	value := *n.value
	return &value
}

func (n nullableNumber) float64() (*float64, error) {
	if n.value == nil {
		return nil, nil
	}
	value, err := n.value.Float64()
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (n nullableNumber) int64() (*int64, error) {
	if n.value == nil {
		return nil, nil
	}
	value, ok := new(big.Rat).SetString(n.value.String())
	if !ok || !value.IsInt() || !value.Num().IsInt64() {
		return nil, fmt.Errorf("numeric value %q is not an int64", n.value.String())
	}
	result := value.Num().Int64()
	return &result, nil
}

func (n nullableNumber) int32() (*int32, error) {
	value, err := n.int64()
	if err != nil || value == nil {
		return nil, err
	}
	if *value < math.MinInt32 || *value > math.MaxInt32 {
		return nil, fmt.Errorf("numeric value %q is not an int32", n.value.String())
	}
	result := int32(*value)
	return &result, nil
}

// unmarshaler accumulates conversion errors while assigning related numeric
// fields in a custom JSON unmarshaler.
type unmarshaler struct {
	err error
}

func (u *unmarshaler) price(v nullableNumber) *int32 {
	if u.err != nil {
		return nil
	}
	result, err := v.int32()
	u.err = err
	return result
}

func (u *unmarshaler) volume(v nullableNumber) *int64 {
	if u.err != nil {
		return nil
	}
	result, err := v.int64()
	u.err = err
	return result
}

func (u *unmarshaler) integer(v nullableNumber) int64 {
	result := u.volume(v)
	if result == nil {
		return 0
	}
	return *result
}

func (u *unmarshaler) jsonNumber(v nullableNumber) *json.Number {
	if u.err != nil {
		return nil
	}
	return v.jsonNumber()
}
