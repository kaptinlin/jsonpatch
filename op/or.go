package op

import "github.com/kaptinlin/jsonpatch/internal"

// OrOperation represents an OR operation that combines multiple predicate operations.
type OrOperation struct {
	BaseOp
	Operations []any `json:"apply"` // Array of operations to apply
}

// NewOr creates a new OR operation.
func NewOr(path []string, ops []any) *OrOperation {
	return &OrOperation{
		BaseOp:     NewBaseOp(path),
		Operations: ops,
	}
}

// Op returns the operation type.
func (oo *OrOperation) Op() internal.OpType {
	return internal.OpOrType
}

// Code returns the operation code.
func (oo *OrOperation) Code() int {
	return internal.OpOrCode
}

// Ops returns the predicate operations.
func (oo *OrOperation) Ops() []internal.PredicateOp {
	return extractPredicateOps(oo.Operations)
}

// Test performs the OR operation.
func (oo *OrOperation) Test(doc any) (bool, error) {
	// If no operations, return false (empty OR is false)
	if len(oo.Operations) == 0 {
		return false, nil
	}

	// Test all operations
	for _, op := range oo.Operations {
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
func (oo *OrOperation) Apply(doc any) (internal.OpResult[any], error) {
	for _, predicateInterface := range oo.Operations {
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
func (oo *OrOperation) ToJSON() (internal.Operation, error) {
	opsJSON := make([]internal.Operation, 0, len(oo.Operations))
	for _, op := range oo.Operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return internal.Operation{}, ErrInvalidPredicateInOr
		}
		jsonVal, err := predicateOp.ToJSON()
		if err != nil {
			return internal.Operation{}, err
		}
		opsJSON = append(opsJSON, jsonVal)
	}
	return internal.Operation{
		Op:    string(internal.OpOrType),
		Path:  formatPath(oo.Path()),
		Apply: opsJSON,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (oo *OrOperation) ToCompact() (internal.CompactOperation, error) {
	opsCompact := make([]any, 0, len(oo.Operations))
	for _, op := range oo.Operations {
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
	return internal.CompactOperation{internal.OpOrCode, oo.Path(), opsCompact}, nil
}

// Validate validates the OR operation.
func (oo *OrOperation) Validate() error {
	// Empty operations are valid (though they return false)
	for _, op := range oo.Operations {
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
