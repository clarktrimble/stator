package collector_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stator/entity"
	. "stator/stat/collector"
)

var _ = Describe("Collector", func() {
	var (
		du    *DiskUsage
		stats entity.PointsAt
		err   error
	)

	BeforeEach(func() {
		du = &DiskUsage{
			Paths: []string{"/"},
		}
	})

	Describe("collecting ...", func() {
		BeforeEach(func() {
			stats, err = du.Collect()
		})

		When("all goes well", func() {
			It("collects stats", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(stats.Name).To(Equal("du"))
				Expect(stats.Points).To(HaveLen(3))
				// Todo: check moar
			})
		})
	})

})
