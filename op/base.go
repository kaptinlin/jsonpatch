package op

// BaseOp provides common functionality for all operations.
// It stores the target path and optional source path (for move/copy operations).
// All operation types embed BaseOp to inherit path management functionality.
type BaseOp struct {
	path []string
	from []string
}

// NewBaseOp creates a new BaseOp with the given path.
// The path is a JSON Pointer path represented as a slice of string segments.
func NewBaseOp(path []string) BaseOp {
	return BaseOp{path: path}
}

// NewBaseOpWithFrom creates a new BaseOp with both target path and source path.
// This is used for move and copy operations that require both paths.
func NewBaseOpWithFrom(path, from []string) BaseOp {
	return BaseOp{path: path, from: from}
}

// Path returns the operation path.
func (b *BaseOp) Path() []string {
	return b.path
}

// From returns the from path for move/copy operations.
func (b *BaseOp) From() []string {
	return b.from
}

// HasFrom returns true if the operation has a from path.
func (b *BaseOp) HasFrom() bool {
	return len(b.from) > 0
}

// Not returns false by default. Operations that support negation override this.
func (b *BaseOp) Not() bool {
	return false
}
