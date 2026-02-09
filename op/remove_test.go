package op

import (
	"errors"
	"testing"
)

func TestRemove_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]any{"nested": "value"},
	}

	removeOp := NewRemove([]string{"foo"})
	result, err := removeOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	modifiedDoc := result.Doc.(map[string]any)
	if got := result.Old; got != "bar" {
		t.Errorf("result.Old = %v, want %v", got, "bar")
	}
	if _, ok := modifiedDoc["foo"]; ok {
		t.Error("modifiedDoc contains key \"foo\" after remove")
	}
	if _, ok := modifiedDoc["baz"]; !ok {
		t.Error("modifiedDoc missing key \"baz\"")
	}
	if _, ok := modifiedDoc["qux"]; !ok {
		t.Error("modifiedDoc missing key \"qux\"")
	}
}

func TestRemove_Nested(t *testing.T) {
	doc := map[string]any{
		"foo": map[string]any{"bar": "baz", "qux": 123},
	}

	removeOp := NewRemove([]string{"foo", "bar"})
	result, err := removeOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	modifiedDoc := result.Doc.(map[string]any)
	foo := modifiedDoc["foo"].(map[string]any)
	if got := result.Old; got != "baz" {
		t.Errorf("result.Old = %v, want %v", got, "baz")
	}
	if _, ok := foo["bar"]; ok {
		t.Error("foo contains key \"bar\" after remove")
	}
	if _, ok := foo["qux"]; !ok {
		t.Error("foo missing key \"qux\"")
	}
}

func TestRemove_Array(t *testing.T) {
	doc := []any{"first", "second", "third"}

	removeOp := NewRemove([]string{"1"})
	result, err := removeOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	modifiedArray := result.Doc.([]any)
	if got := result.Old; got != "second" {
		t.Errorf("result.Old = %v, want %v", got, "second")
	}
	if len(modifiedArray) != 2 {
		t.Fatalf("len(modifiedArray) = %d, want %d", len(modifiedArray), 2)
	}
	if modifiedArray[0] != "first" {
		t.Errorf("modifiedArray[0] = %v, want %v", modifiedArray[0], "first")
	}
	if modifiedArray[1] != "third" {
		t.Errorf("modifiedArray[1] = %v, want %v", modifiedArray[1], "third")
	}
}

func TestRemove_NonExistent(t *testing.T) {
	doc := map[string]any{"foo": "bar"}

	removeOp := NewRemove([]string{"qux"})
	_, err := removeOp.Apply(doc)
	if err == nil {
		t.Error("Apply() expected error for non-existent path")
	}
	if !errors.Is(err, ErrPathNotFound) {
		t.Errorf("Apply() error = %v, want %v", err, ErrPathNotFound)
	}
}

func TestRemove_EmptyPath(t *testing.T) {
	doc := map[string]any{"foo": "bar"}

	removeOp := NewRemove([]string{})
	_, err := removeOp.Apply(doc)
	if err == nil {
		t.Error("Apply() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		t.Errorf("Apply() error = %v, want %v", err, ErrPathEmpty)
	}
}
