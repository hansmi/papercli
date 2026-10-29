package document

import (
	"context"
	"errors"

	"github.com/alecthomas/kingpin/v2"
	"github.com/hansmi/papercli/internal/cli"
	"github.com/hansmi/papercli/internal/printer"
	plclient "github.com/hansmi/paperhooks/pkg/client"
	"github.com/hansmi/paperhooks/pkg/kpflag"
)

type listClient interface {
	ListAllDocuments(context.Context, plclient.ListDocumentsOptions, func(context.Context, plclient.Document) error) error
}

type listHandler struct {
	client listClient

	popts printer.Options
	opts  plclient.ListDocumentsOptions
}

func newListHandler() *listHandler {
	return &listHandler{}
}

func (h *listHandler) Setup(cmd *kingpin.CmdClause) {
	cmd.Help("List documents.")

	printer.RegisterFlags(cmd, &h.popts)
	kpflag.RegisterListDocumentsOptions(cmd, &h.opts)
}

func (h *listHandler) Run(ctx context.Context, cctx cli.Context) error {
	var err error

	client := h.client

	if client == nil {
		if client, err = cctx.Client(); err != nil {
			return err
		}
	}

	p, err := printer.New(cctx.Output(), h.popts)
	if err != nil {
		return err
	}

	defer func() {
		err = errors.Join(err, p.Flush())
	}()

	return client.ListAllDocuments(ctx, h.opts, func(ctx context.Context, doc plclient.Document) error {
		return p.Encode(doc)
	})
}
