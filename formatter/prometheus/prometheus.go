package prometheus

import (
	"bytes"
	"fmt"
	"strings"

	"stator/entity"
)

// Prometheus formats stats for consumption by Prometheus.
//
// For example:
// # TYPE http_requests_total counter
// # HELP http_requests_total The total number of HTTP requests.
// http_requests_total{method="post",code="200"} 1027 1395066363000
// http_requests_total{method="post",code="400"}    3 1395066363000
//
// In the spirit of: https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format
//
// The following advice will be applicable at some scale?
// https://prometheus.io/docs/instrumenting/writing_exporters/#target-labels-not-static-scraped-labels
type Prometheus struct {
}

// Format formats stats.
func (om Prometheus) Format(pa entity.PointsAt) []byte {

	out := map[string][]string{}
	ordered := []string{}

	// gather data under common header

	for i := range pa.Points {
		hdr, dtm := headerDatum(pa, i)

		data, ok := out[hdr]
		if ok {
			data = append(data, dtm)
		} else {
			ordered = append(ordered, hdr)
			data = []string{dtm}
		}
		out[hdr] = data
	}

	// buffer in original order

	var buf bytes.Buffer
	for _, hdr := range ordered {
		buf.WriteString(hdr)
		for _, dtm := range out[hdr] {
			buf.WriteString(dtm)
		}
	}

	return buf.Bytes()
}

// unexported

func headerDatum(pa entity.PointsAt, idx int) (hdr, dtm string) {

	pt := pa.Points[idx]
	lbl := label(append(pa.Labels, pt.Labels...))

	name := fmt.Sprintf("%s_%s", pa.Name, pt.Name)
	if pt.Unit != "" {
		name = fmt.Sprintf("%s_%s", name, pt.Unit)
	}

	hdr = header(name, pt)
	dtm = fmt.Sprintf("%s{%s} %s %d\n", name, lbl, pt.Value, pa.Stamp.UnixMilli())
	return
}

func header(name string, pt entity.Point) string {

	builder := &strings.Builder{}

	fmt.Fprintf(builder, "\n")
	fmt.Fprintf(builder, "# HELP %s %s\n", name, pt.Desc)
	fmt.Fprintf(builder, "# TYPE %s %s\n", name, pt.Type)

	return builder.String()
}

func label(labels entity.Labels) string {

	strs := []string{}
	for _, label := range labels {
		strs = append(strs, fmt.Sprintf(`%s="%s"`, label.Key, label.Val))
	}

	return strings.Join(strs, ",")
}
