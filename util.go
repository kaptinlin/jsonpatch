package jsonpatch

import (
	"regexp"

	"github.com/kaptinlin/jsonpatch/internal"
)

// CreateMatcherDefault is the default regex matcher factory.
// It creates a RegexMatcher function from a pattern and ignoreCase flag.
// If the pattern is invalid, returns a matcher that always returns false.
// This aligns with json-joy's createMatcherDefault.
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

// ValidateOp validates an operation using the Op interface.
func ValidateOp(op internal.Op) error {
	return op.Validate()
}

// GetOpType returns the operation type using the Op interface.
func GetOpType(op internal.Op) internal.OpType {
	return op.Op()
}

// GetOpCode returns the operation code using the Op interface.
func GetOpCode(op internal.Op) int {
	return op.Code()
}

// GetOpPath returns the operation path using the Op interface.
func GetOpPath(op internal.Op) []string {
	return op.Path()
}

// ApplyOpDirect applies an operation directly using the Op interface.
func ApplyOpDirect(op internal.Op, doc any) (internal.OpResult[any], error) {
	return op.Apply(doc)
}

// ToJSON converts an operation to JSON format using the Op interface.
func ToJSON(op internal.Op) (internal.Operation, error) {
	return op.ToJSON()
}

// ToCompact converts an operation to compact format using the Op interface.
func ToCompact(op internal.Op) (internal.CompactOperation, error) {
	return op.ToCompact()
}

// TestPredicate tests a predicate operation using the PredicateOp interface.
func TestPredicate(predicate internal.PredicateOp, doc any) (bool, error) {
	return predicate.Test(doc)
}

// IsNotPredicate checks if a predicate is negated using the PredicateOp interface.
func IsNotPredicate(predicate internal.PredicateOp) bool {
	return predicate.Not()
}

// GetSecondOrderOps returns sub-operations from a second-order predicate using the SecondOrderPredicateOp interface.
func GetSecondOrderOps(predicate internal.SecondOrderPredicateOp) []internal.PredicateOp {
	return predicate.Ops()
}
