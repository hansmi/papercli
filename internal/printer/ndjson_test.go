package printer

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNdjson(t *testing.T) {
	var buf strings.Builder

	p := newNDJSON(&buf)

	for _, value := range []any{
		nil,
		map[string]string{
			"a": "hello",
			"b": "world",
		},
		"test",
	} {
		if err := p.Encode(value); err != nil {
			t.Errorf("Encode(%v) failed: %v", value, err)
		}
	}

	if err := p.Flush(); err != nil {
		t.Errorf("Flush() failed: %v", err)
	}

	want := strings.Join([]string{
		"null\n",
		`{"a":"hello","b":"world"}` + "\n",
		`"test"` + "\n",
	}, "")

	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("Output diff (-want +got):\n%s", diff)
	}
}
