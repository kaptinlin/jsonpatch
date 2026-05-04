package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestMainShowsExpectedOutput(t *testing.T) {
	// This test redirects process-wide stdout.
	output := captureOutput(t, main)
	assertOutputContains(t, output, "All operations perfectly decoded")
	assertOutputContains(t, output, "encoded, err := compact.Encode(operations)")
}

func TestRoundTripMessageReportsDecodedOperationMismatch(t *testing.T) {
	t.Parallel()

	ops := []internal.Op{op.NewAdd([]string{"user", "name"}, "Grace")}
	if got, want := roundTripMessage(ops, ops), "✅ All operations perfectly decoded!"; got != want {
		t.Fatalf("roundTripMessage() = %q, want %q", got, want)
	}

	mismatches := []struct {
		name    string
		decoded []internal.Op
	}{
		{name: "shorter", decoded: nil},
		{name: "operation type", decoded: []internal.Op{op.NewReplace([]string{"user", "name"}, "Grace")}},
		{name: "path length", decoded: []internal.Op{op.NewAdd([]string{"user"}, "Grace")}},
		{name: "path segment", decoded: []internal.Op{op.NewAdd([]string{"user", "email"}, "Grace")}},
	}
	for _, tc := range mismatches {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got, want := roundTripMessage(ops, tc.decoded), "❌ Some operations failed to decode correctly"; got != want {
				t.Fatalf("roundTripMessage() = %q, want %q", got, want)
			}
		})
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	return string(output)
}

func assertOutputContains(t *testing.T, output, want string) {
	t.Helper()

	if !strings.Contains(output, want) {
		t.Fatalf("output does not contain %q:\n%s", want, output)
	}
}
