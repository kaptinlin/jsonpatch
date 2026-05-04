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
	assertOutputContains(t, output, "Round-trip successful")
	assertOutputContains(t, output, "codec := binary.Codec{}")
}

func TestRoundTripMessageReportsOperationCount(t *testing.T) {
	t.Parallel()

	ops := []internal.Op{op.NewAdd([]string{"name"}, "Grace")}
	if got, want := roundTripMessage(ops, ops), "✅ Round-trip successful (operation count matches)"; got != want {
		t.Fatalf("roundTripMessage() = %q, want %q", got, want)
	}
	if got, want := roundTripMessage(ops, nil), "❌ Round-trip failed (operation count mismatch)"; got != want {
		t.Fatalf("roundTripMessage() = %q, want %q", got, want)
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
