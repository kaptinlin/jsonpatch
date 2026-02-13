package jsonpatch_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kaptinlin/jsonpatch"
)

func TestValidateOperations(t *testing.T) {
	t.Parallel()
	t.Run("throws on not an array", func(t *testing.T) {
		t.Parallel()
		err := jsonpatch.ValidateOperations(nil, false)
		if !errors.Is(err, jsonpatch.ErrNotArray) {
			assert.Equal(t, jsonpatch.ErrNotArray, err, "ValidateOperations(nil) error")
		}
	})

	t.Run("throws on empty array", func(t *testing.T) {
		t.Parallel()
		err := jsonpatch.ValidateOperations([]jsonpatch.Operation{}, false)
		if !errors.Is(err, jsonpatch.ErrEmptyPatch) {
			assert.Equal(t, jsonpatch.ErrEmptyPatch, err, "ValidateOperations([]) error")
		}
	})

	t.Run("throws on no operation path", func(t *testing.T) {
		t.Parallel()
		ops := []jsonpatch.Operation{{}}
		err := jsonpatch.ValidateOperations(ops, false)
		if err == nil {
			t.Fatal("ValidateOperations() error = nil, want error")
		}
		if !errors.Is(err, jsonpatch.ErrMissingOp) {
			assert.Equal(t, jsonpatch.ErrMissingOp, err, "ValidateOperations() error")
		}
	})

	t.Run("throws on no operation code", func(t *testing.T) {
		t.Parallel()
		ops := []jsonpatch.Operation{{Path: "/"}}
		err := jsonpatch.ValidateOperations(ops, false)
		if err == nil {
			t.Fatal("ValidateOperations() error = nil, want error")
		}
		if !errors.Is(err, jsonpatch.ErrMissingOp) {
			assert.Equal(t, jsonpatch.ErrMissingOp, err, "ValidateOperations() error")
		}
	})

	t.Run("throws on invalid operation code", func(t *testing.T) {
		t.Parallel()
		ops := []jsonpatch.Operation{{Path: "/", Op: "123"}}
		err := jsonpatch.ValidateOperations(ops, false)
		if err == nil {
			t.Fatal("ValidateOperations() error = nil, want error")
		}
		if !errors.Is(err, jsonpatch.ErrInvalidOperation) {
			assert.Equal(t, jsonpatch.ErrInvalidOperation, err, "ValidateOperations() error")
		}
	})

	t.Run("succeeds on valid operation", func(t *testing.T) {
		t.Parallel()
		ops := []jsonpatch.Operation{{Op: "add", Path: "/test", Value: 123}}
		err := jsonpatch.ValidateOperations(ops, false)
		if err != nil {
			t.Errorf("ValidateOperations() error = %v, want nil", err)
		}
	})

	t.Run("throws on second invalid operation", func(t *testing.T) {
		t.Parallel()
		ops := []jsonpatch.Operation{
			{Op: "add", Path: "/test", Value: 123},
			{Op: "test", Path: "/test"},
		}
		err := jsonpatch.ValidateOperations(ops, false)
		if err == nil {
			t.Fatal("ValidateOperations() error = nil, want error")
		}
		if !errors.Is(err, jsonpatch.ErrMissingValue) {
			assert.Equal(t, jsonpatch.ErrMissingValue, err, "ValidateOperations() error")
		}
	})

	t.Run("throws if JSON pointer does not start with forward slash", func(t *testing.T) {
		t.Parallel()
		ops := []jsonpatch.Operation{
			{Op: "add", Path: "/test", Value: 123},
			{Op: "test", Path: "test", Value: 1},
		}
		err := jsonpatch.ValidateOperations(ops, false)
		if err == nil {
			t.Fatal("ValidateOperations() error = nil, want error")
		}
		if !errors.Is(err, jsonpatch.ErrInvalidJSONPointer) {
			assert.Equal(t, jsonpatch.ErrInvalidJSONPointer, err, "ValidateOperations() error")
		}
	})
}
