package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpMoveOperation represents a move operation that moves a value from one path to another.
type OpMoveOperation struct {
	BaseOp
	FromPath []string `json:"from"` // Source path
}

// NewOpMoveOperation creates a new OpMoveOperation operation.
func NewOpMoveOperation(path, from []string) *OpMoveOperation {
	return &OpMoveOperation{
		BaseOp:   NewBaseOpWithFrom(path, from),
		FromPath: from,
	}
}

// Op returns the operation type.
func (o *OpMoveOperation) Op() internal.OpType {
	return internal.OpMoveType
}

// Code returns the operation code.
func (o *OpMoveOperation) Code() int {
	return internal.OpMoveCode
}

// From returns the source path.
func (o *OpMoveOperation) From() []string {
	return o.FromPath
}

// Apply applies the move operation.
func (o *OpMoveOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Validate that from path exists
	if !pathExists(doc, o.FromPath) {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Check if trying to move a parent into its own child
	if isPrefix(o.FromPath, o.Path()) {
		return internal.OpResult[any]{}, ErrCannotMoveIntoChildren
	}

	// Check if this is a move within the same array
	if len(o.FromPath) > 1 && len(o.Path()) > 1 &&
		len(o.FromPath) == len(o.Path()) &&
		pathEquals(o.FromPath[:len(o.FromPath)-1], o.Path()[:len(o.Path())-1]) {
		// Check if the parent is actually an array
		arrayPath := o.FromPath[:len(o.FromPath)-1]
		if parent, err := getValue(doc, arrayPath); err == nil {
			if _, ok := parent.([]interface{}); ok {
				// Same array movement - use special logic
				return o.applySameArrayMove(doc)
			}
		}
	}

	// Get the value to move
	value, err := getValue(doc, o.FromPath)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Get the old value at the target path (what will be replaced)
	var oldValue interface{}

	// For array insertion, we need to handle the old value specially
	if len(o.Path()) > 0 {
		lastKey := o.Path()[len(o.Path())-1]
		if index, err := parseArrayIndex(lastKey); err == nil {
			// This is targeting an array index
			if len(o.Path()) > 1 {
				parentPath := o.Path()[:len(o.Path())-1]
				if parent, err := getValue(doc, parentPath); err == nil {
					if slice, ok := parent.([]interface{}); ok {
						// If inserting at a valid index within bounds, get the displaced element
						if index >= 0 && index < len(slice) {
							oldValue = slice[index]
						}
						// If appending (index == len), oldValue remains nil
					}
				}
			}
		} else {
			// Regular path, check if it exists
			if pathExists(doc, o.Path()) {
				oldValue, _ = getValue(doc, o.Path())
			}
		}
	}

	// Delete from source path first
	if len(o.FromPath) > 0 {
		parent, key, err := navigateToParent(doc, o.FromPath)
		if err != nil {
			return internal.OpResult[any]{}, err
		}

		// Handle array deletion specially
		if slice, ok := parent.([]interface{}); ok {
			if index, ok := key.(int); ok {
				if index < 0 || index >= len(slice) {
					return internal.OpResult[any]{}, fmt.Errorf("%w: index %d out of range", ErrIndexOutOfRange, index)
				}
				// Create a new slice without the element
				newSlice := make([]interface{}, 0, len(slice)-1)
				newSlice = append(newSlice, slice[:index]...)
				newSlice = append(newSlice, slice[index+1:]...)

				// We need to replace the array in its parent context
				// For now, let's handle root case and nested case differently
				if len(o.FromPath) == 1 {
					// This is modifying the root array, but we can't change doc directly
					return internal.OpResult[any]{}, ErrCannotModifyRootArray
				}
				// Get grandparent and update the parent
				grandParentPath := o.FromPath[:len(o.FromPath)-2]
				grandParentKey := o.FromPath[len(o.FromPath)-2]
				if len(grandParentPath) == 0 {
					// Parent is in root
					if docMap, ok := doc.(map[string]interface{}); ok {
						docMap[grandParentKey] = newSlice
					} else {
						return internal.OpResult[any]{}, ErrCannotUpdateParent
					}
				} else {
					grandParent, err := getValue(doc, grandParentPath)
					if err != nil {
						return internal.OpResult[any]{}, err
					}
					if grandParentMap, ok := grandParent.(map[string]interface{}); ok {
						grandParentMap[grandParentKey] = newSlice
					} else {
						return internal.OpResult[any]{}, ErrCannotUpdateGrandparent
					}
				}
			} else {
				return internal.OpResult[any]{}, ErrInvalidKeyTypeSlice
			}
		} else {
			// Use deleteFromParent for maps
			err = deleteFromParent(parent, key)
			if err != nil {
				return internal.OpResult[any]{}, err
			}
		}
	}

	// Set at target path (use insert mode for cross-array moves)
	if len(o.Path()) == 0 {
		// Moving to root - replace entire document
		return internal.OpResult[any]{Doc: value, Old: doc}, nil
	}
	err = insertValueAtPath(doc, o.Path(), value)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// applySameArrayMove handles movement within the same array
func (o *OpMoveOperation) applySameArrayMove(doc any) (internal.OpResult[any], error) {
	// Parse indices
	fromIndex, err := parseArrayIndex(o.FromPath[len(o.FromPath)-1])
	if err != nil {
		return internal.OpResult[any]{}, fmt.Errorf("%w: invalid from index: %w", ErrInvalidIndex, err)
	}

	toIndex, err := parseArrayIndex(o.Path()[len(o.Path())-1])
	if err != nil {
		return internal.OpResult[any]{}, fmt.Errorf("%w: invalid to index: %w", ErrInvalidIndex, err)
	}

	// Get the array directly
	arrayPath := o.FromPath[:len(o.FromPath)-1]
	array, err := getValue(doc, arrayPath)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	slice, ok := array.([]interface{})
	if !ok {
		return internal.OpResult[any]{}, ErrNotAnArray
	}

	// Validate indices
	if fromIndex < 0 || fromIndex >= len(slice) {
		return internal.OpResult[any]{}, ErrIndexOutOfRange
	}
	if toIndex < 0 || toIndex >= len(slice) {
		return internal.OpResult[any]{}, ErrIndexOutOfRange
	}

	// Get the value to move and old value at target
	value := slice[fromIndex]
	oldValue := slice[toIndex]

	// Create new array with the move using proper algorithm
	newSlice := make([]interface{}, 0, len(slice))

	if fromIndex < toIndex {
		// Moving forward: copy elements before fromIndex, skip fromIndex,
		// copy until toIndex, insert value, copy rest
		newSlice = append(newSlice, slice[:fromIndex]...)
		newSlice = append(newSlice, slice[fromIndex+1:toIndex+1]...)
		newSlice = append(newSlice, value)
		newSlice = append(newSlice, slice[toIndex+1:]...)
	} else {
		// Moving backward: copy elements before toIndex, insert value,
		// copy until fromIndex (skip fromIndex), copy rest
		newSlice = append(newSlice, slice[:toIndex]...)
		newSlice = append(newSlice, value)
		newSlice = append(newSlice, slice[toIndex:fromIndex]...)
		newSlice = append(newSlice, slice[fromIndex+1:]...)
	}

	// Update the array in its parent context
	if len(arrayPath) == 0 {
		// This would be modifying the root array, but we can't change doc directly
		return internal.OpResult[any]{}, ErrCannotModifyRootArray
	}
	// Get the parent of the array and update it
	arrayParent, arrayKey, err := navigateToParent(doc, arrayPath)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if mapParent, ok := arrayParent.(map[string]interface{}); ok {
		if keyStr, ok := arrayKey.(string); ok {
			mapParent[keyStr] = newSlice
		} else {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeMap
		}
	} else {
		return internal.OpResult[any]{}, ErrUnsupportedParentType
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
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
func (o *OpMoveOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":   string(internal.OpMoveType),
		"path": formatPath(o.Path()),
		"from": formatPath(o.FromPath),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpMoveOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMoveCode, o.Path(), o.FromPath}, nil
}

// Validate validates the move operation.
func (o *OpMoveOperation) Validate() error {
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
