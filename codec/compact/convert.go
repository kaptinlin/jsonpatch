package compact

// boolAt safely extracts a bool value at the given index, returning false if absent.
func boolAt(raw Op, index int) bool {
	if len(raw) <= index {
		return false
	}
	return toBool(raw[index])
}

// float64At safely extracts a float64 value at the given index.
func float64At(raw Op, index int) (float64, error) {
	if len(raw) <= index {
		return 0, nil
	}
	return toFloat64(raw[index])
}

// toBool converts a value to bool.
func toBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case int:
		return val != 0
	default:
		return false
	}
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
