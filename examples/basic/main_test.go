package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/kaptinlin/jsonpatch/op"
)

func TestMainShowsExpectedOutput(t *testing.T) {
	// This test redirects process-wide stdout.
	output := captureOutput(t, main)
	assertOutputContains(t, output, "Jane")
	assertOutputContains(t, output, "john@example.com")
}

func TestRunReportsPatchErrors(t *testing.T) {
	// This test writes to process-wide stdout.
	err := run([]jsoncodec.Operation{{Op: "replace", Path: "/missing", Value: "value"}})
	if !errors.Is(err, op.ErrPathNotFound) {
		t.Fatalf("run() error = %v, want ErrPathNotFound", err)
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
