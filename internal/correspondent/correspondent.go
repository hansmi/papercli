package correspondent

import (
	"github.com/hansmi/papercli/internal/cli"
)

func RegisterCommands(reg *cli.CommandRegistry, g cli.CommandGroup) {
	base := g.Command("correspondent", "")

	reg.Add(base.Command("list", ""), newListHandler())
}
