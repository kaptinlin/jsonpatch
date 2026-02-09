package binary

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
	op.NewFlip([]string{"flag"}),
	op.NewInc([]string{"counter"}, 5),
	op.NewDefined([]string{"exists"}),
	op.NewUndefined([]string{"missing"}),
	op.NewContains([]string{"text"}, "search"),
	op.NewStarts([]string{"prefix"}, "start"),
	op.NewEnds([]string{"suffix"}, "end"),
	op.NewLess([]string{"number"}, 100.0),
	op.NewMore([]string{"number"}, 1.0),
	op.NewIn([]string{"array"}, []any{"a", "b", "c"}),
	op.NewStrIns([]string{"text"}, 0, "prefix"),
	op.NewStrDel([]string{"text"}, 0, 5),
	op.NewSplit([]string{"text"}, 5, nil),
	op.NewExtend([]string{"object"}, map[string]any{"key": "value"}, false),
	op.NewMerge([]string{"array"}, 1, map[string]any{"merged": true}),
}

func BenchmarkEncode(b *testing.B) {
	codec := &Codec{}

	b.ResetTimer()
	for b.Loop() {
		_, err := codec.Encode(benchmarkOps)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	codec := &Codec{}
	encoded, err := codec.Encode(benchmarkOps)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for b.Loop() {
		_, err := codec.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	codec := &Codec{}

	b.ResetTimer()
	for b.Loop() {
		encoded, err := codec.Encode(benchmarkOps)
		if err != nil {
			b.Fatal(err)
		}

		_, err = codec.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSingleOperationTypes(b *testing.B) {
	codec := &Codec{}

	tests := []struct {
		name string
		op   internal.Op
	}{
		{"Add", op.NewAdd([]string{"foo"}, "bar")},
		{"Remove", op.NewRemove([]string{"baz"})},
		{"Replace", op.NewReplace([]string{"qux"}, 42)},
		{"Move", op.NewMove([]string{"dest"}, []string{"src"})},
		{"Copy", op.NewCopy([]string{"copy_dest"}, []string{"copy_src"})},
		{"Test", op.NewTest([]string{"check"}, "expected")},
		{"Flip", op.NewFlip([]string{"flag"})},
		{"Inc", op.NewInc([]string{"counter"}, 5)},
		{"Defined", op.NewDefined([]string{"exists"})},
		{"Undefined", op.NewUndefined([]string{"missing"})},
		{"Contains", op.NewContains([]string{"text"}, "search")},
		{"Starts", op.NewStarts([]string{"prefix"}, "start")},
		{"Ends", op.NewEnds([]string{"suffix"}, "end")},
		{"Less", op.NewLess([]string{"number"}, 100.0)},
		{"More", op.NewMore([]string{"number"}, 1.0)},
		{"In", op.NewIn([]string{"array"}, []any{"a", "b", "c"})},
		{"StrIns", op.NewStrIns([]string{"text"}, 0, "prefix")},
		{"StrDel", op.NewStrDel([]string{"text"}, 0, 5)},
		{"Split", op.NewSplit([]string{"text"}, 5, nil)},
		{"Extend", op.NewExtend([]string{"object"}, map[string]any{"key": "value"}, false)},
		{"Merge", op.NewMerge([]string{"array"}, 1, map[string]any{"merged": true})},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ops := []internal.Op{tt.op}
			b.ResetTimer()
			for b.Loop() {
				encoded, err := codec.Encode(ops)
				if err != nil {
					b.Fatal(err)
				}

				_, err = codec.Decode(encoded)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
