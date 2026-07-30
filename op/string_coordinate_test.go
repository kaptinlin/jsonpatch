package op

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPositionalStringOperationsUseCodePointCoordinates(t *testing.T) {
	t.Parallel()

	const text = "A界🙂e\u0301Z"

	t.Run("str_ins", func(t *testing.T) {
		t.Parallel()
		result, err := NewStrIns(nil, 2, "X").Apply(text)
		require.NoError(t, err)
		assert.Equal(t, "A界X🙂e\u0301Z", result.Doc)
	})

	t.Run("str_del length", func(t *testing.T) {
		t.Parallel()
		result, err := NewStrDel(nil, 2, 1).Apply(text)
		require.NoError(t, err)
		assert.Equal(t, "A界e\u0301Z", result.Doc)
	})

	t.Run("str_del string", func(t *testing.T) {
		t.Parallel()
		result, err := NewStrDelWithStr(nil, 3, "e\u0301").Apply(text)
		require.NoError(t, err)
		assert.Equal(t, "A界🙂Z", result.Doc)
	})

	t.Run("split", func(t *testing.T) {
		t.Parallel()
		result, err := NewSplit(nil, 2, nil).Apply(text)
		require.NoError(t, err)
		assert.Equal(t, []any{"A界", "🙂e\u0301Z"}, result.Doc)
	})

	t.Run("test_string supplementary code point", func(t *testing.T) {
		t.Parallel()
		operation := NewTestString(nil, "🙂", 2, false)
		matched, err := operation.Test(text)
		require.NoError(t, err)
		assert.True(t, matched)
		_, err = operation.Apply(text)
		require.NoError(t, err)
	})

	t.Run("test_string combining sequence counts two code points", func(t *testing.T) {
		t.Parallel()
		operation := NewTestString(nil, "e\u0301", 3, false)
		matched, err := operation.Test(text)
		require.NoError(t, err)
		assert.True(t, matched)
		_, err = operation.Apply(text)
		require.NoError(t, err)
	})

	t.Run("test_string_len", func(t *testing.T) {
		t.Parallel()
		atLeastSix := NewTestStringLen(nil, 6)
		matched, err := atLeastSix.Test(text)
		require.NoError(t, err)
		assert.True(t, matched)
		_, err = atLeastSix.Apply(text)
		require.NoError(t, err)

		atLeastSeven := NewTestStringLen(nil, 7)
		matched, err = atLeastSeven.Test(text)
		require.NoError(t, err)
		assert.False(t, matched)
		_, err = atLeastSeven.Apply(text)
		require.ErrorIs(t, err, ErrStringLengthMismatch)
	})
}

func TestPositionalStringOperationBoundaries(t *testing.T) {
	t.Parallel()

	t.Run("ASCII", func(t *testing.T) {
		t.Parallel()
		operation := NewTestString(nil, "b", 1, false)
		matched, err := operation.Test("abc")
		require.NoError(t, err)
		assert.True(t, matched)
		_, err = operation.Apply("abc")
		require.NoError(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		inserted, err := NewStrIns(nil, 0, "x").Apply("")
		require.NoError(t, err)
		assert.Equal(t, "x", inserted.Doc)

		deleted, err := NewStrDel(nil, 0, 1).Apply("")
		require.NoError(t, err)
		assert.Equal(t, "", deleted.Doc)

		split, err := NewSplit(nil, 0, nil).Apply("")
		require.NoError(t, err)
		assert.Equal(t, []any{"", ""}, split.Doc)

		testString := NewTestString(nil, "", 0, false)
		matched, err := testString.Test("")
		require.NoError(t, err)
		assert.True(t, matched)
		_, err = testString.Apply("")
		require.NoError(t, err)

		testLength := NewTestStringLen(nil, 0)
		matched, err = testLength.Test("")
		require.NoError(t, err)
		assert.True(t, matched)
		_, err = testLength.Apply("")
		require.NoError(t, err)
	})

	t.Run("negative", func(t *testing.T) {
		t.Parallel()

		inserted, err := NewStrIns(nil, -1, "X").Apply("A🙂Z")
		require.NoError(t, err)
		assert.Equal(t, "A🙂XZ", inserted.Doc)

		deleted, err := NewStrDel(nil, -1, 1).Apply("A🙂Z")
		require.NoError(t, err)
		assert.Equal(t, "A🙂", deleted.Doc)

		split, err := NewSplit(nil, -1, nil).Apply("A🙂Z")
		require.NoError(t, err)
		assert.Equal(t, []any{"A🙂", "Z"}, split.Doc)

		testString := NewTestString(nil, "A", -1, false)
		matched, err := testString.Test("A🙂Z")
		assert.False(t, matched)
		require.ErrorIs(t, err, ErrPositionOutOfStringRange)
		_, err = testString.Apply("A🙂Z")
		require.ErrorIs(t, err, ErrPositionOutOfStringRange)

		require.ErrorIs(t, NewTestStringLen(nil, -1).Validate(), ErrLengthNegative)
	})

	t.Run("end boundary", func(t *testing.T) {
		t.Parallel()
		const text = "A🙂"

		inserted, err := NewStrIns(nil, 2, "Z").Apply(text)
		require.NoError(t, err)
		assert.Equal(t, "A🙂Z", inserted.Doc)

		deleted, err := NewStrDel(nil, 2, 1).Apply(text)
		require.NoError(t, err)
		assert.Equal(t, text, deleted.Doc)

		split, err := NewSplit(nil, 2, nil).Apply(text)
		require.NoError(t, err)
		assert.Equal(t, []any{text, ""}, split.Doc)

		testString := NewTestString(nil, "", 2, false)
		matched, err := testString.Test(text)
		require.NoError(t, err)
		assert.True(t, matched)
		_, err = testString.Apply(text)
		require.NoError(t, err)
	})

	t.Run("test_string out of range", func(t *testing.T) {
		t.Parallel()
		operation := NewTestString(nil, "", 3, false)
		matched, err := operation.Test("A🙂")
		assert.False(t, matched)
		require.ErrorIs(t, err, ErrPositionOutOfStringRange)
		_, err = operation.Apply("A🙂")
		require.ErrorIs(t, err, ErrPositionOutOfStringRange)
	})
}
