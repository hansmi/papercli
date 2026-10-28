package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/hansmi/paperhooks/pkg/kpflag"
)

type FormatName string

const (
	FormatNdjson FormatName = "ndjson"
)

var allFormats = []FormatName{
	FormatNdjson,
}

type Options struct {
	Format FormatName
}

func RegisterFlags(g kpflag.FlagGroup, opts *Options) {
	var available []string

	for _, i := range allFormats {
		available = append(available, string(i))
	}

	g.Flag("output_format", fmt.Sprintf("Output format (supported: %v)", allFormats)).
		Default(string(FormatNdjson)).
		PlaceHolder("FMT").
		NoEnvar().
		EnumVar((*string)(&opts.Format), available...)
}

// Printer is a generic interface for encoding structured data.
type Printer interface {
	Encode(any) error
	Flush() error
}

func New(w io.Writer, opts Options) (Printer, error) {
	switch opts.Format {
	case FormatNdjson:
		return newNDJSON(w), nil
	}

	return nil, fmt.Errorf("%w: unknown format %q", os.ErrInvalid, opts.Format)
}
