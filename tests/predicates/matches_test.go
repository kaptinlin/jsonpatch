package ops_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kaptinlin/jsonpatch"
	operrors "github.com/kaptinlin/jsonpatch/op"
)

func TestMatchesOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			t.Parallel()
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "",
				Value: "\\d+",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch("123", patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc; got != "123" {
				assert.Equal(t, "123", got, "ApplyPatch() doc")
			}
		})

		t.Run("fails when does not match the string", func(t *testing.T) {
			t.Parallel()
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "",
				Value: "\\d+",
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch("asdf", patch, jsonpatch.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
			if !errors.Is(err, operrors.ErrStringMismatch) {
				assert.Equal(t, operrors.ErrStringMismatch, err, "ApplyPatch() error")
			}
		})

		t.Run("succeeds with case insensitive matching", func(t *testing.T) {
			t.Parallel()
			op := jsonpatch.Operation{
				Op:         "matches",
				Path:       "",
				Value:      "HELLO",
				IgnoreCase: true,
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch("hello world", patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc; got != "hello world" {
				assert.Equal(t, "hello world", got, "ApplyPatch() doc")
			}
		})

		t.Run("fails with case sensitive matching", func(t *testing.T) {
			t.Parallel()
			op := jsonpatch.Operation{
				Op:         "matches",
				Path:       "",
				Value:      "HELLO",
				IgnoreCase: false,
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch("hello world", patch, jsonpatch.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})
	})

	t.Run("nested path", func(t *testing.T) {
		t.Parallel()
		t.Run("matches email pattern", func(t *testing.T) {
			t.Parallel()
			doc := map[string]any{
				"email": "user@example.com",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/email",
				Value: "[a-z]+@[a-z]+\\.[a-z]+",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc["email"]; got != "user@example.com" {
				assert.Equal(t, "user@example.com", got, "result.Doc[email]")
			}
		})

		t.Run("fails with invalid email pattern", func(t *testing.T) {
			t.Parallel()
			doc := map[string]any{
				"email": "invalid-email",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/email",
				Value: "^[a-z]+@[a-z]+\\.[a-z]+$",
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})

		t.Run("matches phone number pattern", func(t *testing.T) {
			t.Parallel()
			doc := map[string]any{
				"phone": "123-456-7890",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/phone",
				Value: "\\d{3}-\\d{3}-\\d{4}",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc["phone"]; got != "123-456-7890" {
				assert.Equal(t, "123-456-7890", got, "result.Doc[phone]")
			}
		})
	})

	t.Run("array element", func(t *testing.T) {
		t.Parallel()
		t.Run("matches string in array", func(t *testing.T) {
			t.Parallel()
			doc := map[string]any{
				"items": []any{"apple", "banana", "cherry"},
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/items/1",
				Value: "^b.*a$",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			items := result.Doc["items"].([]any)
			if got := items[1]; got != "banana" {
				assert.Equal(t, "banana", got, "items[1]")
			}
		})
	})

	t.Run("complex patterns", func(t *testing.T) {
		t.Parallel()
		t.Run("matches URL pattern", func(t *testing.T) {
			t.Parallel()
			doc := map[string]any{
				"website": "https://example.com",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/website",
				Value: "^https?://",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc["website"]; got != "https://example.com" {
				assert.Equal(t, "https://example.com", got, "result.Doc[website]")
			}
		})

		t.Run("matches UUID pattern", func(t *testing.T) {
			t.Parallel()
			doc := map[string]any{
				"id": "123e4567-e89b-12d3-a456-426614174000",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/id",
				Value: "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc["id"]; got != "123e4567-e89b-12d3-a456-426614174000" {
				assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", got, "result.Doc[id]")
			}
		})
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Parallel()
		t.Run("empty pattern matches empty string", func(t *testing.T) {
			t.Parallel()
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "",
				Value: "",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch("", patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc; got != "" {
				assert.Equal(t, "", got, "ApplyPatch() doc")
			}
		})

		t.Run("dot matches any character", func(t *testing.T) {
			t.Parallel()
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "",
				Value: ".",
			}
			patch := []jsonpatch.Operation{op}
			result, err := jsonpatch.ApplyPatch("x", patch, jsonpatch.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc; got != "x" {
				assert.Equal(t, "x", got, "ApplyPatch() doc")
			}
		})

		t.Run("fails on non-string value", func(t *testing.T) {
			t.Parallel()
			doc := map[string]any{
				"number": 123,
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/number",
				Value: "\\d+",
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
			if !errors.Is(err, operrors.ErrNotString) {
				assert.Equal(t, operrors.ErrNotString, err, "ApplyPatch() error")
			}
		})

		t.Run("fails on missing path", func(t *testing.T) {
			t.Parallel()
			doc := map[string]any{
				"field": "value",
			}
			op := jsonpatch.Operation{
				Op:    "matches",
				Path:  "/missing",
				Value: ".*",
			}
			patch := []jsonpatch.Operation{op}
			_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
			if !errors.Is(err, operrors.ErrPathNotFound) {
				assert.Equal(t, operrors.ErrPathNotFound, err, "ApplyPatch() error")
			}
		})
	})
}
