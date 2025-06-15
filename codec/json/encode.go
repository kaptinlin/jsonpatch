// Package json implements JSON codec for JSON Patch operations.
package json

import (
	"github.com/go-json-experiment/json"
	"github.com/kaptinlin/jsonpatch/internal"
)

// Encode converts operations to JSON format.
func Encode(ops []internal.Op) ([]map[string]interface{}, error) {
	operations := make([]map[string]interface{}, 0, len(ops))
	for _, operation := range ops {
		jsonOp, err := operation.ToJSON()
		if err != nil {
			return nil, err
		}
		operations = append(operations, jsonOp)
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
