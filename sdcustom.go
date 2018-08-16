package sdcustom

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/monitoring/v3"
)

// Create new custom metric
func Create(project string, md *monitoring.MetricDescriptor) error {
	ctx := context.Background()
	s, err := createService(ctx)
	if err != nil {
		return err
	}

	pid := projectResource(project)
	_, err = s.Projects.MetricDescriptors.Create(pid, md).Do()
	return err
}

// List custom metrics of project
// fileter syntax: see https://cloud.google.com/monitoring/api/v3/filters#filter_syntax
func List(project, group string) (*monitoring.ListMetricDescriptorsResponse, error) {
	ctx := context.Background()
	s, err := createService(ctx)
	if err != nil {
		return nil, err
	}

	mt := "custom.googleapis.com/" + group
	filter := fmt.Sprintf(`metric.type = starts_with("%s")`, mt)

	pid := projectResource(project)
	return s.Projects.MetricDescriptors.List(pid).Filter(filter).Do()
}

// Get the TimeSeries for the value specified by metric type
func Get(project, metric string) (*monitoring.ListTimeSeriesResponse, error) {
	ctx := context.Background()
	s, err := createService(ctx)
	if err != nil {
		return nil, err
	}

	pid := projectResource(project)
	cond := fmt.Sprintf("metric.type=\"custom.googleapis.com/%s\"", metric)
	st := time.Now().UTC().Add(time.Minute * -5).Format(time.RFC3339Nano)
	end := time.Now().UTC().Format(time.RFC3339Nano)

	return s.Projects.TimeSeries.List(pid).
		Filter(cond).
		IntervalStartTime(st).
		IntervalEndTime(end).
		Do()
}

// Delete custom metric
func Delete(metric string) error {
	ctx := context.Background()
	s, err := createService(ctx)
	if err != nil {
		return err
	}

	_, err = s.Projects.MetricDescriptors.Delete(metric).Do()
	return err
}

func projectResource(project string) string {
	return "projects/" + project
}

func createService(ctx context.Context) (*monitoring.Service, error) {
	hc, err := google.DefaultClient(ctx, monitoring.MonitoringScope)
	if err != nil {
		return nil, err
	}
	s, err := monitoring.New(hc)
	if err != nil {
		return nil, err
	}
	return s, nil
}
