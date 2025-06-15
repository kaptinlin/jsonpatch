package jsonpatch_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
)

// TestScenarious runs scenario tests similar to TypeScript's patch.scenarious.spec.ts
// Original TypeScript: .reference/json-joy/src/json-patch/__tests__/patch.scenarious.spec.ts

// TestCannotAddKeyToEmptyDocument tests that adding key to empty document fails
func TestCannotAddKeyToEmptyDocument(t *testing.T) {
	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":    "add",
			"path":  "/foo",
			"value": 123,
		},
	}

	options := jsonpatch.ApplyPatchOptions{Mutate: true}
	_, err := jsonpatch.ApplyPatch(nil, patch, options)
	if err == nil {
		t.Error("Expected error when adding key to empty document")
	}
}

// TestCanOverwriteEmptyDocument tests that overwriting empty document works
func TestCanOverwriteEmptyDocument(t *testing.T) {
	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":    "add",
			"path":  "/foo",
			"value": 123,
		},
	}

	options := jsonpatch.ApplyPatchOptions{Mutate: true}
	result, err := jsonpatch.ApplyPatch(map[string]interface{}{}, patch, options)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}

	expected := map[string]interface{}{"foo": 123}
	resultMap, ok := result.Doc.(map[string]interface{})
	if !ok {
		t.Fatalf("Result is not a map: %T", result.Doc)
	}

	if resultMap["foo"] != expected["foo"] {
		t.Errorf("Expected %+v, got %+v", expected, resultMap)
	}
}

// TestCannotAddValueToNonexistingPath tests that adding to nonexisting path fails
func TestCannotAddValueToNonexistingPath(t *testing.T) {
	doc := map[string]interface{}{"foo": 123}
	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":    "add",
			"path":  "/foo/bar/baz",
			"value": "test",
		},
	}

	options := jsonpatch.ApplyPatchOptions{Mutate: true}
	_, err := jsonpatch.ApplyPatch(doc, patch, options)
	if err == nil {
		t.Error("Expected error when adding value to nonexisting path")
	}
}
