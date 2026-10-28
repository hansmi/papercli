package printer

import (
	"encoding/json"
	"io"
)

type ndjson struct {
	w   io.Writer
	enc *json.Encoder
}

func newNDJSON(w io.Writer) *ndjson {
	return &ndjson{
		w:   w,
		enc: json.NewEncoder(w),
	}
}

func (e *ndjson) Encode(v any) error {
	return e.enc.Encode(v)
}

func (e *ndjson) Flush() error {
	return nil
}
