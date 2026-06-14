package op

import (
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// StrInsOperation represents a string insert operation.
// path: target path
// pos: insert position (rune index)
// str: string to insert
// Only supports string type fields.
type StrInsOperation struct {
	BaseOp
	Pos int    `json:"pos"` // Insert position
	Str string `json:"str"` // String to insert
}

// NewStrIns creates a new string insert operation.
func NewStrIns(path []string, pos float64, str string) *StrInsOperation {
	return &StrInsOperation{
		BaseOp: NewBaseOp(path),
		Pos:    int(pos),
		Str:    str,
	}
}

// Op returns the operation type.
func (si *StrInsOperation) Op() internal.OpType {
	return internal.OpStrInsType
}

// Apply applies the string insert operation.
func (si *StrInsOperation) Apply(doc any) (internal.OpResult[any], error) {
	path := si.Path()
	target, err := value(doc, path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	targetStr, ok := target.(string)
	if target == nil {
		if si.Pos != 0 {
			return internal.OpResult[any]{}, ErrInvalidPosition
		}
	} else if !ok {
		return internal.OpResult[any]{}, ErrNotString
	}

	result := si.applyStrIns(targetStr)
	if len(path) == 0 {
		return internal.OpResult[any]{Doc: result, Old: target}, nil
	}

	if err := setValueAtPath(doc, path, result); err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

func clampStringPosition(pos, length int) int {
	pos = min(pos, length)
	if pos < 0 {
		return max(length+pos, 0)
	}
	return pos
}

func (si *StrInsOperation) applyStrIns(str string) string {
	runes := []rune(str)
	pos := clampStringPosition(si.Pos, len(runes))

	var builder strings.Builder
	builder.Grow(len(str) + len(si.Str))

	if pos > 0 {
		builder.WriteString(string(runes[:pos]))
	}
	builder.WriteString(si.Str)
	if pos < len(runes) {
		builder.WriteString(string(runes[pos:]))
	}

	return builder.String()
}

// Validate validates the string insert operation.
// Negative positions are valid (JS slice semantics: count from end).
func (si *StrInsOperation) Validate() error {
	return nil
}
