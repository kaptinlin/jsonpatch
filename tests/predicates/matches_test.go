package ops_test

import (
	"errors"
	"testing"

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
				t.Errorf("ApplyPatch() doc = %v, want %v", got, "123")
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
				t.Errorf("ApplyPatch() error = %v, want %v", err, operrors.ErrStringMismatch)
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
				t.Errorf("ApplyPatch() doc = %v, want %v", got, "hello world")
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
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc["email"]; got != "user@example.com" {
				t.Errorf("result.Doc[email] = %v, want %v", got, "user@example.com")
			}
		})

		t.Run("fails with invalid email pattern", func(t *testing.T) {
			t.Parallel()
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
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})

		t.Run("matches phone number pattern", func(t *testing.T) {
			t.Parallel()
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
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc["phone"]; got != "123-456-7890" {
				t.Errorf("result.Doc[phone] = %v, want %v", got, "123-456-7890")
			}
		})
	})

	t.Run("array element", func(t *testing.T) {
		t.Parallel()
		t.Run("matches string in array", func(t *testing.T) {
			t.Parallel()
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
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			items := result.Doc["items"].([]interface{})
			if got := items[1]; got != "banana" {
				t.Errorf("items[1] = %v, want %v", got, "banana")
			}
		})
	})

	t.Run("complex patterns", func(t *testing.T) {
		t.Parallel()
		t.Run("matches URL pattern", func(t *testing.T) {
			t.Parallel()
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
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc["website"]; got != "https://example.com" {
				t.Errorf("result.Doc[website] = %v, want %v", got, "https://example.com")
			}
		})

		t.Run("matches UUID pattern", func(t *testing.T) {
			t.Parallel()
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
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
			if got := result.Doc["id"]; got != "123e4567-e89b-12d3-a456-426614174000" {
				t.Errorf("result.Doc[id] = %v, want %v", got, "123e4567-e89b-12d3-a456-426614174000")
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
				t.Errorf("ApplyPatch() doc = %q, want %q", got, "")
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
				t.Errorf("ApplyPatch() doc = %v, want %v", got, "x")
			}
		})

		t.Run("fails on non-string value", func(t *testing.T) {
			t.Parallel()
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
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
			if !errors.Is(err, operrors.ErrNotString) {
				t.Errorf("ApplyPatch() error = %v, want %v", err, operrors.ErrNotString)
			}
		})

		t.Run("fails on missing path", func(t *testing.T) {
			t.Parallel()
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
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
			if !errors.Is(err, operrors.ErrPathNotFound) {
				t.Errorf("ApplyPatch() error = %v, want %v", err, operrors.ErrPathNotFound)
			}
		})
	})
}
