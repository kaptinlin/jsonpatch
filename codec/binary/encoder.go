package binary

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/internal"
)

// encodeOps writes the operation count followed by each compact operation tree.
func encodeOps(w *msgp.Writer, operations []internal.Op) error {
	if err := w.WriteArrayHeader(uint32(len(operations))); err != nil { //nolint:gosec // Operation count is bounded by practical limits.
		return err
	}
	for _, operation := range operations {
		if err := encodeOp(w, operation); err != nil {
			return err
		}
	}
	return nil
}

func encodeOp(w *msgp.Writer, operation internal.Op) error {
	compactOperation, ok := operation.(internal.CompactOp)
	if !ok {
		return fmt.Errorf("operation %T: %w", operation, ErrUnsupportedOp)
	}
	tree, err := compactOperation.ToCompact()
	if err != nil {
		return err
	}
	return w.WriteIntf(tree)
}
