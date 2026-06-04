package op

import (
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// StrDelOperation represents a string delete operation.
// path: target path
// pos: start position (rune index)
// len: number of runes to delete (when HasStr is false)
// str: specific string to delete (when HasStr is true, takes precedence)
// Only supports string type fields.
type StrDelOperation struct {
	BaseOp
	Pos    int    `json:"pos"` // Delete position
	Len    int    `json:"len"` // Number of characters to delete
	Str    string `json:"str"` // Specific string to delete (optional)
	HasStr bool   // true when str mode is explicitly set (distinguishes "" from unset)
}

// NewStrDel creates a new string delete operation with length.
func NewStrDel(path []string, pos, length float64) *StrDelOperation {
	return &StrDelOperation{
		BaseOp: NewBaseOp(path),
		Pos:    int(pos),
		Len:    int(length),
		Str:    "",
	}
}

// NewStrDelWithStr creates a new string delete operation with specific string.
func NewStrDelWithStr(path []string, pos float64, str string) *StrDelOperation {
	return &StrDelOperation{
		BaseOp: NewBaseOp(path),
		Pos:    int(pos),
		Len:    len([]rune(str)),
		Str:    str,
		HasStr: true,
	}
}

// Op returns the operation type.
func (sd *StrDelOperation) Op() internal.OpType {
	return internal.OpStrDelType
}

// Code returns the operation code.
func (sd *StrDelOperation) Code() int {
	return internal.OpStrDelCode
}

// Apply applies the string delete operation.
func (sd *StrDelOperation) Apply(doc any) (internal.OpResult[any], error) {
	path := sd.Path()
	target, err := value(doc, path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	targetStr, ok := target.(string)
	if !ok {
		return internal.OpResult[any]{}, ErrNotString
	}

	result := sd.applyStrDel(targetStr)
	if len(path) == 0 {
		return internal.OpResult[any]{Doc: result, Old: target}, nil
	}

	if err := setValueAtPath(doc, path, result); err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

func (sd *StrDelOperation) applyStrDel(val string) string {
	runes := []rune(val)
	length := len(runes)
	pos := clampStringPosition(sd.Pos, length)

	deletionLength := sd.Len
	if sd.HasStr {
		deletionLength = len([]rune(sd.Str))
	}
	if deletionLength <= 0 {
		return val
	}

	end := min(pos+deletionLength, length)
	if pos >= length || pos == end {
		return val
	}

	var builder strings.Builder
	builder.Grow(len(val) - (end - pos))

	if pos > 0 {
		builder.WriteString(string(runes[:pos]))
	}
	if end < length {
		builder.WriteString(string(runes[end:]))
	}

	return builder.String()
}

// ToJSON serializes the operation to JSON format.
func (sd *StrDelOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpStrDelType),
		Path: formatPath(sd.Path()),
		Pos:  sd.Pos,
		Len:  sd.Len,
	}

	if sd.HasStr {
		result.Str = sd.Str
		result.Len = 0
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (sd *StrDelOperation) ToCompact() (internal.CompactOperation, error) {
	if sd.HasStr {
		return internal.CompactOperation{internal.OpStrDelCode, sd.Path(), sd.Pos, sd.Str}, nil
	}
	return internal.CompactOperation{internal.OpStrDelCode, sd.Path(), sd.Pos, 0, sd.Len}, nil
}

// Validate validates the string delete operation.
// Negative positions are valid (JS slice semantics: count from end).
func (sd *StrDelOperation) Validate() error {
	if sd.Len < 0 {
		return ErrLengthNegative
	}
	if !sd.HasStr && sd.Len <= 0 {
		return ErrMissingStrOrLen
	}
	return nil
}
