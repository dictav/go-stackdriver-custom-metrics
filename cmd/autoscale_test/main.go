package main

import (
	"context"
	"flag"
	"os"
	"time"

	sdcustom "github.com/dictav/go-stackdriver-custom-metrics"
)

var (
	project = flag.String("project", "", "GCP Project ID")
	zone    = flag.String("zone", "asia-northeast1-a", "GCP Zone")
	metric  = flag.String("metric", "custom.googleapis.com/autoscaling/count", "Custom Metric Name")
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

	m, err := sdcustom.NewMetricReporter(ctx, *project, *zone, *metric, instance)
	if err != nil {
		panic(err)
	}
	m.SetInterval(10 * time.Second)
	m.Start()
	defer m.Stop()

	t1 := time.After(5 * time.Minute)
	t2 := time.After(10 * time.Minute)
	dl := time.After(20 * time.Minute)

LOOP:
	for {
		select {
		case <-t1:
			m.Add(8)
		case <-t2:
			m.Add(-8)
		case <-dl:
			println("bye")
			break LOOP
		}
	}
}
