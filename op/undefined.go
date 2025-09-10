package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// UndefinedOperation represents an undefined operation that checks if a path doesn't exist.
type UndefinedOperation struct {
	PredicateOpBase
}

// NewOpUndefinedOperation creates a new undefined operation.
func NewOpUndefinedOperation(path []string, not bool) *UndefinedOperation {
	return &UndefinedOperation{
		PredicateOpBase: NewPredicateOp(path, not),
	}
}

// Op returns the operation type.
func (o *UndefinedOperation) Op() internal.OpType {
	return internal.OpUndefinedType
}

// Code returns the operation code.
func (o *UndefinedOperation) Code() int {
	return internal.OpUndefinedCode
}

// checkPathUndefined is a helper function that checks if a path is undefined
func (o *UndefinedOperation) checkPathUndefined(doc interface{}) bool {
	_, err := getValue(doc, o.path)
	// Path doesn't exist means undefined is true
	result := err != nil
	if o.not {
		result = !result
	}
	return result
}

// Test performs the undefined operation.
func (o *UndefinedOperation) Test(doc interface{}) (bool, error) {
	return o.checkPathUndefined(doc), nil
}

// Not returns false (undefined operation doesn't support not modifier).
func (o *UndefinedOperation) Not() bool {
	return false
}

// Apply applies the undefined operation.
func (o *UndefinedOperation) Apply(doc any) (internal.OpResult[any], error) {
	if !o.checkPathUndefined(doc) {
		return internal.OpResult[any]{}, ErrUndefinedTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *UndefinedOperation) ToJSON() (internal.Operation, error) {
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
func (o *UndefinedOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpUndefinedCode, o.path, o.not}, nil
}

// Validate validates the undefined operation.
func (o *UndefinedOperation) Validate() error {
	if len(o.path) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Path returns the path for the undefined operation.
func (o *UndefinedOperation) Path() []string {
	return o.path
}

// Short aliases for common use
var (
	// NewUndefined creates a new undefined operation
	NewUndefined = NewOpUndefinedOperation
)
