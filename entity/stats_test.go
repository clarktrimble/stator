package entity_test

// Note: tapping out to "_test" in order to dodge ginkgo's "Label"

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ste "stator/entity"
)

func TestEntity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Entity Suite")
}

var _ = Describe("Stats", func() {
	var (
		str string
	)

	Describe("formatting an unsigned integer value", func() {
		var (
			val ste.Uint
		)

		BeforeEach(func() {
			val = ste.Uint{Data: 99}
			str = val.String()
		})

		When("all goes well", func() {
			It("formats nicely", func() {
				Expect(str).To(Equal("99"))
			})
		})
	})

	Describe("formatting an floating point value", func() {
		var (
			val ste.Float
		)

		BeforeEach(func() {
			val = ste.Float{Data: 99.999999}
			str = val.String()
		})

		When("all goes well", func() {
			It("formats nicely", func() {
				Expect(str).To(Equal("100.00"))
			})
		})
	})

})
