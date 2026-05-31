package op

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemove_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]any{"nested": "value"},
	}

	removeOp := NewRemove([]string{"foo"})
	result, err := removeOp.Apply(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() unexpected error: %v", err))
	}

	modifiedDoc := result.Doc.(map[string]any)
	if got := result.Old; got != "bar" {
		assert.Equal(t, "bar", got, "result.Old")
	}
	if _, ok := modifiedDoc["foo"]; ok {
		assert.Fail(t, "modifiedDoc contains key \"foo\" after remove")
	}
	if _, ok := modifiedDoc["baz"]; !ok {
		assert.Fail(t, "modifiedDoc missing key \"baz\"")
	}
	if _, ok := modifiedDoc["qux"]; !ok {
		assert.Fail(t, "modifiedDoc missing key \"qux\"")
	}
}

func TestRemove_Nested(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": map[string]any{"bar": "baz", "qux": 123},
	}

	removeOp := NewRemove([]string{"foo", "bar"})
	result, err := removeOp.Apply(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() unexpected error: %v", err))
	}

	modifiedDoc := result.Doc.(map[string]any)
	foo := modifiedDoc["foo"].(map[string]any)
	if got := result.Old; got != "baz" {
		assert.Equal(t, "baz", got, "result.Old")
	}
	if _, ok := foo["bar"]; ok {
		assert.Fail(t, "foo contains key \"bar\" after remove")
	}
	if _, ok := foo["qux"]; !ok {
		assert.Fail(t, "foo missing key \"qux\"")
	}
}

func TestRemove_Array(t *testing.T) {
	t.Parallel()
	doc := []any{"first", "second", "third"}

	removeOp := NewRemove([]string{"1"})
	result, err := removeOp.Apply(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() unexpected error: %v", err))
	}

	modifiedArray := result.Doc.([]any)
	if got := result.Old; got != "second" {
		assert.Equal(t, "second", got, "result.Old")
	}
	want := []any{"first", "third"}
	if diff := cmp.Diff(want, modifiedArray); diff != "" {
		t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
	}
}

func TestRemove_NestedArrayUpdatesParent(t *testing.T) {
	t.Parallel()

	doc := map[string]any{"items": []any{"first", "second", "third"}}
	result, err := NewRemove([]string{"items", "1"}).Apply(doc)
	require.NoError(t, err)

	want := map[string]any{"items": []any{"first", "third"}}
	if diff := cmp.Diff(want, result.Doc); diff != "" {
		t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, "second", result.Old)
}

func TestRemove_RejectsScalarDocument(t *testing.T) {
	t.Parallel()

	_, err := NewRemove([]string{"name"}).Apply("Ada")
	require.ErrorIs(t, err, ErrCannotRemoveFromValue)
}

func TestRemove_RejectsNestedArrayIndexOutOfRange(t *testing.T) {
	t.Parallel()

	doc := map[string]any{"items": []any{"first"}}
	_, err := NewRemove([]string{"items", "1"}).Apply(doc)
	require.ErrorIs(t, err, ErrIndexOutOfRange)
}

func TestRemove_NonExistent(t *testing.T) {
	t.Parallel()
	doc := map[string]any{"foo": "bar"}

	removeOp := NewRemove([]string{"qux"})
	_, err := removeOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "Apply() expected error for non-existent path")
	}
	if !errors.Is(err, ErrPathNotFound) {
		assert.Equal(t, ErrPathNotFound, err, "Apply() error")
	}
}

func TestRemove_EmptyPath(t *testing.T) {
	t.Parallel()
	doc := map[string]any{"foo": "bar"}

	removeOp := NewRemove([]string{})
	result, err := removeOp.Apply(doc)
	require.NoError(t, err)
	assert.Nil(t, result.Doc)
	assert.Equal(t, doc, result.Old)
}
