package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpUndefinedOperation represents an undefined operation that checks if a path doesn't exist.
type OpUndefinedOperation struct {
	PredicateOpBase
}

// NewOpUndefinedOperation creates a new undefined operation.
func NewOpUndefinedOperation(path []string, not bool) *OpUndefinedOperation {
	return &OpUndefinedOperation{
		PredicateOpBase: NewPredicateOp(path, not),
	}
}

// Op returns the operation type.
func (o *OpUndefinedOperation) Op() internal.OpType {
	return internal.OpUndefinedType
}

// Code returns the operation code.
func (o *OpUndefinedOperation) Code() int {
	return internal.OpUndefinedCode
}

// checkPathUndefined is a helper function that checks if a path is undefined
func (o *OpUndefinedOperation) checkPathUndefined(doc interface{}) bool {
	_, err := getValue(doc, o.path)
	// Path doesn't exist means undefined is true
	result := err != nil
	if o.not {
		result = !result
	}
	return result
}

// Test performs the undefined operation.
func (o *OpUndefinedOperation) Test(doc interface{}) (bool, error) {
	return o.checkPathUndefined(doc), nil
}

// Not returns false (undefined operation doesn't support not modifier).
func (o *OpUndefinedOperation) Not() bool {
	return false
}

// Apply applies the undefined operation.
func (o *OpUndefinedOperation) Apply(doc any) (internal.OpResult, error) {
	if !o.checkPathUndefined(doc) {
		return internal.OpResult{}, ErrUndefinedTestFailed
	}
	return internal.OpResult{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *OpUndefinedOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":   string(internal.OpUndefinedType),
		"path": formatPath(o.path),
	}
	if o.not {
		result["not"] = true
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpUndefinedOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpUndefinedCode, o.path, o.not}, nil
}

// Validate validates the undefined operation.
func (o *OpUndefinedOperation) Validate() error {
	if len(o.path) == 0 {
		return ErrPathEmpty
	}
	return nil
}

func (o *OpUndefinedOperation) Path() []string {
	return o.path
}
