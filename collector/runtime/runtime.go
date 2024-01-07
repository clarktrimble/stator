package runtime

import (
	"os"
	"runtime/metrics"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"stator/entity"
)

const (
	name = "gort"
)

var (
	collectible = []metric{
		{
			long: "/cpu/classes/total:cpu-seconds",
			name: "cpu_total",
			unit: "seconds",
			desc: "Total available CPU time",
		},
		{
			long: "/cpu/classes/user:cpu-seconds",
			name: "cpu_user",
			unit: "seconds",
			desc: "CPU time spent running user Go code",
		},
		{
			long: "/cpu/classes/idle:cpu-seconds",
			name: "cpu_idle",
			unit: "seconds",
			desc: "Unused available CPU time",
		},
		{
			long: "/cpu/classes/gc/total:cpu-seconds",
			name: "cpu_gc",
			unit: "seconds",
			desc: "CPU time spent performing GC tasks",
		},
		{
			long: "/memory/classes/total:bytes",
			name: "mem_total",
			unit: "bytes",
			desc: "All memory mapped into the current process",
		},
		{
			long: "/memory/classes/heap/objects:bytes",
			name: "mem_heap",
			unit: "bytes",
			desc: "Memory occupied by live and yet to be marked free objects",
		},
		{
			long: "/memory/classes/heap/stacks:bytes",
			name: "mem_stack",
			unit: "bytes",
			desc: "Memory allocated from the heap that is reserved for stack space",
		},
		{
			long: "/sched/goroutines:goroutines",
			name: "goroutines",
			unit: "count",
			desc: "Count of live goroutines",
		},
		{
			long: "/sync/mutex/wait/total:seconds",
			name: "mutex_wait",
			unit: "seconds",
			desc: "Time spent blocked on a mutex",
		},
	}
)

// Runtime collects go runtime stats.
type Runtime struct {
	AppId string
	RunId string
}

// Collect collects stats.
func (rt Runtime) Collect(ts time.Time) (pa entity.PointsAt, err error) {

	samples := make([]metrics.Sample, len(collectible))
	for i := range collectible {
		samples[i].Name = collectible[i].long
	}

	metrics.Read(samples)

	points, err := toPoints(samples)
	if err != nil {
		return
	}

	pa = entity.PointsAt{
		Name:  name,
		Stamp: ts,
		Labels: entity.Labels{
			{Key: "app_id", Val: rt.AppId},
			{Key: "run_id", Val: rt.RunId},
			{Key: "process_id", Val: strconv.Itoa(os.Getpid())},
		},
		Points: points,
	}
	return
}

// unexported

type metric struct {
	long string
	name string
	unit string
	desc string
}

func toPoints(samples []metrics.Sample) (points []entity.Point, err error) {

	points = make([]entity.Point, len(collectible))
	for i, sample := range samples {

		var value entity.Value
		switch sample.Value.Kind() {
		case metrics.KindUint64:
			value = entity.Uint{Data: sample.Value.Uint64()}
		case metrics.KindFloat64:
			value = entity.Float{Data: sample.Value.Float64()}
		default:
			err = errors.Errorf("unknown go runtime stat type for: %s", sample.Name)
			return
		}

		points[i] = entity.Point{
			Name:  collectible[i].name,
			Desc:  collectible[i].desc,
			Unit:  collectible[i].unit,
			Type:  "gauge",
			Value: value,
		}
	}

	return
}
