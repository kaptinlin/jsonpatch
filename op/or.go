package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpOrOperation represents an OR operation that combines multiple predicate operations.
type OpOrOperation struct {
	BaseOp
	Operations []interface{} `json:"ops"` // Array of operations
}

// NewOpOrOperation creates a new OpOrOperation operation.
func NewOpOrOperation(path []string, ops []interface{}) *OpOrOperation {
	return &OpOrOperation{
		BaseOp:     NewBaseOp(path),
		Operations: ops,
	}
}

// Op returns the operation type.
func (o *OpOrOperation) Op() internal.OpType {
	return internal.OpOrType
}

// Code returns the operation code.
func (o *OpOrOperation) Code() int {
	return internal.OpOrCode
}

// Ops returns the predicate operations.
func (o *OpOrOperation) Ops() []internal.PredicateOp {
	ops := make([]internal.PredicateOp, 0, len(o.Operations))
	for _, op := range o.Operations {
		if predicateOp, ok := op.(internal.PredicateOp); ok {
			ops = append(ops, predicateOp)
		}
	}
	return ops
}

// Test performs the OR operation.
func (o *OpOrOperation) Test(doc interface{}) (bool, error) {
	// If no operations, return false (empty OR is false)
	if len(o.Operations) == 0 {
		return false, nil
	}

	// Test all operations
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return false, ErrInvalidPredicateInOr
		}
		ok, err := predicateOp.Test(doc)
		if err != nil {
			return false, err
		}
		if ok {
			// If any operation passes, OR passes
			return true, nil
		}
	}

	// No operations passed
	return false, nil
}

// Apply applies the OR operation to the document.
func (o *OpOrOperation) Apply(doc any) (internal.OpResult[any], error) {
	for _, predicateInterface := range o.Operations {
		predicate, ok := predicateInterface.(internal.PredicateOp)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidPredicateInOr
		}
		result, err := predicate.Apply(doc)
		if err == nil {
			return result, nil
		}
	}
	return internal.OpResult[any]{}, ErrOrTestFailed
}

// ToJSON serializes the operation to JSON format.
func (o *OpOrOperation) ToJSON() (internal.Operation, error) {
	opsJSON := make([]interface{}, 0, len(o.Operations))
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, ErrInvalidPredicateInOr
		}
		jsonVal, err := predicateOp.ToJSON()
		if err != nil {
			return nil, err
		}
		opsJSON = append(opsJSON, jsonVal)
	}
	return internal.Operation{
		"op":    string(internal.OpOrType),
		"path":  formatPath(o.Path()),
		"apply": opsJSON,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpOrOperation) ToCompact() (internal.CompactOperation, error) {
	opsCompact := make([]interface{}, 0, len(o.Operations))
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, ErrInvalidPredicateInOr
		}
		compact, err := predicateOp.ToCompact()
		if err != nil {
			return nil, err
		}
		opsCompact = append(opsCompact, compact)
	}
	return internal.CompactOperation{internal.OpOrCode, o.Path(), opsCompact}, nil
}

// Validate validates the OR operation.
func (o *OpOrOperation) Validate() error {
	if len(o.Operations) == 0 {
		return ErrOrNoOperands
	}
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return ErrInvalidPredicateInOr
		}
		if err := predicateOp.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Path returns the path for the OR operation.
func (o *OpOrOperation) Path() []string {
	return o.path
}

// Not returns false since this is not a NOT operation.
func (o *OpOrOperation) Not() bool {
	return false
}
