package correspondent

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hansmi/papercli/internal/clitest"
	plclient "github.com/hansmi/paperhooks/pkg/client"
)

type fakeListClient struct {
	items []plclient.Correspondent
	err   error
}

func (c *fakeListClient) ListAllCorrespondents(ctx context.Context, opts plclient.ListCorrespondentsOptions, handle func(context.Context, plclient.Correspondent) error) error {
	if c.err != nil {
		return c.err
	}

	for _, i := range c.items {
		if err := handle(ctx, i); err != nil {
			return err
		}
	}

	return nil
}

func newListHandlerForTest(patch func(*listHandler)) *listHandler {
	h := newListHandler()
	h.client = &fakeListClient{}
	h.popts.Format = "ndjson"

	if patch != nil {
		patch(h)
	}

	return h
}

func TestListHandler(t *testing.T) {
	errTest := errors.New("test error")

	for _, tc := range []struct {
		name    string
		h       *listHandler
		want    *regexp.Regexp
		wantErr error
	}{
		{
			name: "empty",
			h:    newListHandlerForTest(nil),
			want: regexp.MustCompile(`^$`),
		},
		{
			name: "error",
			h: newListHandlerForTest(func(h *listHandler) {
				h.client = &fakeListClient{
					err: errTest,
				}
			}),
			wantErr: errTest,
		},
		{
			name: "items",
			h: newListHandlerForTest(func(h *listHandler) {
				var items []plclient.Correspondent

				for i := range int64(3) {
					items = append(items, plclient.Correspondent{
						ID: i,
					})
				}

				h.client = &fakeListClient{
					items: items,
				}
			}),
			want: regexp.MustCompile(`(?m)^\{.*\blast_correspondence\b.*\}$`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			t.Cleanup(cancel)

			var buf strings.Builder

			cctx := clitest.NewContext(t)
			cctx.SetOutput(&buf)

			err := tc.h.Run(ctx, cctx)

			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Run() error diff (-want +got):\n%s", diff)
			}

			if err == nil && !tc.want.MatchString(buf.String()) {
				t.Errorf("Output doesn't match %q: %q", tc.want.String(), buf.String())
			}
		})
	}
}
