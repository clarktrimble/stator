package wave

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stator/entity"
)

func TestWave(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wave Suite")
}

var _ = Describe("Wave", func() {
	var (
		wv  *Wave
		pa  entity.PointsAt
		err error
	)

	BeforeEach(func() {
		wv = New()
	})

	Describe("creating a collecotr", func() {
		When("all goes well", func() {
			It("creates one", func() {
				Expect(wv.count).To(Equal(0))
				Expect(wv.series).To(HaveLen(3))
			})
		})
	})

	Describe("collecting stats", func() {
		BeforeEach(func() {
			pa, err = wv.Collect(time.Time{})
		})

		When("all goes well", func() {
			It("collects stats", func() {
				Expect(err).ToNot(HaveOccurred())

				Expect(pa.Name).To(Equal("wave"))
				Expect(pa.Points).To(HaveLen(3))
			})
		})
	})

})
