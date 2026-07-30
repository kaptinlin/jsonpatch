package op_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/op"
)

func TestCompositePredicatesUsePublicTypedConstructors(t *testing.T) {
	t.Parallel()

	doc := map[string]any{"name": "Ada", "role": "admin"}
	all := op.NewAnd(nil,
		op.NewDefined([]string{"name"}),
		op.NewContains([]string{"role"}, "admin"),
	)
	anyMatch := op.NewOr(nil,
		op.NewUndefined([]string{"deleted"}),
		op.NewContains([]string{"role"}, "owner"),
	)

	matched, err := all.Test(doc)
	require.NoError(t, err)
	assert.True(t, matched)

	matched, err = anyMatch.Test(doc)
	require.NoError(t, err)
	assert.True(t, matched)
}

func TestNotUsesExplicitParentPathAndOperand(t *testing.T) {
	t.Parallel()

	negated := op.NewNot(
		[]string{"profile"},
		op.NewDefined([]string{"profile", "deleted"}),
	)

	matched, err := negated.Test(map[string]any{
		"profile": map[string]any{"name": "Ada"},
	})
	require.NoError(t, err)
	assert.True(t, matched)
}
