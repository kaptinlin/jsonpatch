package compact

import (
	"fmt"

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
		code = int(v)
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
