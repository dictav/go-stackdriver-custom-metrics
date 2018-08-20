package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	"github.com/google/subcommands"

	"github.com/dictav/go-stackdriver-custom-metrics"
)

type listCmd struct {
	project string
	group   string
}

func (*listCmd) Name() string     { return "list" }
func (*listCmd) Synopsis() string { return "list a custom metric" }
func (*listCmd) Usage() string {
	return `list -project=<PROJECT> [-group=<GROUP>]
  Print custom metric descriptor tempalte.
`
}

func (cmd *listCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.project, "project", "", "GCP Project")
	f.StringVar(&cmd.group, "group", "", "Filter by metric group")
}

func (cmd *listCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(cmd.project) == 0 {
		println(cmd.Usage())
		return subcommands.ExitUsageError
	}

	res, err := sdcustom.List(cmd.project, cmd.group)
	if err != nil {
		println(err.Error())
		return subcommands.ExitFailure
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&res); err != nil {
		println(err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
