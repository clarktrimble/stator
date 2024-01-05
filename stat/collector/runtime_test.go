package collector_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stator/entity"
	. "stator/stat/collector"
)

func TestCollector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collector Suite")
}

var _ = Describe("Runtime", func() {
	var (
		rt    *Runtime
		stats entity.PointsAt
		err   error
	)

	BeforeEach(func() {
		rt = &Runtime{
			AppId: "boxie",
		}
	})

	Describe("collecting runtime stats", func() {
		BeforeEach(func() {
			stats, err = rt.Collect()
		})

		When("all goes well", func() {
			It("collects stats", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(stats.Name).To(Equal("gort"))
				Expect(stats.Points).To(HaveLen(9))
				// Todo: check moar
			})
		})
	})

})
