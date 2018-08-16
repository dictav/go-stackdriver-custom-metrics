package main

import (
	"context"
	"encoding/json"
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
	metric  string
	debug   bool

	enc *json.Encoder
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
	f.BoolVar(&cmd.debug, "debug", false, "Enable Debug Mode")
}

func (cmd *watchCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	timer := time.NewTicker(5 * time.Second)
	if cmd.debug {
		cmd.enc = json.NewEncoder(os.Stderr)
		cmd.enc.SetIndent("", "  ")
	}

	println("start watching...")
Loop:
	for {
		select {
		case <-timer.C:
			v, err := sdcustom.Get(cmd.project, cmd.metric)
			if err != nil {
				println(err.Error())
				continue
			}

			cmd.logMetric(v)
		case <-c:
			break Loop
		}
	}

	return subcommands.ExitSuccess
}

func (cmd *watchCmd) logMetric(res *monitoring.ListTimeSeriesResponse) {
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
	if cmd.enc != nil {
		cmd.enc.Encode(&res)
	}
}
