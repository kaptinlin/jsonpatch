package jsonpatch

import (
	"regexp"

	"github.com/kaptinlin/jsonpatch/internal"
)

// CreateMatcherDefault creates a regex matcher from a pattern and case sensitivity flag.
// Returns a matcher that always returns false if pattern compilation fails.
// This aligns with json-joy's createMatcherDefault behavior.
func CreateMatcherDefault(pattern string, ignoreCase bool) RegexMatcher {
	flags := ""
	if ignoreCase {
		flags = "(?i)"
	}

	regex, err := regexp.Compile(flags + pattern)
	if err != nil {
		// Return a matcher that always returns false if compilation fails
		return func(_ string) bool { return false }
	}

	return regex.MatchString
}

// ValidateOp validates an operation.
func ValidateOp(op internal.Op) error {
	return op.Validate()
}

// GetOpType returns the operation type.
func GetOpType(op internal.Op) internal.OpType {
	return op.Op()
}

// GetOpCode returns the operation code.
func GetOpCode(op internal.Op) int {
	return op.Code()
}

// GetOpPath returns the operation path.
func GetOpPath(op internal.Op) []string {
	return op.Path()
}

// ApplyOpDirect applies an operation directly to a document.
func ApplyOpDirect(op internal.Op, doc any) (internal.OpResult[any], error) {
	return op.Apply(doc)
}

// ToJSON converts an operation to JSON format.
func ToJSON(op internal.Op) (internal.Operation, error) {
	return op.ToJSON()
}

// ToCompact converts an operation to compact format.
func ToCompact(op internal.Op) (internal.CompactOperation, error) {
	return op.ToCompact()
}

// TestPredicate tests a predicate operation against a document.
func TestPredicate(predicate internal.PredicateOp, doc any) (bool, error) {
	return predicate.Test(doc)
}

// IsNotPredicate reports whether a predicate is negated.
func IsNotPredicate(predicate internal.PredicateOp) bool {
	return predicate.Not()
}

// GetSecondOrderOps returns sub-operations from a second-order predicate.
func GetSecondOrderOps(predicate internal.SecondOrderPredicateOp) []internal.PredicateOp {
	return predicate.Ops()
}
