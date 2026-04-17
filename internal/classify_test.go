package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsJSONPatchOperation(t *testing.T) {
	t.Parallel()
	valid := []string{"add", "remove", "replace", "move", "copy", "test"}
	for _, op := range valid {
		if !IsJSONPatchOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsJSONPatchOperation(%q) = false, want true", op))
		}
	}

	invalid := []string{"", "inc", "flip", "and", "defined", "unknown"}
	for _, op := range invalid {
		if IsJSONPatchOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsJSONPatchOperation(%q) = true, want false", op))
		}
	}
}

func TestIsFirstOrderPredicateOperation(t *testing.T) {
	t.Parallel()
	valid := []string{
		"test", "defined", "undefined", "test_type",
		"test_string", "test_string_len", "contains",
		"ends", "starts", "in", "less", "more", "matches",
	}
	for _, op := range valid {
		if !IsFirstOrderPredicateOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsFirstOrderPredicateOperation(%q) = false, want true", op))
		}
	}

	invalid := []string{"", "add", "and", "or", "not", "inc"}
	for _, op := range invalid {
		if IsFirstOrderPredicateOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsFirstOrderPredicateOperation(%q) = true, want false", op))
		}
	}
}

func TestIsSecondOrderPredicateOperation(t *testing.T) {
	t.Parallel()
	valid := []string{"and", "or", "not"}
	for _, op := range valid {
		if !IsSecondOrderPredicateOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsSecondOrderPredicateOperation(%q) = false, want true", op))
		}
	}

	invalid := []string{"", "add", "test", "defined", "inc"}
	for _, op := range invalid {
		if IsSecondOrderPredicateOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsSecondOrderPredicateOperation(%q) = true, want false", op))
		}
	}
}

func TestIsPredicateOperation(t *testing.T) {
	t.Parallel()
	firstOrder := []string{
		"test", "defined", "undefined", "test_type",
		"test_string", "test_string_len", "contains",
		"ends", "starts", "in", "less", "more", "matches",
	}
	for _, op := range firstOrder {
		if !IsPredicateOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsPredicateOperation(%q) = false, want true (first-order)", op))
		}
	}

	secondOrder := []string{"and", "or", "not"}
	for _, op := range secondOrder {
		if !IsPredicateOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsPredicateOperation(%q) = false, want true (second-order)", op))
		}
	}

	invalid := []string{"", "add", "remove", "inc", "flip", "unknown"}
	for _, op := range invalid {
		if IsPredicateOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsPredicateOperation(%q) = true, want false", op))
		}
	}
}

func TestIsJSONPatchExtendedOperation(t *testing.T) {
	t.Parallel()
	valid := []string{"str_ins", "str_del", "flip", "inc", "split", "merge", "extend"}
	for _, op := range valid {
		if !IsJSONPatchExtendedOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsJSONPatchExtendedOperation(%q) = false, want true", op))
		}
	}

	invalid := []string{"", "add", "test", "and", "defined", "unknown"}
	for _, op := range invalid {
		if IsJSONPatchExtendedOperation(op) {
			assert.Fail(t, fmt.Sprintf("IsJSONPatchExtendedOperation(%q) = true, want false", op))
		}
	}
}
