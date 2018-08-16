package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/google/subcommands"
	"os"

	"google.golang.org/api/monitoring/v3"
)

type templateCmd struct {
}

func (*templateCmd) Name() string     { return "template" }
func (*templateCmd) Synopsis() string { return "template a custom metric" }
func (*templateCmd) Usage() string {
	return `template:
  Print custom metric descriptor tempalte.
`
}

func (p *templateCmd) SetFlags(f *flag.FlagSet) {
}

func (p *templateCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	t := &monitoring.MetricDescriptor{
		Description: "A detailed description of the metric, which can be used in documentation",
		DisplayName: "A concise name for the metric, which can be displayed in user interfaces",
		/*
			Labels: []*monitoring.LabelDescriptor{
				{
					Description: "A human-readable description for the label",
					Key:         "The label key",
					ValueType:   "STRING/BOOL/INT64",
				},
			},
		*/
		MetricKind: "GUAGE/DELTA/CUMULATIVE",
		Name:       "The resource name of the metric descriptor",
		Unit:       "items",
		ValueType:  "INT64/DOUBLE/STRING/MONEY",
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&t); err != nil {
		println(err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
