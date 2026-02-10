package jsonpatch_test

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch"
)

func TestValidateOperations(t *testing.T) {
	t.Parallel()
	t.Run("throws on not an array", func(t *testing.T) {
		t.Parallel()
		err := jsonpatch.ValidateOperations(nil, false)
		if !errors.Is(err, jsonpatch.ErrNotArray) {
			t.Errorf("ValidateOperations(nil) error = %v, want %v", err, jsonpatch.ErrNotArray)
		}
	})

	t.Run("throws on empty array", func(t *testing.T) {
		t.Parallel()
		err := jsonpatch.ValidateOperations([]jsonpatch.Operation{}, false)
		if !errors.Is(err, jsonpatch.ErrEmptyPatch) {
			t.Errorf("ValidateOperations([]) error = %v, want %v", err, jsonpatch.ErrEmptyPatch)
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
			t.Errorf("ValidateOperations() error = %v, want %v", err, jsonpatch.ErrMissingOp)
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
			t.Errorf("ValidateOperations() error = %v, want %v", err, jsonpatch.ErrMissingOp)
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
			t.Errorf("ValidateOperations() error = %v, want %v", err, jsonpatch.ErrInvalidOperation)
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
			t.Errorf("ValidateOperations() error = %v, want %v", err, jsonpatch.ErrMissingValue)
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
			t.Errorf("ValidateOperations() error = %v, want %v", err, jsonpatch.ErrInvalidJSONPointer)
		}
	})
}
