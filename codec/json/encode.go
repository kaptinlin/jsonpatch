package json

import (
	"github.com/go-json-experiment/json"
	"github.com/kaptinlin/jsonpatch/internal"
)

// Encode converts operations to JSON format.
func Encode(ops []internal.Op) ([]internal.Operation, error) {
	operations := make([]internal.Operation, len(ops))
	for i, o := range ops {
		jsonOp, err := o.ToJSON()
		if err != nil {
			return nil, err
		}
		operations[i] = jsonOp
	}
	return operations, nil
}

// EncodeJSON converts operations to JSON bytes.
func EncodeJSON(ops []internal.Op) ([]byte, error) {
	operations, err := Encode(ops)
	if err != nil {
		return nil, err
	}
	return json.Marshal(operations)
}
