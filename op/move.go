package op

import (
	"errors"
	"fmt"
	"slices"

	"github.com/kaptinlin/jsonpatch/internal"
)

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

// Apply applies the move operation following RFC 6902: remove then add.
func (m *MoveOperation) Apply(doc any) (internal.OpResult[any], error) {
	if slices.Equal(m.path, m.from) {
		if _, err := value(doc, m.from); err != nil {
			return internal.OpResult[any]{}, err
		}
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
	result, err := addOp.Apply(removeResult.Doc)
	if err == nil {
		return result, nil
	}
	if _, _, restoreErr := addAtPath(removeResult.Doc, m.from, removeResult.Old); restoreErr != nil {
		return internal.OpResult[any]{}, errors.Join(err, fmt.Errorf("restore move source: %w", restoreErr))
	}
	return internal.OpResult[any]{}, err
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

// Validate validates the move operation.
func (m *MoveOperation) Validate() error {
	if isPrefix(m.from, m.path) {
		return ErrCannotMoveIntoChildren
	}
	return nil
}
