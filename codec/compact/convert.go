package compact

import (
	"math"
	"math/big"
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
)

func trueMarkerAt(raw Op, index int) (bool, error) {
	if len(raw) <= index {
		return false, nil
	}
	switch marker := raw[index].(type) {
	case int:
		if marker == 1 {
			return true, nil
		}
	case int64:
		if marker == 1 {
			return true, nil
		}
	case float64:
		if marker == 1 {
			return true, nil
		}
	}
	return false, ErrInvalidTrueMarker
}

// toFloat64 converts a value to float64.
func toFloat64(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	default:
		return 0, ErrNotFloat64
	}
}

func finiteNumber(v any) (float64, error) {
	number, err := toFloat64(v)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, ErrNumberNotFinite
	}
	return number, nil
}

func integerNumber(v any) (int, error) {
	return exactInteger(v)
}

func exactInteger(v any) (int, error) {
	switch value := v.(type) {
	case int:
		return value, nil
	case int8:
		return int(value), nil
	case int16:
		return int(value), nil
	case int32:
		return int(value), nil
	case int64:
		if strconv.IntSize == 32 && (value < math.MinInt32 || value > math.MaxInt32) {
			return 0, ErrNumberOutOfRange
		}
		return int(value), nil
	case uint:
		if uint64(value) > uint64(^uint(0)>>1) {
			return 0, ErrNumberOutOfRange
		}
		return int(value), nil
	case uint8:
		return int(value), nil
	case uint16:
		return int(value), nil
	case uint32:
		if strconv.IntSize == 32 && uint64(value) > math.MaxInt32 {
			return 0, ErrNumberOutOfRange
		}
		return int(value), nil
	case uint64:
		if value > uint64(^uint(0)>>1) {
			return 0, ErrNumberOutOfRange
		}
		return int(value), nil
	case float32:
		return integerFromFloat(float64(value))
	case float64:
		return integerFromFloat(value)
	case jsontext.Value:
		return integerFromJSON(value)
	default:
		return 0, ErrNotFloat64
	}
}

func integerFromJSON(raw jsontext.Value) (int, error) {
	text := string(raw)
	if value, err := strconv.ParseInt(text, 10, strconv.IntSize); err == nil {
		return int(value), nil
	}
	number, ok := new(big.Rat).SetString(text)
	if !ok || !number.IsInt() {
		return 0, ErrNumberNotInteger
	}
	integer := number.Num()
	if !integer.IsInt64() {
		return 0, ErrNumberOutOfRange
	}
	value := integer.Int64()
	if strconv.IntSize == 32 && (value < math.MinInt32 || value > math.MaxInt32) {
		return 0, ErrNumberOutOfRange
	}
	return int(value), nil
}

func integerFromFloat(number float64) (int, error) {
	if math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, ErrNumberNotFinite
	}
	if number != math.Trunc(number) {
		return 0, ErrNumberNotInteger
	}
	limit := math.Ldexp(1, strconv.IntSize-1)
	if number < -limit || number >= limit {
		return 0, ErrNumberOutOfRange
	}
	return int(number), nil
}

func isZeroMarker(v any) bool {
	switch marker := v.(type) {
	case int:
		return marker == 0
	case int64:
		return marker == 0
	case float64:
		return marker == 0
	default:
		return false
	}
}

// toStringSlice converts a value to []string.
func toStringSlice(v any) ([]string, error) {
	arr, ok := v.([]any)
	if !ok {
		return nil, ErrExpectedArray
	}
	result := make([]string, len(arr))
	for i, item := range arr {
		s, ok := item.(string)
		if !ok {
			return nil, ErrExpectedString
		}
		result[i] = s
	}
	return result, nil
}
