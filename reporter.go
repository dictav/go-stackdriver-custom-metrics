package sdcustom

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/api/monitoring/v3"
)

const (
	defaultInterval = 55 * time.Second
	minInterval     = 10 * time.Second
)

// MetricReporter struct
type MetricReporter struct {
	project    string
	zone       string
	metric     string
	instance   string
	monitoring *monitoring.Service
	value      chan int64
	timer      *time.Ticker
	logger     Logger
	interval   time.Duration
}

// NewMetricReporter creates a new MetricReporter
func NewMetricReporter(ctx context.Context, project, zone, metric, instance string) (*MetricReporter, error) {
	s, err := createService(ctx)
	if err != nil {
		return nil, err
	}

	l := log.New(os.Stderr, "", log.LstdFlags)

	v := make(chan int64, 1)
	v <- 0
	m := &MetricReporter{
		project:    project,
		zone:       zone,
		metric:     metric,
		instance:   instance,
		monitoring: s,
		value:      v,
		logger:     wrapLogger(l),
		interval:   defaultInterval,
	}

	if err := m.send(); err != nil {
		return nil, err
	}

	return m, nil
}

// Add metric value
func (m *MetricReporter) Add(n int64) {
	if m == nil {
		return
	}

	v := <-m.value
	m.value <- v + n
}

// Done reduce metric value
func (m *MetricReporter) Done(n int64) {
	if m == nil {
		return
	}

	v := <-m.value
	m.value <- v - n
}

// SetInterval for sending metric
func (m *MetricReporter) SetInterval(t time.Duration) {
	if t < minInterval {
		m.interval = minInterval
	} else {
		m.interval = t
	}
}

// Start reporting
func (m *MetricReporter) Start() {
	if m == nil {
		return
	}
	if m.timer == nil {
		t := m.interval
		if t < m.interval {
			t = minInterval
		}

		m.timer = time.NewTicker(t)
		go func() {
			for range m.timer.C {
				if err := m.send(); err != nil {
					m.logger.Print("Could not write time series value:", err)
				}
			}
		}()
	}
}

// Stop reporting
func (m *MetricReporter) Stop() {
	if m == nil || m.timer == nil {
		return
	}
	m.timer.Stop()
	m.timer = nil
}

// send a value for the custom metric created
func (m *MetricReporter) send() error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	v := <-m.value
	m.value <- v
	m.logger.Debug("send:", v)
	timeseries := monitoring.TimeSeries{
		Metric: &monitoring.Metric{
			Type: m.metric,
		},
		Resource: &monitoring.MonitoredResource{
			Labels: map[string]string{
				"instance_id": m.instance,
				"zone":        m.zone,
				"project_id":  m.project,
			},
			Type: "gce_instance",
		},
		Points: []*monitoring.Point{
			{
				Interval: &monitoring.TimeInterval{
					EndTime: now,
				},
				Value: &monitoring.TypedValue{
					Int64Value: &v,
				},
			},
		},
	}

	req := monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{&timeseries},
	}

	pid := projectResource(m.project)
	_, err := m.monitoring.Projects.TimeSeries.Create(pid, &req).Do()
	return err
}
