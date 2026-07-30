package op_test

import "github.com/kaptinlin/jsonpatch/op"

var (
	_ func([]string, int, string) *op.StrInsOperation           = op.NewStrIns
	_ func([]string, int, int) *op.StrDelOperation              = op.NewStrDel
	_ func([]string, int, string) *op.StrDelOperation           = op.NewStrDelWithStr
	_ func([]string, int, any) *op.SplitOperation               = op.NewSplit
	_ func([]string, int, map[string]any) *op.MergeOperation    = op.NewMerge
	_ func([]string, string, int, bool) *op.TestStringOperation = op.NewTestString
	_ func([]string, int) *op.TestStringLenOperation            = op.NewTestStringLen
	_ func([]string, int, bool) *op.TestStringLenOperation      = op.NewTestStringLenWithNot
)
