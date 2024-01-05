package collector

import (
	"syscall"
	"time"

	"stator/entity"

	"github.com/pkg/errors"
)

type DiskUsage struct {
	Paths []string
}

func (du DiskUsage) Collect() (pa entity.PointsAt, err error) {

	// Todo: or pass now in??
	now := time.Now() // Todo: UTC somewhere above?

	points := []entity.Point{}
	for _, path := range du.Paths {

		fs := &syscall.Statfs_t{}
		err = syscall.Statfs(path, fs)
		if err != nil {
			err = errors.Wrapf(err, "failed to get disk usage for %s", path)
			return
		}

		bs := uint64(fs.Frsize)
		reserved := fs.Bfree - fs.Bavail

		total := (fs.Blocks - reserved) * bs
		available := fs.Bavail * bs
		percentUsed := float64(total-available) * 100 / float64(total)

		// Todo: account for multiple fs please!!!
		// name, desc, everything but labels and value are common across filesytems
		// a-and can be emitted together

		points = append(points, []entity.Point{
			{
				Name: "size_total",
				Desc: "Total size of the filesystem",
				Unit: "bytes",
				Type: "gauge",
				Labels: entity.Labels{
					{Key: "path", Val: path},
				},
				Value: entity.Uint{total},
			},
			{
				Name: "available",
				Desc: "Available space on the filesystem",
				Unit: "bytes",
				Type: "gauge",
				Labels: entity.Labels{
					{Key: "path", Val: path},
				},
				Value: entity.Uint{available},
			},
			{
				Name: "used",
				Desc: "Percentage of space on the filesystem in use",
				Unit: "percent",
				Type: "gauge",
				Labels: entity.Labels{
					{Key: "path", Val: path},
				},
				Value: entity.Float{percentUsed},
			},
		}...)
	}

	pa = entity.PointsAt{
		Name:  "du",
		Stamp: now,
		//Labels: labels,
		Points: points,
	}
	return
}

/*
// Datum is one of many.
type Datum struct {
	Labels    Labels
	Timestamp *time.Time
	Value     float64
}

// Field is named and typed data.
type Field struct {
	Name string
	Type string
	Data []Datum
}

// Fields converts Point to Fields.
func (pp *Point) Fields() (fields Fields) {

	fields = Fields{}
	for key, val := range pp.KeyVal {
		fields = append(fields, Field{
			Name: key,
			Type: pp.Type,
			Data: []Datum{{
				Labels:    pp.Labels,
				Timestamp: pp.Timestamp,
				Value:     val,
			}},
		})
	}

	return
}

  reservedBlocks := s.Bfree - s.Bavail
	info = Info{
		Total: uint64(s.Frsize) * (s.Blocks - reservedBlocks),
		Free:  uint64(s.Frsize) * s.Bavail,
		Files: s.Files,
		Ffree: s.Ffree,
		//nolint:unconvert
		FSType: getFSType(int64(s.Type)),
	}

func usage(path string) (disk DiskStatus) {

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}
*/
