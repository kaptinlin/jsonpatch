package binary

// normalizeMap recursively converts map[any]any to map[string]any
// and normalizes nested values in maps and slices.
func normalizeMap(v any) any {
	switch m := v.(type) {
	case map[any]any:
		res := make(map[string]any, len(m))
		for key, val := range m {
			if k, ok := key.(string); ok {
				res[k] = normalizeMap(val)
			}
		}
		return res
	case map[string]any:
		for k, val := range m {
			m[k] = normalizeMap(val)
		}
		return m
	case []any:
		for i, val := range m {
			m[i] = normalizeMap(val)
		}
		return m
	default:
		return v
	}
}
