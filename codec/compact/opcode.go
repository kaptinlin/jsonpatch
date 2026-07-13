package compact

import (
	"fmt"
	"math"

	"github.com/kaptinlin/jsonpatch/internal"
)

// resolveOpType determines the operation type from the opcode.
func resolveOpType(opcode any) (internal.OpType, error) {
	if s, ok := opcode.(string); ok {
		if spec, exists := internal.LookupOperation(internal.OpType(s)); exists {
			return spec.Type, nil
		}
		return "", fmt.Errorf("%w: %s", ErrUnknownStringCode, s)
	}

	var code int
	switch v := opcode.(type) {
	case int:
		code = v
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) || v != math.Trunc(v) {
			return "", fmt.Errorf("%w: %v", ErrInvalidCodeType, v)
		}
		code = int(v)
		if float64(code) != v {
			return "", fmt.Errorf("%w: %v", ErrInvalidCodeType, v)
		}
	case Code:
		code = int(v)
	default:
		return "", fmt.Errorf("%w: %T", ErrInvalidCodeType, opcode)
	}

	if spec, exists := internal.LookupOperationCode(code); exists {
		return spec.Type, nil
	}
	return "", fmt.Errorf("%w: %d", ErrUnknownNumericCode, code)
}
