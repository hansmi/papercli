package printer

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFormatsSorted(t *testing.T) {
	if !slices.IsSorted(allFormats) {
		t.Errorf("List of formats is not sorted: %v", allFormats)
	}
}

func TestRegisterFlags(t *testing.T) {
	app := kingpin.New(t.Name(), "")

	var opts Options

	RegisterFlags(app, &opts)

	if _, err := app.Parse(nil); err != nil {
		t.Errorf("Parsing arguments failed: %v", err)
	}
}

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		name         string
		opts         Options
		wantMinBytes int
		wantErr      error
	}{
		{
			name:    "empty",
			wantErr: os.ErrInvalid,
		},
		{
			name: "ndjson",
			opts: Options{
				Format: FormatNdjson,
			},
			wantMinBytes: 20,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var buf strings.Builder

			got, err := New(&buf, tc.opts)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("New() error diff (-want +got):\n%s", diff)
			}

			if err == nil {
				if err := got.Encode(map[string]int{
					"first":  1,
					"second": 2,
				}); err != nil {
					t.Errorf("Encode() failed: %v", err)
				}

				if err := got.Flush(); err != nil {
					t.Errorf("Flush() failed: %v", err)
				}

				if count := buf.Len(); count < tc.wantMinBytes {
					t.Errorf("%d bytes were written, want at least %d: %q", count, tc.wantMinBytes, buf.String())
				}
			}
		})
	}
}
