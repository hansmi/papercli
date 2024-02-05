package cli

import "github.com/alecthomas/kingpin/v2"

// Implemented by [kingpin.CommandLine] (an instance of [kingpin.Application]).
type CommandGroup interface {
	Command(name, help string) *kingpin.CmdClause
}

var _ CommandGroup = (*kingpin.Application)(nil)
var _ CommandGroup = kingpin.CommandLine
