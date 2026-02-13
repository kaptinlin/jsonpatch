package op

import "github.com/kaptinlin/jsonpatch/internal"

// MoveOperation represents a move operation that moves a value from one path to another.
type MoveOperation struct {
	BaseOp
}

// NewMove creates a new move operation.
func NewMove(path, from []string) *MoveOperation {
	return &MoveOperation{
		BaseOp: NewBaseOpWithFrom(path, from),
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

// Apply applies the move operation following RFC 6902: remove then add.
func (m *MoveOperation) Apply(doc any) (internal.OpResult[any], error) {
	if pathEquals(m.path, m.from) {
		return internal.OpResult[any]{Doc: doc, Old: nil}, nil
	}

	if isPrefix(m.from, m.path) {
		return internal.OpResult[any]{}, ErrCannotMoveIntoChildren
	}

	// Move = remove + add (RFC 6902)
	removeOp := NewRemove(m.from)
	removeResult, err := removeOp.Apply(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	addOp := NewAdd(m.path, removeResult.Old)
	return addOp.Apply(removeResult.Doc)
}

// isPrefix checks if prefix is a prefix of path.
func isPrefix(prefix, path []string) bool {
	if len(prefix) >= len(path) {
		return false
	}
	for i, p := range prefix {
		if path[i] != p {
			return false
		}
	}
	return true
}

// ToJSON serializes the operation to JSON format.
func (m *MoveOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpMoveType),
		Path: formatPath(m.path),
		From: formatPath(m.from),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (m *MoveOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMoveCode, m.path, m.from}, nil
}

// Validate validates the move operation.
func (m *MoveOperation) Validate() error {
	if len(m.path) == 0 {
		return ErrPathEmpty
	}
	if len(m.from) == 0 {
		return ErrFromPathEmpty
	}
	if pathEquals(m.path, m.from) {
		return ErrPathsIdentical
	}
	if isPrefix(m.from, m.path) {
		return ErrCannotMoveIntoChildren
	}
	return nil
}
