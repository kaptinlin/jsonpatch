package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// MoveOperation represents a move operation that moves a value from one path to another.
type MoveOperation struct {
	BaseOp
	FromPath []string `json:"from"` // Source path
}

// NewOpMoveOperation creates a new OpMoveOperation operation.
func NewOpMoveOperation(path, from []string) *MoveOperation {
	return &MoveOperation{
		BaseOp:   NewBaseOpWithFrom(path, from),
		FromPath: from,
	}
}

// Op returns the operation type.
func (o *MoveOperation) Op() internal.OpType {
	return internal.OpMoveType
}

// Code returns the operation code.
func (o *MoveOperation) Code() int {
	return internal.OpMoveCode
}

// From returns the source path.
func (o *MoveOperation) From() []string {
	return o.FromPath
}

// Apply applies the move operation following RFC 6902: remove then add.
func (o *MoveOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Check that path and from are not the same
	if pathEquals(o.Path(), o.FromPath) {
		// Moving to the same location has no effect according to JSON Patch spec
		return internal.OpResult[any]{Doc: doc, Old: nil}, nil
	}

	// Check if trying to move a parent into its own child
	if isPrefix(o.FromPath, o.Path()) {
		return internal.OpResult[any]{}, ErrCannotMoveIntoChildren
	}

	// Following TypeScript reference: move = remove + add
	// First, remove from source path
	removeOp := NewOpRemoveOperation(o.FromPath)
	removeResult, err := removeOp.Apply(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Then, add the removed value to target path
	addOp := NewAdd(o.Path(), removeResult.Old)
	addResult, err := addOp.Apply(removeResult.Doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return addResult, nil
}

// isPrefix checks if prefix is a prefix of path
func isPrefix(prefix, path []string) bool {
	if len(prefix) >= len(path) {
		return false
	}
	for i, p := range prefix {
		if i >= len(path) || path[i] != p {
			return false
		}
	}
	return true
}

// ToJSON serializes the operation to JSON format.
func (o *MoveOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":   string(internal.OpMoveType),
		"path": formatPath(o.Path()),
		"from": formatPath(o.FromPath),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *MoveOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMoveCode, o.Path(), o.FromPath}, nil
}

// Validate validates the move operation.
func (o *MoveOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(o.FromPath) == 0 {
		return ErrFromPathEmpty
	}
	// Check that path and from are not the same
	if pathEquals(o.Path(), o.FromPath) {
		return ErrPathsIdentical
	}
	// Check for moving into own children
	if isPrefix(o.FromPath, o.Path()) {
		return ErrCannotMoveIntoChildren
	}
	return nil
}

// Short aliases for common use
var (
	// NewMove creates a new move operation
	NewMove = NewOpMoveOperation
)
