package op

import "github.com/kaptinlin/jsonpatch/internal"

// AndOperation represents an and operation that combines multiple predicate operations.
type AndOperation struct {
	BaseOp
	Operations []interface{} `json:"apply"` // Array of operations to apply
}

// NewAnd creates a new AND operation.
func NewAnd(path []string, ops []interface{}) *AndOperation {
	return &AndOperation{
		BaseOp:     NewBaseOp(path),
		Operations: ops,
	}
}

// Op returns the operation type.
func (o *AndOperation) Op() internal.OpType {
	return internal.OpAndType
}

// Code returns the operation code.
func (o *AndOperation) Code() int {
	return internal.OpAndCode
}

// Ops returns the predicate operations.
func (o *AndOperation) Ops() []internal.PredicateOp {
	return extractPredicateOps(o.Operations)
}

// Test performs the AND operation.
func (o *AndOperation) Test(doc interface{}) (bool, error) {
	// If no operations, return true (vacuous truth - empty AND is true)
	if len(o.Operations) == 0 {
		return true, nil
	}

	// Test all operations
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return false, ErrInvalidPredicateInAnd
		}
		ok, err := predicateOp.Test(doc)
		if err != nil {
			return false, err
		}
		if !ok {
			// If any operation fails, AND fails
			return false, nil
		}
	}

	// All operations passed
	return true, nil
}

// Apply applies the AND operation.
func (o *AndOperation) Apply(doc any) (internal.OpResult[any], error) {
	ok, err := o.Test(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	if !ok {
		return internal.OpResult[any]{}, ErrAndTestFailed
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *AndOperation) ToJSON() (internal.Operation, error) {
	operations := make([]internal.Operation, 0, len(o.Operations))
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return internal.Operation{}, ErrInvalidPredicateInAnd
		}
		jsonVal, err := predicateOp.ToJSON()
		if err != nil {
			return internal.Operation{}, err
		}
		operations = append(operations, jsonVal)
	}
	return internal.Operation{
		Op:    string(internal.OpAndType),
		Path:  formatPath(o.path),
		Apply: operations,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *AndOperation) ToCompact() (internal.CompactOperation, error) {
	opsCompact := make([]interface{}, 0, len(o.Operations))
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, ErrInvalidPredicateInAnd
		}
		compact, err := predicateOp.ToCompact()
		if err != nil {
			return nil, err
		}
		opsCompact = append(opsCompact, compact)
	}
	return internal.CompactOperation{internal.OpAndCode, o.path, opsCompact}, nil
}

// Validate validates the AND operation.
func (o *AndOperation) Validate() error {
	// Empty operations are valid (vacuous truth)
	for _, op := range o.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return ErrInvalidPredicateInAnd
		}
		if err := predicateOp.Validate(); err != nil {
			return err
		}
	}
	return nil
}

