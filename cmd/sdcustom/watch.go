package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/google/subcommands"
	"google.golang.org/api/monitoring/v3"

	sdcustom "github.com/dictav/go-stackdriver-custom-metrics"
)

type watchCmd struct {
	project string
	zone    string
	group   string
	metric  string
}

func (*watchCmd) Name() string     { return "watch" }
func (*watchCmd) Synopsis() string { return "watch a custom metric" }
func (*watchCmd) Usage() string {
	return `watch:
  Watch custom metric 
`
}

func (cmd *watchCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.project, "project", "", "GCP Project ID")
	f.StringVar(&cmd.metric, "metric", "", "Custom Metric name")
}

func (cmd *watchCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	timer := time.NewTicker(5 * time.Second)

Loop:
	for {
		select {
		case <-timer.C:
			v, err := sdcustom.Get(cmd.project, cmd.metric)
			if err != nil {
				println(err.Error())
				continue
			}

			logMetric(v)
		case <-c:
			break Loop
		}
	}

	return subcommands.ExitSuccess
}

func logMetric(res *monitoring.ListTimeSeriesResponse) {
	if len(res.TimeSeries) == 0 {
		return
	}

	ts := res.TimeSeries[0]
	if len(ts.Points) == 0 {
		return
	}

	var v interface{}
	val := ts.Points[0].Value
	instance := ts.Resource.Labels["instance_id"]
	switch {
	case val.BoolValue != nil:
		v = *val.BoolValue
	case val.DoubleValue != nil:
		v = *val.DoubleValue
	case val.Int64Value != nil:
		v = *val.Int64Value
	case val.StringValue != nil:
		v = *val.StringValue
	}

	log.Printf("%s: %v, labels=%s", instance, v, ts.Resource.Labels)
}
