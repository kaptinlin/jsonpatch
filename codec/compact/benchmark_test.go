package compact

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

var benchmarkOps = []internal.Op{
	op.NewAdd([]string{"foo"}, "bar"),
	op.NewRemove([]string{"baz"}),
	op.NewReplace([]string{"qux"}, 42),
	op.NewMove([]string{"dest"}, []string{"src"}),
	op.NewCopy([]string{"copy_dest"}, []string{"copy_src"}),
	op.NewTest([]string{"check"}, "expected"),
	op.NewOpFlipOperation([]string{"flag"}),
	op.NewOpIncOperation([]string{"counter"}, 5),
	op.NewOpDefinedOperation([]string{"exists"}),
	op.NewOpUndefinedOperation([]string{"missing"}),
}

func BenchmarkEncode(b *testing.B) {
	encoder := NewEncoder()

	b.ResetTimer()
	for b.Loop() {
		_, err := encoder.EncodeSlice(benchmarkOps)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	encoder := NewEncoder()
	encoded, err := encoder.EncodeSlice(benchmarkOps)
	if err != nil {
		b.Fatal(err)
	}

	decoder := NewDecoder()

	b.ResetTimer()
	for b.Loop() {
		_, err := decoder.DecodeSlice(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeJSON(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_, err := EncodeJSON(benchmarkOps)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeJSON(b *testing.B) {
	jsonData, err := EncodeJSON(benchmarkOps)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := DecodeJSON(jsonData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		encoded, err := Encode(benchmarkOps)
		if err != nil {
			b.Fatal(err)
		}

		_, err = Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStringVsNumericOpcode(b *testing.B) {
	singleOp := []internal.Op{op.NewAdd([]string{"foo"}, "bar")}

	b.Run("Numeric", func(b *testing.B) {
		encoder := NewEncoder() // default uses numeric opcodes
		b.ResetTimer()
		for b.Loop() {
			_, err := encoder.EncodeSlice(singleOp)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("String", func(b *testing.B) {
		encoder := NewEncoder(WithStringOpcode(true))
		b.ResetTimer()
		for b.Loop() {
			_, err := encoder.EncodeSlice(singleOp)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
