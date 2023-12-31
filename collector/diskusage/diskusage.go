package diskusage

import (
	"syscall"
	"time"

	"github.com/pkg/errors"

	"stator/entity"
)

const (
	name = "du"
)

// DiskUsage collects disk usage stats.
type DiskUsage struct {
	Paths []string
}

// Collect collects stats.
func (du DiskUsage) Collect(ts time.Time) (pa entity.PointsAt, err error) {

	pa = entity.PointsAt{
		Name:   name,
		Stamp:  ts,
		Points: []entity.Point{},
	}

	for _, path := range du.Paths {

		var size, avail uint64
		var used float64

		size, avail, used, err = duStats(path)
		if err != nil {
			return
		}

		labels := entity.Labels{{Key: "path", Val: path}}

		pa.Points = append(pa.Points, []entity.Point{
			{
				Name:   "size",
				Desc:   "Total size of the filesystem",
				Unit:   "bytes",
				Type:   "gauge",
				Labels: labels,
				Value:  entity.Uint{Data: size},
			},
			{
				Name:   "available",
				Desc:   "Available space on the filesystem",
				Unit:   "bytes",
				Type:   "gauge",
				Labels: labels,
				Value:  entity.Uint{Data: avail},
			},
			{
				Name:   "used",
				Desc:   "Percentage of space on the filesystem in use",
				Unit:   "percent",
				Type:   "gauge",
				Labels: labels,
				Value:  entity.Float{Data: used},
			},
		}...)
	}

	return
}

// unexported

func duStats(path string) (size, avail uint64, used float64, err error) {

	fs := &syscall.Statfs_t{}
	err = syscall.Statfs(path, fs)
	if err != nil {
		err = errors.Wrapf(err, "failed to get disk usage for %s", path)
		return
	}

	bs := uint64(fs.Frsize)
	size = fs.Blocks * bs
	avail = fs.Bavail * bs

	reserved := fs.Bfree - fs.Bavail
	total := (fs.Blocks - reserved) * bs
	used = float64(total-avail) * 100 / float64(total)

	return
}
