package stat

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"stator/stat/entity"
)

// Logger specifies a logger.
type Logger interface {
	Error(ctx context.Context, msg string, err error, kv ...any)
}

// Formatter specifies a stats formatter.
type Formatter interface {
	Format(stats entity.PointsAt) (data []byte)
}

// Collector specifies a stats collector.
type Collector interface {
	Collect(time.Time) (stats entity.PointsAt, err error)
	// Note: cache as needed _within_ any collector as needed!
}

// Svc handles requests for stats
type Svc struct {
	Logger     Logger
	Formatter  Formatter
	Collectors []Collector
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
