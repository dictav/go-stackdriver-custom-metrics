package main

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"

	sdcustom "github.com/dictav/go-stackdriver-custom-metrics"
)

var (
	project  = flag.String("project", "", "GCP Project ID")
	zone     = flag.String("zone", "asia-northeast1-a", "GCP Zone")
	group    = flag.String("group", "autoscale-test", "GCP Autoscaling Group")
	metric   = flag.String("metric", "autoscaling/count", "Custom Metric Name")
	refValue = flag.Int64("value", 80, "Reference value")
	maxIns   = flag.Int("max-instances", 5, "Max instances")
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

	cs, err := getComputeService(ctx)
	if err != nil {
		println(err.Error())
		os.Exit(2)
	}

	m, err := sdcustom.NewMetricReporter(ctx, *project, *zone, *metric, instance)
	if err != nil {
		println(err.Error())
		os.Exit(2)
	}
	m.Start()
	defer m.Stop()

	i := time.NewTicker(15 * time.Second)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

LOOP:
	for {
		select {
		case <-i.C:
			n := numberOfInstances(cs, *project, *zone, *group)
			if n >= *maxIns {
				m.Set(rand.Int63n(10))
			} else {
				m.Set(*refValue / int64(n))
			}
		case <-sig:
			break LOOP
		}
	}
}

func getComputeService(ctx context.Context) (*compute.Service, error) {
	client, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	return compute.New(client)
}

func numberOfInstances(cs *compute.Service, project, zone, group string) int {
	req := &compute.InstanceGroupsListInstancesRequest{InstanceState: "ALL"}
	igs := compute.NewInstanceGroupsService(cs)
	list, err := igs.ListInstances(project, zone, group, req).Do()
	if err != nil {
		println("error:", err.Error())
		return 0
	}

	return len(list.Items)
}
