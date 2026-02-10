package internal

import "testing"

func TestIsJSONPatchOperation(t *testing.T) {
	t.Parallel()
	valid := []string{"add", "remove", "replace", "move", "copy", "test"}
	for _, op := range valid {
		if !IsJSONPatchOperation(op) {
			t.Errorf("IsJSONPatchOperation(%q) = false, want true", op)
		}
	}

	invalid := []string{"", "inc", "flip", "and", "defined", "unknown"}
	for _, op := range invalid {
		if IsJSONPatchOperation(op) {
			t.Errorf("IsJSONPatchOperation(%q) = true, want false", op)
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
			t.Errorf("IsFirstOrderPredicateOperation(%q) = false, want true", op)
		}
	}

	invalid := []string{"", "add", "and", "or", "not", "inc"}
	for _, op := range invalid {
		if IsFirstOrderPredicateOperation(op) {
			t.Errorf("IsFirstOrderPredicateOperation(%q) = true, want false", op)
		}
	}
}

func TestIsSecondOrderPredicateOperation(t *testing.T) {
	t.Parallel()
	valid := []string{"and", "or", "not"}
	for _, op := range valid {
		if !IsSecondOrderPredicateOperation(op) {
			t.Errorf("IsSecondOrderPredicateOperation(%q) = false, want true", op)
		}
	}

	invalid := []string{"", "add", "test", "defined", "inc"}
	for _, op := range invalid {
		if IsSecondOrderPredicateOperation(op) {
			t.Errorf("IsSecondOrderPredicateOperation(%q) = true, want false", op)
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
			t.Errorf("IsPredicateOperation(%q) = false, want true (first-order)", op)
		}
	}

	secondOrder := []string{"and", "or", "not"}
	for _, op := range secondOrder {
		if !IsPredicateOperation(op) {
			t.Errorf("IsPredicateOperation(%q) = false, want true (second-order)", op)
		}
	}

	invalid := []string{"", "add", "remove", "inc", "flip", "unknown"}
	for _, op := range invalid {
		if IsPredicateOperation(op) {
			t.Errorf("IsPredicateOperation(%q) = true, want false", op)
		}
	}
}

func TestIsJSONPatchExtendedOperation(t *testing.T) {
	t.Parallel()
	valid := []string{"str_ins", "str_del", "flip", "inc", "split", "merge", "extend"}
	for _, op := range valid {
		if !IsJSONPatchExtendedOperation(op) {
			t.Errorf("IsJSONPatchExtendedOperation(%q) = false, want true", op)
		}
	}

	invalid := []string{"", "add", "test", "and", "defined", "unknown"}
	for _, op := range invalid {
		if IsJSONPatchExtendedOperation(op) {
			t.Errorf("IsJSONPatchExtendedOperation(%q) = true, want false", op)
		}
	}
}
