package main

import (
	"context"
	"flag"
	"os"
	"time"

	sdcustom "github.com/dictav/go-stackdriver-custom-metrics"
)

var (
	projectID = flag.String("project", "", "GCP Project ID")
	zone      = flag.String("zone", "asia-northeast1-a", "GCP Zone")
	metric    = flag.String("metric", "custom.googleapis.com/autoscaling/count", "Custom Metric Name")

	interval = flag.Duration("interval", 5*time.Second, "Count up interval")
	deadline = flag.Duration("deadline", 30*time.Minute, "Deadline")
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		println("Usage: autoscale_test <INSTANCE>", len(args))
		os.Exit(1)
	}
	instance := args[0]

	m, err := sdcustom.NewMetricReporter(ctx, *projectID, *zone, *metric, instance)
	if err != nil {
		panic(err)
	}
	m.SetInterval(10 * time.Second)
	m.Start()
	defer m.Stop()

	t := time.NewTicker(*interval)
	dl := time.After(*deadline)

LOOP:
	for {
		select {
		case <-t.C:
			m.Add(1)
		case <-dl:
			println("bye")
			break LOOP
		}
	}
}
