package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

func panicIfError(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&applyCmd{}, "")
	subcommands.Register(&listCmd{}, "")
	subcommands.Register(&templateCmd{}, "")
	subcommands.Register(&watchCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	ret := subcommands.Execute(ctx)
	if ret == subcommands.ExitUsageError {
		println(subcommands.CommandsCommand().Usage())
	}
	os.Exit(int(ret))
}
