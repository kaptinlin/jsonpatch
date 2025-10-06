// Package main demonstrates the compact codec functionality.
// This example shows how to use the compact codec to achieve significant space savings
// compared to standard JSON Patch format.
package main

import (
	"fmt"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"log"

	"github.com/kaptinlin/jsonpatch/codec/compact"
	codecjson "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func main() {
	fmt.Println("=== Compact Codec Demo ===")

	// Create sample operations
	ops := []internal.Op{
		op.NewAdd([]string{"users", "0", "name"}, "John Doe"),
		op.NewAdd([]string{"users", "0", "email"}, "john@example.com"),
		op.NewReplace([]string{"users", "0", "age"}, 30),
		op.NewRemove([]string{"temp"}),
		op.NewMove([]string{"users", "0", "profile"}, []string{"users", "0", "info"}),
		op.NewCopy([]string{"users", "1"}, []string{"users", "0"}),
		op.NewTest([]string{"version"}, "1.0"),
	}

	fmt.Printf("Operations to encode: %d\n\n", len(ops))

	// Standard JSON format
	fmt.Println("1. Standard JSON Format:")
	jsonOps, err := codecjson.Encode(ops)
	if err != nil {
		log.Fatal(err)
	}
	jsonBytes, _ := json.Marshal(jsonOps, jsontext.Multiline(true))
	fmt.Printf("Size: %d bytes\n", len(jsonBytes))
	fmt.Printf("Content:\n%s\n\n", jsonBytes)

	// Compact format with numeric opcodes
	fmt.Println("2. Compact Format (Numeric Opcodes):")
	compactOps, err := compact.Encode(ops)
	if err != nil {
		log.Fatal(err)
	}
	compactBytes, _ := json.Marshal(compactOps, jsontext.Multiline(true))
	fmt.Printf("Size: %d bytes\n", len(compactBytes))
	fmt.Printf("Content:\n%s\n\n", compactBytes)

	// Compact format with string opcodes
	fmt.Println("3. Compact Format (String Opcodes):")
	compactStringOps, err := compact.Encode(ops, compact.WithStringOpcode(true))
	if err != nil {
		log.Fatal(err)
	}
	compactStringBytes, _ := json.Marshal(compactStringOps, jsontext.Multiline(true))
	fmt.Printf("Size: %d bytes\n", len(compactStringBytes))
	fmt.Printf("Content:\n%s\n\n", compactStringBytes)

	// Calculate space savings
	standardSize := float64(len(jsonBytes))
	compactSize := float64(len(compactBytes))
	compactStringSize := float64(len(compactStringBytes))

	numericSavings := (1 - compactSize/standardSize) * 100
	stringSavings := (1 - compactStringSize/standardSize) * 100

	fmt.Println("=== Space Savings Analysis ===")
	fmt.Printf("Standard JSON:     %d bytes (100%%)\n", len(jsonBytes))
	fmt.Printf("Compact (numeric): %d bytes (%.1f%% savings)\n", len(compactBytes), numericSavings)
	fmt.Printf("Compact (string):  %d bytes (%.1f%% savings)\n", len(compactStringBytes), stringSavings)
	fmt.Println()

	// Demonstrate round-trip compatibility
	fmt.Println("=== Round-trip Compatibility Test ===")

	// Decode compact operations back to operations
	decodedOps, err := compact.Decode(compactOps)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original operations:  %d\n", len(ops))
	fmt.Printf("Decoded operations:   %d\n", len(decodedOps))

	// Verify all operations match
	allMatch := true
	for i, originalOp := range ops {
		if i >= len(decodedOps) {
			allMatch = false
			break
		}
		decodedOp := decodedOps[i]
		if originalOp.Op() != decodedOp.Op() {
			allMatch = false
			break
		}
		// Check paths match
		origPath := originalOp.Path()
		decodedPath := decodedOp.Path()
		if len(origPath) != len(decodedPath) {
			allMatch = false
			break
		}
		for j, segment := range origPath {
			if decodedPath[j] != segment {
				allMatch = false
				break
			}
		}
	}

	if allMatch {
		fmt.Println("✅ All operations perfectly decoded!")
	} else {
		fmt.Println("❌ Some operations failed to decode correctly")
	}

	fmt.Println("\n=== Usage Examples ===")

	// Simple encoding
	fmt.Println("// Simple encoding:")
	fmt.Println("encoded, err := compact.Encode(operations)")
	fmt.Println()

	// With options
	fmt.Println("// With string opcodes:")
	fmt.Println("encoded, err := compact.Encode(operations, compact.WithStringOpcode(true))")
	fmt.Println()

	// JSON marshaling
	fmt.Println("// Direct JSON marshaling:")
	fmt.Println("jsonData, err := compact.EncodeJSON(operations)")
	fmt.Println("decoded, err := compact.DecodeJSON(jsonData)")
	fmt.Println()

	// Using encoder/decoder structs
	fmt.Println("// Using encoder/decoder objects:")
	fmt.Println("encoder := compact.NewEncoder(compact.WithStringOpcode(true))")
	fmt.Println("decoder := compact.NewDecoder()")
	fmt.Println("encoded, err := encoder.EncodeSlice(operations)")
	fmt.Println("decoded, err := decoder.DecodeSlice(encoded)")

	fmt.Println("\n=== Demo Complete ===")
}
