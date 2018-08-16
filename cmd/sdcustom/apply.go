package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path"

	"github.com/google/subcommands"
	"google.golang.org/api/monitoring/v3"

	sdcustom "github.com/dictav/go-stackdriver-custom-metrics"
)

const customMetric = "custom.googleapis.com"

type applyCmd struct {
	project string
	zone    string
	metric  string
	file    string
}

func (*applyCmd) Name() string     { return "apply" }
func (*applyCmd) Synopsis() string { return "Apply a custom metric" }
func (*applyCmd) Usage() string {
	return `apply -project <project> -metric <metric> -file <definition file> [-zone <zone>]:
	Create a new custom metric
`
}

func (cmd *applyCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.project, "project", "", "GCP Project ID")
	f.StringVar(&cmd.zone, "zone", "asia-northeast1-a", "GCP Zone")
	f.StringVar(&cmd.metric, "metric", "", "Custom Metric Definition file")
	f.StringVar(&cmd.file, "file", "", "custom metric definication file")
}

func (cmd *applyCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(cmd.project) == 0 || len(cmd.zone) == 0 || len(cmd.metric) == 0 || len(cmd.file) == 0 {
		return subcommands.ExitUsageError
	}

	fd, err := os.Open(cmd.file)
	if err != nil {
		println("here1", err.Error())
		return subcommands.ExitFailure
	}
	defer func() {
		if err := fd.Close(); err != nil {
			println(err.Error())
		}
	}()

	md := monitoring.MetricDescriptor{}
	if err := json.NewDecoder(fd).Decode(&md); err != nil {
		println("here2", err.Error())
		return subcommands.ExitFailure
	}
	md.Type = path.Join(customMetric, cmd.metric)

	err = sdcustom.Create(cmd.project, &md)
	if err != nil {
		println("here2", err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess

}
