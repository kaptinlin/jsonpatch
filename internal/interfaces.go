package internal

// Op is the executable interface for JSON Patch operations.
type Op interface {
	// Op returns the operation type string (e.g. "add", "remove").
	Op() OpType
	// Path returns the JSON Pointer path as a string slice.
	Path() []string
	// Apply applies the operation to the document, returning the result and any error.
	Apply(doc any) (OpResult[any], error)
	// Validate checks that the operation parameters are valid.
	Validate() error
}

// CloneOp is an operation that can freeze itself for compiled patch storage.
type CloneOp interface {
	Op
	// Clone returns an independent executable copy of the operation.
	Clone() (Op, error)
}

// JSONOp is an operation that can project itself to JSON operation form.
type JSONOp interface {
	Op
	// ToJSON serializes the operation to standard JSON Patch format.
	ToJSON() (Operation, error)
}

// CompactOp is an operation that can project itself to compact array form.
type CompactOp interface {
	Op
	// Code returns the numeric code for the operation.
	Code() int
	// ToCompact serializes the operation to compact array format.
	ToCompact() (CompactOperation, error)
}

// PredicateOp is the interface for predicate (test-type) operations.
type PredicateOp interface {
	Op
	// Test evaluates the predicate against the document.
	Test(doc any) (bool, error)
	// Not reports whether the predicate is negated.
	Not() bool
}

// SecondOrderPredicateOp is the interface for composite predicates
// that combine multiple sub-predicates (and, or, not).
type SecondOrderPredicateOp interface {
	PredicateOp
	// Ops returns the sub-predicate operations.
	Ops() []PredicateOp
}

// Codec is the interface for encoding and decoding JSON Patch operations.
type Codec interface {
	// Decode decodes a single JSON operation into an Op.
	Decode(operation Operation) (Op, error)
	// DecodeSlice decodes a slice of JSON operations into Ops.
	DecodeSlice(operations []Operation) ([]Op, error)
	// Encode encodes an Op into a JSON operation.
	Encode(op Op) (Operation, error)
	// EncodeSlice encodes a slice of Ops into JSON operations.
	EncodeSlice(ops []Op) ([]Operation, error)
}
