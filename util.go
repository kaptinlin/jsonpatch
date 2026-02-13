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
func ValidateOp(operation internal.Op) error {
	return operation.Validate()
}

// GetOpType returns the operation type using the Op interface.
func GetOpType(operation internal.Op) internal.OpType {
	return operation.Op()
}

// GetOpCode returns the operation code using the Op interface.
func GetOpCode(operation internal.Op) int {
	return operation.Code()
}

// GetOpPath returns the operation path using the Op interface.
func GetOpPath(operation internal.Op) []string {
	return operation.Path()
}

// ApplyOpDirect applies an operation directly using the Op interface.
func ApplyOpDirect(operation internal.Op, doc any) (internal.OpResult[any], error) {
	return operation.Apply(doc)
}

// ToJSON converts an operation to JSON format using the Op interface.
func ToJSON(operation internal.Op) (internal.Operation, error) {
	return operation.ToJSON()
}

// ToCompact converts an operation to compact format using the Op interface.
func ToCompact(operation internal.Op) (internal.CompactOperation, error) {
	return operation.ToCompact()
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
