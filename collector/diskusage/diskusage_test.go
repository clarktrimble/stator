package diskusage

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stator/entity"
)

func TestDiskUsage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DiskUsage Suite")
}

var _ = Describe("DiskUsage", func() {
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
			stats, err = du.Collect(time.Time{})
		})

		When("all goes well", func() {
			It("collects stats", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(stats.Name).To(Equal("du"))
				Expect(stats.Points).To(HaveLen(3))
			})
		})
	})

})
