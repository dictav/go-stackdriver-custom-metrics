package main

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	sdcustom "github.com/dictav/go-stackdriver-custom-metrics"
)

type deleteCmd struct {
}

func (*deleteCmd) Name() string     { return "delete" }
func (*deleteCmd) Synopsis() string { return "delete a custom metric" }
func (*deleteCmd) Usage() string {
	return `delete <metric> [...]:
	Create a new custom metric
`
}

func (cmd *deleteCmd) SetFlags(f *flag.FlagSet) {
}

func (cmd *deleteCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	metric := ""
	args := f.Args()
	if len(args) == 0 {
		return subcommands.ExitUsageError
	}

	err := sdcustom.Delete(metric)
	if err != nil {
		println(err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess

}
