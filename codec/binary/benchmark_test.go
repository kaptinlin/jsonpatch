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
	op.NewOpFlipOperation([]string{"flag"}),
	op.NewOpIncOperation([]string{"counter"}, 5),
	op.NewOpDefinedOperation([]string{"exists"}),
	op.NewOpUndefinedOperation([]string{"missing"}),
	op.NewOpContainsOperation([]string{"text"}, "search"),
	op.NewOpStartsOperation([]string{"prefix"}, "start"),
	op.NewOpEndsOperation([]string{"suffix"}, "end"),
	op.NewOpLessOperation([]string{"number"}, 100.0),
	op.NewOpMoreOperation([]string{"number"}, 1.0),
	op.NewOpInOperation([]string{"array"}, []interface{}{"a", "b", "c"}),
	op.NewOpStrInsOperation([]string{"text"}, 0, "prefix"),
	op.NewOpStrDelOperation([]string{"text"}, 0, 5),
	op.NewOpSplitOperation([]string{"text"}, 5, nil),
	op.NewOpExtendOperation([]string{"object"}, map[string]interface{}{"key": "value"}, false),
	op.NewOpMergeOperation([]string{"array"}, 1, map[string]interface{}{"merged": true}),
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
		{"Flip", op.NewOpFlipOperation([]string{"flag"})},
		{"Inc", op.NewOpIncOperation([]string{"counter"}, 5)},
		{"Defined", op.NewOpDefinedOperation([]string{"exists"})},
		{"Undefined", op.NewOpUndefinedOperation([]string{"missing"})},
		{"Contains", op.NewOpContainsOperation([]string{"text"}, "search")},
		{"Starts", op.NewOpStartsOperation([]string{"prefix"}, "start")},
		{"Ends", op.NewOpEndsOperation([]string{"suffix"}, "end")},
		{"Less", op.NewOpLessOperation([]string{"number"}, 100.0)},
		{"More", op.NewOpMoreOperation([]string{"number"}, 1.0)},
		{"In", op.NewOpInOperation([]string{"array"}, []interface{}{"a", "b", "c"})},
		{"StrIns", op.NewOpStrInsOperation([]string{"text"}, 0, "prefix")},
		{"StrDel", op.NewOpStrDelOperation([]string{"text"}, 0, 5)},
		{"Split", op.NewOpSplitOperation([]string{"text"}, 5, nil)},
		{"Extend", op.NewOpExtendOperation([]string{"object"}, map[string]interface{}{"key": "value"}, false)},
		{"Merge", op.NewOpMergeOperation([]string{"array"}, 1, map[string]interface{}{"merged": true})},
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
