package internal

// Op is the unified interface for all JSON Patch operations.
// Any operation must implement these methods.
type Op interface {
	// Op returns the operation type string (e.g. "add", "remove", etc.)
	Op() OpType
	// Code returns the numeric code for the operation (consistent with constants)
	Code() int
	// Path returns the JSON Pointer path for the operation (slice form)
	Path() []string
	// Apply applies the operation to the document, returning the new document and old value
	Apply(doc any) (OpResult[any], error)
	// ToJSON returns the standard JSON Patch format
	ToJSON() (Operation, error)
	// ToCompact returns the compact array format
	ToCompact() (CompactOperation, error)
	// Validate validates the operation parameters
	Validate() error
}

// PredicateOp is the interface for all predicate (test-type) operations.
type PredicateOp interface {
	Op
	// Test tests the operation on the document, returning whether it passed
	Test(doc any) (bool, error)
	// Not returns whether the operation is a negation predicate
	Not() bool
}

// SecondOrderPredicateOp is the interface for second-order predicate operations that combine multiple predicates.
type SecondOrderPredicateOp interface {
	PredicateOp
	// Ops returns all sub-predicate operations
	Ops() []PredicateOp
}

// Codec is the interface for codecs that provide unified encoding and decoding functionality.
type Codec interface {
	// Decode decodes a JSON operation object into an Op instance
	Decode(operation Operation) (Op, error)
	// DecodeSlice decodes an array of operations into an Op slice
	DecodeSlice(operations []Operation) ([]Op, error)
	// Encode encodes an Op instance into a JSON operation object
	Encode(op Op) (Operation, error)
	// EncodeSlice encodes an Op slice into an array of operations
	EncodeSlice(ops []Op) ([]Operation, error)
}
