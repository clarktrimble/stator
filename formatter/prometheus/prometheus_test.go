package prometheus

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stator/entity"
)

func TestPrometheus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Prometheus Suite")
}

var _ = Describe("Prometheus", func() {

	Describe("formatting stats", func() {

		var (
			pa  entity.PointsAt
			om  Prometheus
			out []byte
		)

		BeforeEach(func() {

			pa = entity.PointsAt{
				Name:   "common",
				Stamp:  time.Time{},
				Labels: entity.Labels{{Key: "cid", Val: "valero"}},
				Points: []entity.Point{
					{
						Name:   "particular",
						Desc:   "Dummy point for test.",
						Unit:   "bytes",
						Type:   "gauge",
						Labels: entity.Labels{{Key: "path", Val: "/boot"}},
						Value:  entity.Uint{Data: 99},
					},
					{
						Name:   "particular",
						Desc:   "Dummy point for test.",
						Unit:   "bytes",
						Type:   "gauge",
						Labels: entity.Labels{{Key: "path", Val: "/different"}},
						Value:  entity.Uint{Data: 9999},
					},
				},
			}

			out = om.Format(pa)
		})

		When("all goes well", func() {
			It("formats them with aplomb", func() {
				Expect(string(out)).To(Equal(expected))
			})
		})
	})
})

var expected = `
# HELP common_particular_bytes Dummy point for test.
# TYPE common_particular_bytes gauge
common_particular_bytes{cid="valero",path="/boot"} 99 -62135596800000
common_particular_bytes{cid="valero",path="/different"} 9999 -62135596800000
`
