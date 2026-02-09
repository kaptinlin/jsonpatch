package op

import (
	"fmt"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestStringOperation represents a test operation that checks if a value is a string and matches a pattern.
type TestStringOperation struct {
	BaseOp
	Str        string `json:"str"`                   // Expected string value
	Pos        int    `json:"pos"`                   // Position within string
	NotFlag    bool   `json:"not,omitempty"`         // Whether to negate the result
	IgnoreCase bool   `json:"ignore_case,omitempty"` // Whether to ignore case
}

// NewTestString creates a new test string operation.
func NewTestString(path []string, expectedValue string) *TestStringOperation {
	return &TestStringOperation{
		BaseOp:  NewBaseOp(path),
		Str:     expectedValue,
		Pos:     0,     // Default position
		NotFlag: false, // Default not flag
	}
}

// NewTestStringWithPos creates a new test string operation with position.
func NewTestStringWithPos(path []string, expectedValue string, pos float64) *TestStringOperation {
	return &TestStringOperation{
		BaseOp:  NewBaseOp(path),
		Str:     expectedValue,
		Pos:     int(pos),
		NotFlag: false, // Default not flag
	}
}

// NewTestStringWithPosAndNot creates a new test string operation with position and not flag.
func NewTestStringWithPosAndNot(path []string, expectedValue string, pos float64, notFlag bool) *TestStringOperation {
	return &TestStringOperation{
		BaseOp:     NewBaseOp(path),
		Str:        expectedValue,
		Pos:        int(pos),
		NotFlag:    notFlag,
		IgnoreCase: false,
	}
}

// NewTestStringWithIgnoreCase creates a new test string operation with ignore case flag.
func NewTestStringWithIgnoreCase(path []string, expectedValue string, pos float64, notFlag bool, ignoreCase bool) *TestStringOperation {
	return &TestStringOperation{
		BaseOp:     NewBaseOp(path),
		Str:        expectedValue,
		Pos:        int(pos),
		NotFlag:    notFlag,
		IgnoreCase: ignoreCase,
	}
}

// Op returns the operation type.
func (o *TestStringOperation) Op() internal.OpType {
	return internal.OpTestStringType
}

// Code returns the operation code.
func (o *TestStringOperation) Code() int {
	return internal.OpTestStringCode
}

// Test evaluates the test string predicate condition.
func (o *TestStringOperation) Test(doc any) (bool, error) {
	val, err := getValue(doc, o.Path())
	if err != nil {
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	// Convert to string or from byte slice
	str, ok := extractString(val)
	if !ok {
		return false, nil // Return false if not string or byte slice
	}

	// Implement the same logic as json-joy reference
	// const length = (val as string).length;
	// const start = Math.min(this.pos, length);
	// const end = Math.min(this.pos + this.str.length, length);
	// const test = (val as string).substring(start, end) === this.str;
	// return this.not ? !test : test;

	length := len(str)
	start := min(o.Pos, length)
	end := min(o.Pos+len(o.Str), length)

	substring := str[start:end]
	var test bool
	if o.IgnoreCase {
		test = strings.EqualFold(substring, o.Str)
	} else {
		test = substring == o.Str
	}
	return o.NotFlag != test, nil // XOR with NotFlag for negation
}

// Apply applies the test string operation to the document.
func (o *TestStringOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get target value
	val, err := getValue(doc, o.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Check if value is a string or convert byte slice to string
	var str string
	switch v := val.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return internal.OpResult[any]{}, ErrNotString
	}

	pos := o.Pos
	if pos < 0 || pos > len(str) {
		return internal.OpResult[any]{}, ErrPositionOutOfStringRange
	}

	endPos := pos + len(o.Str)
	if endPos > len(str) {
		return internal.OpResult[any]{}, ErrSubstringTooLong
	}

	substring := str[pos:endPos]
	var matches bool
	if o.IgnoreCase {
		matches = strings.EqualFold(substring, o.Str)
	} else {
		matches = substring == o.Str
	}

	// Apply negation logic using XOR (same as Test method)
	shouldPass := matches != o.NotFlag
	if !shouldPass {
		if o.NotFlag {
			// When Not is true and test fails, it means the string DID match when we expected it not to
			return internal.OpResult[any]{}, fmt.Errorf("%w: string matched %q at position %d when NOT expected", ErrSubstringMismatch, o.Str, pos)
		}
		// When Not is false and test fails, it means the string didn't match when we expected it to
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected %q at position %d, got %q", ErrSubstringMismatch, o.Str, pos, substring)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *TestStringOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:         string(internal.OpTestStringType),
		Path:       formatPath(o.Path()),
		Str:        o.Str,
		Pos:        o.Pos,
		Not:        o.NotFlag,
		IgnoreCase: o.IgnoreCase,
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *TestStringOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestStringCode, o.Path(), o.Str}, nil
}

// Not returns whether this operation is a negation predicate.
func (o *TestStringOperation) Not() bool {
	return o.NotFlag
}

// Validate validates the test string operation.
func (o *TestStringOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// NewTestStringFull creates a new test string operation with all parameters.
func NewTestStringFull(path []string, str string, pos float64, not bool) *TestStringOperation {
	return &TestStringOperation{
		BaseOp:     NewBaseOp(path),
		Str:        str,
		Pos:        int(pos),
		NotFlag:    not,
		IgnoreCase: false,
	}
}

