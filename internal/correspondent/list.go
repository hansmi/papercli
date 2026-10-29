package correspondent

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
	ListAllCorrespondents(context.Context, plclient.ListCorrespondentsOptions, func(context.Context, plclient.Correspondent) error) error
}

type listHandler struct {
	client listClient

	popts printer.Options
	opts  plclient.ListCorrespondentsOptions
}

func newListHandler() *listHandler {
	return &listHandler{}
}

func (h *listHandler) Setup(cmd *kingpin.CmdClause) {
	cmd.Help("List correspondents.")

	printer.RegisterFlags(cmd, &h.popts)
	kpflag.RegisterListCorrespondentsOptions(cmd, &h.opts)
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

	return client.ListAllCorrespondents(ctx, h.opts, func(ctx context.Context, c plclient.Correspondent) error {
		return p.Encode(c)
	})
}
