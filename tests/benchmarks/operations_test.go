package jsonpatch_test

import (
	"testing"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func BenchmarkBasicOperations(b *testing.B) {
	testCases := []struct {
		name string
		doc  any
		ops  []jsoncodec.Operation
	}{
		{
			name: "add_simple_value",
			doc:  map[string]any{"foo": "bar"},
			ops: []jsoncodec.Operation{
				{Op: "add", Path: "/baz", Value: "qux"},
			},
		},
		{
			name: "replace_value",
			doc:  map[string]any{"foo": "bar", "baz": "qux"},
			ops: []jsoncodec.Operation{
				{Op: "replace", Path: "/foo", Value: "new_value"},
			},
		},
		{
			name: "remove_value",
			doc:  map[string]any{"foo": "bar", "baz": "qux"},
			ops: []jsoncodec.Operation{
				{Op: "remove", Path: "/baz"},
			},
		},
		{
			name: "test_operation",
			doc:  map[string]any{"foo": "bar"},
			ops: []jsoncodec.Operation{
				{Op: "test", Path: "/foo", Value: "bar"},
			},
		},
		{
			name: "copy_operation",
			doc:  map[string]any{"foo": "bar", "baz": map[string]any{"deep": "value"}},
			ops: []jsoncodec.Operation{
				{Op: "copy", From: "/baz", Path: "/copied"},
			},
		},
		{
			name: "move_operation",
			doc:  map[string]any{"foo": "bar", "baz": "qux"},
			ops: []jsoncodec.Operation{
				{Op: "move", From: "/baz", Path: "/moved"},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			patch := compileBenchmarkPatch(b, tc.ops)
			b.ResetTimer()
			for b.Loop() {
				docCopy := cloneDocument(tc.doc)
				applyBenchmarkPatch(b, patch, docCopy, true)
			}
		})
	}
}

func BenchmarkExtendedOperations(b *testing.B) {
	testCases := []struct {
		name string
		doc  any
		ops  []jsoncodec.Operation
	}{
		{
			name: "inc_operation",
			doc:  map[string]any{"counter": 42},
			ops: []jsoncodec.Operation{
				{Op: "inc", Path: "/counter", Inc: 1},
			},
		},
		{
			name: "flip_operation",
			doc:  map[string]any{"enabled": true},
			ops: []jsoncodec.Operation{
				{Op: "flip", Path: "/enabled"},
			},
		},
		{
			name: "strins_operation",
			doc:  map[string]any{"text": "hello world"},
			ops: []jsoncodec.Operation{
				{Op: "str_ins", Path: "/text", Pos: 5, Str: " beautiful"},
			},
		},
		{
			name: "split_operation",
			doc:  map[string]any{"text": "hello,world,test"},
			ops: []jsoncodec.Operation{
				{Op: "split", Path: "/text", Pos: 5},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			patch := compileBenchmarkPatch(b, tc.ops)
			b.ResetTimer()
			for b.Loop() {
				docCopy := cloneDocument(tc.doc)
				applyBenchmarkPatch(b, patch, docCopy, true)
			}
		})
	}
}

func BenchmarkPredicateOperations(b *testing.B) {
	testCases := []struct {
		name string
		doc  any
		ops  []jsoncodec.Operation
	}{
		{
			name: "contains_predicate",
			doc:  map[string]any{"text": "hello world"},
			ops: []jsoncodec.Operation{
				{Op: "contains", Path: "/text", Value: "world"},
			},
		},
		{
			name: "starts_predicate",
			doc:  map[string]any{"text": "hello world"},
			ops: []jsoncodec.Operation{
				{Op: "starts", Path: "/text", Value: "hello"},
			},
		},
		{
			name: "ends_predicate",
			doc:  map[string]any{"text": "hello world"},
			ops: []jsoncodec.Operation{
				{Op: "ends", Path: "/text", Value: "world"},
			},
		},
		{
			name: "type_predicate",
			doc:  map[string]any{"value": 42},
			ops: []jsoncodec.Operation{
				{Op: "type", Path: "/value", Value: "number"},
			},
		},
		{
			name: "less_predicate",
			doc:  map[string]any{"value": 42},
			ops: []jsoncodec.Operation{
				{Op: "less", Path: "/value", Value: 50},
			},
		},
		{
			name: "more_predicate",
			doc:  map[string]any{"value": 42},
			ops: []jsoncodec.Operation{
				{Op: "more", Path: "/value", Value: 30},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			patch := compileBenchmarkPatch(b, tc.ops)
			b.ResetTimer()
			for b.Loop() {
				docCopy := cloneDocument(tc.doc)
				applyBenchmarkPatch(b, patch, docCopy, true)
			}
		})
	}
}

func BenchmarkSecondOrderPredicates(b *testing.B) {
	testCases := []struct {
		name string
		doc  any
		ops  []jsoncodec.Operation
	}{
		{
			name: "not_predicate",
			doc:  map[string]any{"foo": 1, "bar": 2},
			ops: []jsoncodec.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsoncodec.Operation{
						{Op: "test", Path: "/foo", Value: 2},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			patch := compileBenchmarkPatch(b, tc.ops)
			b.ResetTimer()
			for b.Loop() {
				docCopy := cloneDocument(tc.doc)
				applyBenchmarkPatch(b, patch, docCopy, true)
			}
		})
	}
}

func BenchmarkComplexDocument(b *testing.B) {
	complexDoc := map[string]any{
		"users": []any{
			map[string]any{
				"id":   1,
				"name": "Alice",
				"profile": map[string]any{
					"age":   30,
					"email": "alice@example.com",
					"preferences": map[string]any{
						"theme": "dark",
						"lang":  "en",
					},
				},
			},
			map[string]any{
				"id":   2,
				"name": "Bob",
				"profile": map[string]any{
					"age":   25,
					"email": "bob@example.com",
					"preferences": map[string]any{
						"theme": "light",
						"lang":  "es",
					},
				},
			},
		},
		"metadata": map[string]any{
			"version":   "1.0",
			"timestamp": "2024-01-01T00:00:00Z",
			"stats": map[string]any{
				"total_users":   2,
				"active_users":  1,
				"last_modified": "2024-01-01T12:00:00Z",
			},
		},
	}

	testCases := []struct {
		name string
		ops  []jsoncodec.Operation
	}{
		{
			name: "add_new_user",
			ops: []jsoncodec.Operation{
				{
					Op:   "add",
					Path: "/users/-",
					Value: map[string]any{
						"id":   3,
						"name": "Charlie",
						"profile": map[string]any{
							"age":   28,
							"email": "charlie@example.com",
						},
					},
				},
			},
		},
		{
			name: "update_user_preference",
			ops: []jsoncodec.Operation{
				{Op: "replace", Path: "/users/0/profile/preferences/theme", Value: "light"},
			},
		},
		{
			name: "increment_stats",
			ops: []jsoncodec.Operation{
				{Op: "inc", Path: "/metadata/stats/total_users", Inc: 1},
				{Op: "inc", Path: "/metadata/stats/active_users", Inc: 1},
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			patch := compileBenchmarkPatch(b, tc.ops)
			b.ResetTimer()
			for b.Loop() {
				docCopy := cloneDocument(complexDoc)
				applyBenchmarkPatch(b, patch, docCopy, true)
			}
		})
	}
}

func BenchmarkApplyInPlaceVsApply(b *testing.B) {
	doc := map[string]any{
		"data": make([]any, 1000),
		"metadata": map[string]any{
			"size":    1000,
			"version": 1,
		},
	}

	for i := range 1000 {
		doc["data"].([]any)[i] = map[string]any{
			"id":    i,
			"value": i * 2,
		}
	}

	ops := []jsoncodec.Operation{
		{Op: "replace", Path: "/metadata/version", Value: 2},
		{Op: "add", Path: "/metadata/updated", Value: true},
	}

	patch := compileBenchmarkPatch(b, ops)

	b.Run("in_place", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			docCopy := cloneDocument(doc)
			applyBenchmarkPatch(b, patch, docCopy, true)
		}
	})

	b.Run("immutable", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			docCopy := cloneDocument(doc)
			applyBenchmarkPatch(b, patch, docCopy, false)
		}
	})
}

func compileBenchmarkPatch(b testing.TB, operations []jsoncodec.Operation) *jsonpatch.Patch {
	b.Helper()
	patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		b.Fatalf("CompileOperations failed: %v", err)
	}
	return patch
}

func applyBenchmarkPatch(b testing.TB, patch *jsonpatch.Patch, doc any, inPlace bool) {
	b.Helper()
	if inPlace {
		if err := jsonpatch.ApplyInPlace(patch, &doc); err != nil {
			b.Fatalf("ApplyInPlace failed: %v", err)
		}
		return
	}
	if _, err := jsonpatch.Apply(patch, doc); err != nil {
		b.Fatalf("Apply failed: %v", err)
	}
}

func cloneDocument(doc any) any {
	data, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	var result any
	err = json.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result
}
