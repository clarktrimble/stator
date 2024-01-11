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
		cfg   *Config
		du    *DiskUsage
		stats entity.PointsAt
		err   error
	)

	BeforeEach(func() {
		cfg = &Config{
			Paths: []string{"/"},
		}

		du = cfg.New()
	})

	Describe("creating a collector", func() {
		It("collects stats", func() {
			Expect(du).To(Equal(&DiskUsage{
				Paths: []string{"/"},
			}))
		})
	})

	Describe("collecting stats", func() {

		JustBeforeEach(func() {
			stats, err = du.Collect(time.Time{})
		})

		When("all goes well", func() {
			It("collects stats", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(stats.Name).To(Equal("du"))
				Expect(stats.Points).To(HaveLen(3))
			})
		})

		When("no such filesystem", func() {
			BeforeEach(func() {
				du.Paths = []string{"/", "/bargle"}
			})

			It("returns error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

	})

})
