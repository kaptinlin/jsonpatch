package op

import "github.com/kaptinlin/jsonpatch/internal"

// Op interface defines the core operation behavior.
type Op = internal.Op

// Result represents the result of applying an operation.
type Result[T internal.Document] = internal.OpResult[T]

// PredicateOp represents predicate operations used for testing conditions.
type PredicateOp = internal.PredicateOp

// SecondOrderPredicateOp represents operations that combine multiple predicate operations.
type SecondOrderPredicateOp = internal.SecondOrderPredicateOp
