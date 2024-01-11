// Package stator is a metrics service layer.
package stator

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"stator/collector/runtime"
	"stator/entity"
	"stator/formatter/prometheus"
)

//go:generate moq -out mock_test.go . Collector Formatter Router Logger

// Collector specifies a stats collector.
type Collector interface {
	Collect(time.Time) (stats entity.PointsAt, err error)
	// Note: cache as needed _within_ any collector as needed!
}

// Formatter specifies a stats formatter.
type Formatter interface {
	Format(stats entity.PointsAt) (data []byte)
}

// Router specifies an http router.
type Router interface {
	HandleFunc(pattern string, handler http.HandlerFunc)
}

// Logger specifies a logger.
type Logger interface {
	Error(ctx context.Context, msg string, err error, kv ...any)
}

// Svc handles requests for stats
type Svc struct {
	Collectors []Collector
	Formatter  Formatter
	Logger     Logger
}

// ExposeRuntime is a convienience function that creates a stats service
// that collects runtime stats and exposes them via "/metrics" in prometheus format.
func ExposeRuntime(appId, runId string, rtr Router, lgr Logger) (svc *Svc) {

	svc = &Svc{
		Collectors: []Collector{
			&runtime.Runtime{AppId: appId, RunId: runId},
		},
		Formatter: prometheus.Prometheus{},
		Logger:    lgr,
	}

	rtr.HandleFunc("GET /metrics", svc.GetStats)

	return
}

// AddCollector adds a collector.
func (svc *Svc) AddCollector(collector Collector) {

	svc.Collectors = append(svc.Collectors, collector)
}

// GetStats handles http requests for stats
func (svc *Svc) GetStats(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()

	stats := svc.runCollectors(ctx)
	data := svc.format(stats)

	_, err := writer.Write(data)
	if err != nil {
		svc.Logger.Error(ctx, "failed to write stats to response", err)
	}
}

// unexported

func (svc *Svc) runCollectors(ctx context.Context) (stats entity.Stats) {

	stats = entity.Stats{}
	now := time.Now()

	for _, collector := range svc.Collectors {
		pts, err := collector.Collect(now)
		if err != nil {
			svc.Logger.Error(ctx, "failed to collect stats", err)
			continue
		}

		stats = append(stats, pts)
	}

	return
}

func (svc *Svc) format(stats entity.Stats) []byte {

	buf := bytes.Buffer{}
	for _, pa := range stats {
		buf.Write(svc.Formatter.Format(pa))
	}

	return buf.Bytes()
}
