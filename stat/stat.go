package stat

import (
	"bytes"
	"context"
	"net/http"

	"stator/entity"
)

// Logger is the internal logging interface.
type Logger interface {
	Error(ctx context.Context, msg string, err error, kv ...any)
}

// Formatter is the internal formatting interface.
type Formatter interface {
	Format(stats entity.PointsAt) (data []byte)
}

// Collector is the internal collection interface.
type Collector interface {
	Collect() (stats entity.PointsAt, err error)
	// Note: cache as needed _within_ any collector as needed!  yeah?
}

// Todo: left public for unit, will this become true?

// Service handles requests for stats
type Service struct {
	Logger     Logger
	Formatter  Formatter
	Collectors []Collector
}

// Handle handles http requests for stats
func (svc *Service) Handle(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()

	stats := svc.runCollectors(ctx)
	data := svc.format(stats)

	_, err := writer.Write(data)
	if err != nil {
		svc.Logger.Error(ctx, "failed to write stats to response", err)
	}

	// Todo: some statsy headers are good here?
	// Todo: responder is wanted? at all?
	// Todo: stats/points naming is approximate?
}

// unexported

func (svc *Service) runCollectors(ctx context.Context) (stats entity.Stats) {

	stats = entity.Stats{}

	for _, collector := range svc.Collectors {
		pts, err := collector.Collect()
		if err != nil {
			svc.Logger.Error(ctx, "failed to collect stats", err)
			continue
		}

		stats = append(stats, pts)
	}

	return
}

func (svc *Service) format(stats entity.Stats) []byte {

	buf := bytes.Buffer{}
	for _, pts := range stats {
		buf.Write(svc.Formatter.Format(pts))
	}

	return buf.Bytes()
}
