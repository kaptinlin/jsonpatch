// Package main demonstrates binary codec operations using JSON Patch.
package main

import (
	"fmt"
	"log"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/codec/binary"
	codecjson "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func main() {
	fmt.Println("=== Binary Codec Demo ===")

	// Sample operations
	ops := []internal.Op{
		op.NewAdd([]string{"users", "0", "name"}, "John Doe"),
		op.NewReplace([]string{"users", "0", "age"}, 30),
		op.NewFlip([]string{"active"}),
	}

	fmt.Printf("Operations to encode: %d\n\n", len(ops))

	// Standard JSON format size (pretty-printed for clarity)
	jsonOps, err := codecjson.Encode(ops)
	if err != nil {
		log.Fatal(err)
	}
	jsonBytes, _ := json.Marshal(jsonOps)
	fmt.Printf("Standard JSON size: %d bytes\n", len(jsonBytes))

	// Binary codec
	binCodec := binary.Codec{}
	binaryData, err := binCodec.Encode(ops)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Binary (MessagePack) size: %d bytes\n", len(binaryData))
	fmt.Println()

	// Decode back and verify
	decodedOps, err := binCodec.Decode(binaryData)
	if err != nil {
		log.Fatal(err)
	}

	if len(decodedOps) == len(ops) {
		fmt.Println("✅ Round-trip successful (operation count matches)")
	} else {
		fmt.Println("❌ Round-trip failed (operation count mismatch)")
	}

	fmt.Println("\n=== Quick Usage ===")
	fmt.Println("codec := binary.Codec{}")
	fmt.Println("data, _ := codec.Encode(ops)")
	fmt.Println("decoded, _ := codec.Decode(data)")
}
