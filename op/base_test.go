package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBaseOpClonesPaths(t *testing.T) {
	t.Parallel()

	path := []string{"profile", "name"}
	from := []string{"source", "name"}
	base := NewBaseOpWithFrom(path, from)

	path[0] = "mutated"
	from[0] = "mutated"

	wantPath := []string{"profile", "name"}
	wantFrom := []string{"source", "name"}
	if diff := cmp.Diff(wantPath, base.Path()); diff != "" {
		t.Errorf("Path() after source mutation mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantFrom, base.From()); diff != "" {
		t.Errorf("From() after source mutation mismatch (-want +got):\n%s", diff)
	}

	gotPath := base.Path()
	gotFrom := base.From()
	gotPath[0] = "changed"
	gotFrom[0] = "changed"

	if diff := cmp.Diff(wantPath, base.Path()); diff != "" {
		t.Errorf("Path() after returned slice mutation mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantFrom, base.From()); diff != "" {
		t.Errorf("From() after returned slice mutation mismatch (-want +got):\n%s", diff)
	}
}
