package binarytests

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/codec/binary"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	jsonsamples "github.com/kaptinlin/jsonpatch/codec/json/tests"
	"github.com/kaptinlin/jsonpatch/internal"
)

// Unsupported composite operations in binary codec for now.
var unsupportedOps = map[string]struct{}{
	"and": {},
	"or":  {},
	"not": {},
}

func TestAutomaticRoundtrip(t *testing.T) {
	t.Parallel()
	binCodec := binary.Codec{}
	options := internal.JSONPatchOptions{CreateMatcher: jsonpatch.CreateMatcherDefault}

	for name, opMap := range jsonsamples.SampleOperations {
		if opType, ok := opMap["op"].(string); ok {
			if _, unsupported := unsupportedOps[opType]; unsupported {
				continue // Skip composite predicates not yet supported
			}
			// Skip operations with features not yet preserved by binary codec
			if _, hasOld := opMap["oldValue"]; hasOld {
				continue
			}
			if _, hasIgnore := opMap["ignore_case"]; hasIgnore {
				continue
			}
			if opType == "str_del" {
				if _, hasStr := opMap["str"]; hasStr {
					continue
				}
			}
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Step 1: JSON -> Op (json codec)
			jsonOps, err := jsoncodec.Decode([]map[string]any{opMap}, options)
			if err != nil {
				t.Fatalf("json Decode should not error: %v", err)
			}

			// Step 2: Op -> Binary bytes
			encoded, err := binCodec.Encode(jsonOps)
			if err != nil {
				t.Fatalf("binary Encode should not error: %v", err)
			}

			// Step 3: Binary bytes -> Op
			decodedOps, err := binCodec.Decode(encoded)
			if err != nil {
				t.Fatalf("binary Decode should not error: %v", err)
			}

			// Step 4: Validate equality between original decoded ops and binary roundtrip
			if !areOpsEqual(jsonOps, decodedOps) {
				t.Error("roundtrip should preserve ops equality")
			}
		})
	}
}
