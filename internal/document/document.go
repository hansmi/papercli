package document

import (
	"github.com/hansmi/papercli/internal/cli"
)

func RegisterCommands(reg *cli.CommandRegistry, g cli.CommandGroup) {
	base := g.Command("document", "")

	reg.Add(base.Command("upload", ""), newUploadHandler())
}
