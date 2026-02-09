package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// MoveOperation represents a move operation that moves a value from one path to another.
type MoveOperation struct {
	BaseOp
	FromPath []string `json:"from"` // Source path
}

// NewMove creates a new move operation.
func NewMove(path, from []string) *MoveOperation {
	return &MoveOperation{
		BaseOp:   NewBaseOpWithFrom(path, from),
		FromPath: from,
	}
}

// Op returns the operation type.
func (m *MoveOperation) Op() internal.OpType {
	return internal.OpMoveType
}

// Code returns the operation code.
func (m *MoveOperation) Code() int {
	return internal.OpMoveCode
}

// From returns the source path.
func (m *MoveOperation) From() []string {
	return m.FromPath
}

// Apply applies the move operation following RFC 6902: remove then add.
func (m *MoveOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Check that path and from are not the same
	if pathEquals(m.Path(), m.FromPath) {
		// Moving to the same location has no effect according to JSON Patch spec
		return internal.OpResult[any]{Doc: doc, Old: nil}, nil
	}

	// Check if trying to move a parent into its own child
	if isPrefix(m.FromPath, m.Path()) {
		return internal.OpResult[any]{}, ErrCannotMoveIntoChildren
	}

	// Following TypeScript reference: move = remove + add
	// First, remove from source path
	removeOp := NewRemove(m.FromPath)
	removeResult, err := removeOp.Apply(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Then, add the removed value to target path
	addOp := NewAdd(m.Path(), removeResult.Old)
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
func (m *MoveOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpMoveType),
		Path: formatPath(m.Path()),
		From: formatPath(m.FromPath),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (m *MoveOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMoveCode, m.Path(), m.FromPath}, nil
}

// Validate validates the move operation.
func (m *MoveOperation) Validate() error {
	if len(m.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(m.FromPath) == 0 {
		return ErrFromPathEmpty
	}
	// Check that path and from are not the same
	if pathEquals(m.Path(), m.FromPath) {
		return ErrPathsIdentical
	}
	// Check for moving into own children
	if isPrefix(m.FromPath, m.Path()) {
		return ErrCannotMoveIntoChildren
	}
	return nil
}
