package json

import (
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/internal"
)

// Encode converts Op instances to Operation structs.
func Encode(ops []internal.Op) ([]internal.Operation, error) {
	result := make([]internal.Operation, len(ops))
	for i, o := range ops {
		jsonOp, ok := o.(internal.JSONOp)
		if !ok {
			return nil, fmt.Errorf("operation %T: %w", o, ErrUnsupportedOp)
		}
		encoded, err := jsonOp.ToJSON()
		if err != nil {
			return nil, err
		}
		result[i] = encoded
	}
	return result, nil
}

// EncodeJSON converts Op instances to JSON bytes.
func EncodeJSON(ops []internal.Op) ([]byte, error) {
	result, err := Encode(ops)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result)
}
