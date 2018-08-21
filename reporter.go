package sdcustom

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"google.golang.org/api/monitoring/v3"
)

const (
	defaultInterval = 55 * time.Second
	minInterval     = 10 * time.Second
)

// MetricReporter struct
type MetricReporter struct {
	ctx        context.Context
	project    string
	zone       string
	metric     string
	instance   string
	monitoring *monitoring.Service
	value      chan int64
	values     []int64
	timer      *time.Ticker
	logger     Logger
	interval   time.Duration
	mutex      sync.Mutex
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
		ctx:        ctx,
		project:    project,
		zone:       zone,
		metric:     customMetricPrefix + metric,
		instance:   instance,
		monitoring: s,
		value:      v,
		logger:     WrapLogger(l),
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
	m.logger.Debug("add:", v+n)
}

// Done reduce metric value
func (m *MetricReporter) Done(n int64) {
	if m == nil {
		return
	}

	v := <-m.value
	m.value <- v - n
	m.logger.Debug("done:", v-n)
}

// Set metric value
func (m *MetricReporter) Set(n int64) {
	if m == nil {
		return
	}

	<-m.value
	m.value <- n
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
	if m == nil || m.timer != nil {
		return
	}

	t := m.interval
	if t < m.interval {
		t = minInterval
	}

	m.timer = time.NewTicker(t)
	go func() {
		i := time.NewTicker(1 * time.Second)
		defer i.Stop()
		for {
			if m.timer == nil {
				return
			}

			select {
			case <-m.ctx.Done():
				m.Stop() // it is necessary to avoid memory leak
				return
			case <-m.timer.C:
				if err := m.send(); err != nil {
					m.logger.Print("Could not write time series value:", err)
				}
			case _ = <-i.C:
				v := <-m.value
				m.values = append(m.values, v)
				m.value <- v
			}
		}
	}()
}

// Stop reporting
func (m *MetricReporter) Stop() {
	if m == nil {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.timer == nil {
		return
	}
	m.timer.Stop()
	m.timer = nil // It is necessary to detect the end of ticker
}

// SetLogger set logger
func (m *MetricReporter) SetLogger(l Logger) {
	if l == nil {
		return
	}
	m.mutex.Lock()
	m.logger = l
	m.mutex.Unlock()
}

// send a value for the custom metric created
func (m *MetricReporter) send() error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	v := int64(0)
	if len(m.values) > 0 {
		// mutex.Lock is unnecessary because m.values is processed only in the same goroutine
		for _, n := range m.values {
			v += n
		}
		v /= int64(len(m.values))
		m.values = m.values[:0]
	}
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
