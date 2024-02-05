package cli

import (
	"context"

	"github.com/alecthomas/kingpin/v2"
)

type Handler interface {
	Setup(*kingpin.CmdClause)
	Run(context.Context, Context) error
}
