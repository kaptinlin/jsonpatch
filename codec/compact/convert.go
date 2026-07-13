package compact

import "math"

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

func integerNumber(v any) (float64, error) {
	number, err := finiteNumber(v)
	if err != nil {
		return 0, err
	}
	if number != math.Trunc(number) {
		return 0, ErrNumberNotInteger
	}
	return number, nil
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
