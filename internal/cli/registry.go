package cli

import "github.com/alecthomas/kingpin/v2"

type CommandRegistry struct {
	entries map[string]Handler
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		entries: map[string]Handler{},
	}
}

func (r *CommandRegistry) Add(cmd *kingpin.CmdClause, h Handler) {
	h.Setup(cmd)

	r.entries[cmd.FullCommand()] = h
}

func (r *CommandRegistry) Lookup(name string) Handler {
	return r.entries[name]
}
