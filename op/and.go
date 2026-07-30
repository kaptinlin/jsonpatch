package op

import (
	"slices"

	"github.com/kaptinlin/jsonpatch/internal"
)

// AndOperation represents an and operation that combines multiple predicate operations.
type AndOperation struct {
	BaseOp
	Operations []PredicateOp `json:"apply"`
}

// NewAnd creates a new AND operation.
func NewAnd(path []string, predicates ...PredicateOp) *AndOperation {
	return &AndOperation{
		BaseOp:     NewBaseOp(path),
		Operations: slices.Clone(predicates),
	}
}

// Op returns the operation type.
func (ao *AndOperation) Op() internal.OpType {
	return internal.OpAndType
}

// Ops returns the predicate operations.
func (ao *AndOperation) Ops() []internal.PredicateOp {
	return slices.Clone(ao.Operations)
}

// Test performs the AND operation.
func (ao *AndOperation) Test(doc any) (bool, error) {
	// If no operations, return true (vacuous truth - empty AND is true)
	if len(ao.Operations) == 0 {
		return true, nil
	}

	for _, predicate := range ao.Operations {
		if isNilPredicate(predicate) {
			return false, ErrInvalidPredicateInAnd
		}
		ok, err := predicate.Test(doc)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

// Apply applies the AND operation.
func (ao *AndOperation) Apply(doc any) (internal.OpResult[any], error) {
	ok, err := ao.Test(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	if !ok {
		return internal.OpResult[any]{}, predicateFailure(ErrAndTestFailed)
	}
	return internal.OpResult[any]{Doc: doc}, nil
}

// Validate validates the AND operation.
func (ao *AndOperation) Validate() error {
	return validatePredicateOps(ao.Operations, ErrInvalidPredicateInAnd)
}
