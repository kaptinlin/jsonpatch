package op

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestRFCMutatorTargetResolution(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		newOp   func([]string) internal.Op
		path    []string
		wantErr error
		check   func(*testing.T, any)
	}{
		{name: "add root", newOp: newTargetAdd, path: nil, check: checkEqual("new")},
		{name: "add empty object key", newOp: newTargetAdd, path: []string{""}, check: checkMapValue("", "new")},
		{name: "add escaped slash token", newOp: newTargetAdd, path: []string{"a/b"}, check: checkMapValue("a/b", "new")},
		{name: "add escaped tilde token", newOp: newTargetAdd, path: []string{"m~n"}, check: checkMapValue("m~n", "new")},
		{name: "add missing middle", newOp: newTargetAdd, path: []string{"missing", "value"}, wantErr: ErrPathNotFound},
		{name: "add missing final", newOp: newTargetAdd, path: []string{"object", "new"}, check: checkNestedMapValue("object", "new", "new")},
		{name: "add array append marker", newOp: newTargetAdd, path: []string{"items", "-"}, check: checkItems("zero", "one", "new")},
		{name: "add array zero", newOp: newTargetAdd, path: []string{"items", "0"}, check: checkItems("new", "zero", "one")},
		{name: "add array leading zero", newOp: newTargetAdd, path: []string{"items", "01"}, wantErr: ErrPathNotFound},
		{name: "add array plus sign", newOp: newTargetAdd, path: []string{"items", "+1"}, wantErr: ErrPathNotFound},
		{name: "add array negative", newOp: newTargetAdd, path: []string{"items", "-1"}, wantErr: ErrPathNotFound},
		{name: "add array end index", newOp: newTargetAdd, path: []string{"items", "2"}, check: checkItems("zero", "one", "new")},
		{name: "add array out of range", newOp: newTargetAdd, path: []string{"items", "3"}, wantErr: ErrIndexOutOfRange},

		{name: "remove root", newOp: newTargetRemove, path: nil, check: checkEqual(nil)},
		{name: "remove empty object key", newOp: newTargetRemove, path: []string{""}, check: checkMapMissing("")},
		{name: "remove escaped slash token", newOp: newTargetRemove, path: []string{"a/b"}, check: checkMapMissing("a/b")},
		{name: "remove missing middle", newOp: newTargetRemove, path: []string{"missing", "value"}, wantErr: ErrPathNotFound},
		{name: "remove missing final", newOp: newTargetRemove, path: []string{"object", "missing"}, wantErr: ErrPathNotFound},
		{name: "remove array append marker", newOp: newTargetRemove, path: []string{"items", "-"}, wantErr: ErrPathNotFound},
		{name: "remove array zero", newOp: newTargetRemove, path: []string{"items", "0"}, check: checkItems("one")},
		{name: "remove array leading zero", newOp: newTargetRemove, path: []string{"items", "01"}, wantErr: ErrPathNotFound},
		{name: "remove array plus sign", newOp: newTargetRemove, path: []string{"items", "+1"}, wantErr: ErrPathNotFound},
		{name: "remove array negative", newOp: newTargetRemove, path: []string{"items", "-1"}, wantErr: ErrPathNotFound},
		{name: "remove array end", newOp: newTargetRemove, path: []string{"items", "2"}, wantErr: ErrIndexOutOfRange},
		{name: "remove array out of range", newOp: newTargetRemove, path: []string{"items", "3"}, wantErr: ErrIndexOutOfRange},

		{name: "replace root", newOp: newTargetReplace, path: nil, check: checkEqual("new")},
		{name: "replace empty object key", newOp: newTargetReplace, path: []string{""}, check: checkMapValue("", "new")},
		{name: "replace escaped tilde token", newOp: newTargetReplace, path: []string{"m~n"}, check: checkMapValue("m~n", "new")},
		{name: "replace missing middle", newOp: newTargetReplace, path: []string{"missing", "value"}, wantErr: ErrPathNotFound},
		{name: "replace missing final", newOp: newTargetReplace, path: []string{"object", "missing"}, wantErr: ErrPathNotFound},
		{name: "replace array append marker", newOp: newTargetReplace, path: []string{"items", "-"}, wantErr: ErrPathNotFound},
		{name: "replace array zero", newOp: newTargetReplace, path: []string{"items", "0"}, check: checkItems("new", "one")},
		{name: "replace array leading zero", newOp: newTargetReplace, path: []string{"items", "01"}, wantErr: ErrPathNotFound},
		{name: "replace array plus sign", newOp: newTargetReplace, path: []string{"items", "+1"}, wantErr: ErrPathNotFound},
		{name: "replace array negative", newOp: newTargetReplace, path: []string{"items", "-1"}, wantErr: ErrPathNotFound},
		{name: "replace array end", newOp: newTargetReplace, path: []string{"items", "2"}, wantErr: ErrIndexOutOfRange},
		{name: "replace array out of range", newOp: newTargetReplace, path: []string{"items", "3"}, wantErr: ErrIndexOutOfRange},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := tc.newOp(tc.path).Apply(targetResolutionDocument())
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			tc.check(t, result.Doc)
		})
	}
}

func TestMoveCopyTargetResolution(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		newOp   func(path, from []string) internal.Op
		path    []string
		from    []string
		wantErr error
		check   func(*testing.T, any)
	}{
		{name: "copy root to root", newOp: newTargetCopy, check: checkCopiedRoot},
		{name: "copy empty key source", newOp: newTargetCopy, path: []string{"object", "copy"}, from: []string{""}, check: checkNestedMapValue("object", "copy", "empty")},
		{name: "copy array zero source to append target", newOp: newTargetCopy, path: []string{"items", "-"}, from: []string{"items", "0"}, check: checkItems("zero", "one", "zero")},
		{name: "copy missing source", newOp: newTargetCopy, path: []string{"object", "copy"}, from: []string{"missing"}, wantErr: ErrPathNotFound},
		{name: "copy append marker source", newOp: newTargetCopy, path: []string{"object", "copy"}, from: []string{"items", "-"}, wantErr: ErrPathNotFound},
		{name: "copy leading zero source", newOp: newTargetCopy, path: []string{"object", "copy"}, from: []string{"items", "01"}, wantErr: ErrPathNotFound},
		{name: "copy plus sign source", newOp: newTargetCopy, path: []string{"object", "copy"}, from: []string{"items", "+1"}, wantErr: ErrPathNotFound},
		{name: "copy negative source", newOp: newTargetCopy, path: []string{"object", "copy"}, from: []string{"items", "-1"}, wantErr: ErrPathNotFound},
		{name: "copy end source", newOp: newTargetCopy, path: []string{"object", "copy"}, from: []string{"items", "2"}, wantErr: ErrPathNotFound},
		{name: "copy out of range source", newOp: newTargetCopy, path: []string{"object", "copy"}, from: []string{"items", "3"}, wantErr: ErrPathNotFound},
		{name: "copy leading zero target", newOp: newTargetCopy, path: []string{"items", "01"}, from: []string{"items", "0"}, wantErr: ErrPathNotFound},
		{name: "copy missing middle target", newOp: newTargetCopy, path: []string{"missing", "copy"}, from: []string{"items", "0"}, wantErr: ErrPathNotFound},

		{name: "move root same path", newOp: newTargetMove, check: checkCopiedRoot},
		{name: "move empty key source", newOp: newTargetMove, path: []string{"object", "moved"}, from: []string{""}, check: checkMovedEmptyKey},
		{name: "move array zero source to append target", newOp: newTargetMove, path: []string{"items", "-"}, from: []string{"items", "0"}, check: checkItems("one", "zero")},
		{name: "move missing same path", newOp: newTargetMove, path: []string{"missing"}, from: []string{"missing"}, wantErr: ErrPathNotFound},
		{name: "move append marker same path", newOp: newTargetMove, path: []string{"items", "-"}, from: []string{"items", "-"}, wantErr: ErrPathNotFound},
		{name: "move leading zero same path", newOp: newTargetMove, path: []string{"items", "01"}, from: []string{"items", "01"}, wantErr: ErrPathNotFound},
		{name: "move plus sign same path", newOp: newTargetMove, path: []string{"items", "+1"}, from: []string{"items", "+1"}, wantErr: ErrPathNotFound},
		{name: "move negative same path", newOp: newTargetMove, path: []string{"items", "-1"}, from: []string{"items", "-1"}, wantErr: ErrPathNotFound},
		{name: "move end same path", newOp: newTargetMove, path: []string{"items", "2"}, from: []string{"items", "2"}, wantErr: ErrPathNotFound},
		{name: "move out of range same path", newOp: newTargetMove, path: []string{"items", "3"}, from: []string{"items", "3"}, wantErr: ErrPathNotFound},
		{name: "move leading zero target", newOp: newTargetMove, path: []string{"items", "01"}, from: []string{"items", "0"}, wantErr: ErrPathNotFound},
		{name: "move missing middle target", newOp: newTargetMove, path: []string{"missing", "moved"}, from: []string{"items", "0"}, wantErr: ErrPathNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := tc.newOp(tc.path, tc.from).Apply(targetResolutionDocument())
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			tc.check(t, result.Doc)
		})
	}
}

func targetResolutionDocument() map[string]any {
	return map[string]any{
		"":       "empty",
		"a/b":    "slash",
		"m~n":    "tilde",
		"items":  []any{"zero", "one"},
		"object": map[string]any{"present": "old"},
	}
}

func newTargetAdd(path []string) internal.Op {
	return NewAdd(path, "new")
}

func newTargetRemove(path []string) internal.Op {
	return NewRemove(path)
}

func newTargetReplace(path []string) internal.Op {
	return NewReplace(path, "new")
}

func newTargetCopy(path, from []string) internal.Op {
	return NewCopy(path, from)
}

func newTargetMove(path, from []string) internal.Op {
	return NewMove(path, from)
}

func checkEqual(want any) func(*testing.T, any) {
	return func(t *testing.T, got any) {
		t.Helper()
		assert.Equal(t, want, got)
	}
}

func checkMapValue(key string, want any) func(*testing.T, any) {
	return func(t *testing.T, got any) {
		t.Helper()
		doc := got.(map[string]any)
		assert.Equal(t, want, doc[key])
	}
}

func checkMapMissing(key string) func(*testing.T, any) {
	return func(t *testing.T, got any) {
		t.Helper()
		doc := got.(map[string]any)
		assert.NotContains(t, doc, key)
	}
}

func checkNestedMapValue(parent, key string, want any) func(*testing.T, any) {
	return func(t *testing.T, got any) {
		t.Helper()
		doc := got.(map[string]any)
		nested := doc[parent].(map[string]any)
		assert.Equal(t, want, nested[key])
	}
}

func checkItems(want ...any) func(*testing.T, any) {
	return func(t *testing.T, got any) {
		t.Helper()
		doc := got.(map[string]any)
		assert.Equal(t, want, doc["items"])
	}
}

func checkCopiedRoot(t *testing.T, got any) {
	t.Helper()
	assert.Equal(t, targetResolutionDocument(), got)
}

func checkMovedEmptyKey(t *testing.T, got any) {
	t.Helper()
	doc := got.(map[string]any)
	assert.NotContains(t, doc, "")
	checkNestedMapValue("object", "moved", "empty")(t, got)
}
