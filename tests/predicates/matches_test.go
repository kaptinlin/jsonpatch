package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	operrors "github.com/kaptinlin/jsonpatch/op"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchesOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "",
				Value: "\\d+",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch("123", patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			assert.Equal(t, "123", result.Doc)
		})

		t.Run("fails when does not match the string", func(t *testing.T) {
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "",
				Value: "\\d+",
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch("asdf", patch, jsonpatch.WithMutate(true))
			require.Error(t, err)
			assert.ErrorIs(t, err, operrors.ErrStringMismatch)
		})

		t.Run("succeeds with case insensitive matching", func(t *testing.T) {
			op := jsonpatch.Operation{
				Op:         "matches",
				Path:       "",
				Value:      "HELLO",
				IgnoreCase: true,
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch("hello world", patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			assert.Equal(t, "hello world", result.Doc)
		})

		t.Run("fails with case sensitive matching", func(t *testing.T) {
			op := jsonpatch.Operation{
				Op:         "matches",
				Path:       "",
				Value:      "HELLO",
				IgnoreCase: false,
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch("hello world", patch, jsonpatch.WithMutate(true))
			require.Error(t, err)
		})
	})

	t.Run("nested path", func(t *testing.T) {
		t.Run("matches email pattern", func(t *testing.T) {
			doc := map[string]interface{}{
				"email": "user@example.com",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/email",
				Value: "[a-z]+@[a-z]+\\.[a-z]+",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			resultDoc := result.Doc
			assert.Equal(t, "user@example.com", resultDoc["email"])
		})

		t.Run("fails with invalid email pattern", func(t *testing.T) {
			doc := map[string]interface{}{
				"email": "invalid-email",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/email",
				Value: "^[a-z]+@[a-z]+\\.[a-z]+$",
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			require.Error(t, err)
		})

		t.Run("matches phone number pattern", func(t *testing.T) {
			doc := map[string]interface{}{
				"phone": "123-456-7890",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/phone",
				Value: "\\d{3}-\\d{3}-\\d{4}",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			resultDoc := result.Doc
			assert.Equal(t, "123-456-7890", resultDoc["phone"])
		})
	})

	t.Run("array element", func(t *testing.T) {
		t.Run("matches string in array", func(t *testing.T) {
			doc := map[string]interface{}{
				"items": []interface{}{"apple", "banana", "cherry"},
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/items/1",
				Value: "^b.*a$",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			resultDoc := result.Doc
			items := resultDoc["items"].([]interface{})
			assert.Equal(t, "banana", items[1])
		})
	})

	t.Run("complex patterns", func(t *testing.T) {
		t.Run("matches URL pattern", func(t *testing.T) {
			doc := map[string]interface{}{
				"website": "https://example.com",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/website",
				Value: "^https?://",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			resultDoc := result.Doc
			assert.Equal(t, "https://example.com", resultDoc["website"])
		})

		t.Run("matches UUID pattern", func(t *testing.T) {
			doc := map[string]interface{}{
				"id": "123e4567-e89b-12d3-a456-426614174000",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/id",
				Value: "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			resultDoc := result.Doc
			assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", resultDoc["id"])
		})
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Run("empty pattern matches empty string", func(t *testing.T) {
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "",
				Value: "",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch("", patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			assert.Equal(t, "", result.Doc)
		})

		t.Run("dot matches any character", func(t *testing.T) {
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "",
				Value: ".",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch("x", patch, jsonpatch.WithMutate(true))
			require.NoError(t, err)
			assert.Equal(t, "x", result.Doc)
		})

		t.Run("fails on non-string value", func(t *testing.T) {
			doc := map[string]interface{}{
				"number": 123,
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/number",
				Value: "\\d+",
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			require.Error(t, err)
			assert.ErrorIs(t, err, operrors.ErrNotString)
		})

		t.Run("fails on missing path", func(t *testing.T) {
			doc := map[string]interface{}{
				"field": "value",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/missing",
				Value: ".*",
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			require.Error(t, err)
			assert.ErrorIs(t, err, operrors.ErrPathNotFound)
		})
	})
}
