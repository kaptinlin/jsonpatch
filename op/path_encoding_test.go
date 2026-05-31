package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestToJSONEscapesPointerSegments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		op   interface {
			ToJSON() (internal.Operation, error)
		}
		want internal.Operation
	}{
		{
			name: "path segment escapes slash and tilde",
			op:   NewAdd([]string{"a/b", "c~d"}, "value"),
			want: internal.Operation{Op: "add", Path: "/a~1b/c~0d", Value: "value"},
		},
		{
			name: "from segment escapes slash and tilde",
			op:   NewMove([]string{"target/key"}, []string{"source~key"}),
			want: internal.Operation{Op: "move", Path: "/target~1key", From: "/source~0key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.op.ToJSON()
			if err != nil {
				t.Fatalf("ToJSON() error = %v, want nil", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ToJSON() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
