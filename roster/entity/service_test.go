package entity

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEntity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Entity Suite")
}

var _ = Describe("Service", func() {
	var (
		svc Service
	)

	BeforeEach(func() {
		svc = Service{
			Id:          "123",
			Name:        "bargla",
			Tags:        []string{"one", "two"},
			IpAddress:   "1.2.3.4",
			Port:        8082,
			MonitorSpec: "http://%s:%d/monitor",
		}
	})

	Describe("getting a name id combo", func() {
		var (
			nid string
		)

		BeforeEach(func() {
			nid = svc.NameId()
		})

		When("all goes well", func() {
			It("sticks them together", func() {
				Expect(nid).To(Equal("bargla-123"))
			})
		})
	})

	Describe("checking validity", func() {
		var (
			err error
		)

		JustBeforeEach(func() {
			err = svc.Valid()
		})

		When("all goes well", func() {
			It("does not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		When("service is totally invalid", func() {
			BeforeEach(func() {
				svc = Service{}
			})

			It("errors with all the reasons", func() {
				Expect(err).To(MatchError("invalid Service: Id must not be blank,Name must not be blank,IpAddress failed to parse,Port must be between 1 and 65535,MonitorSpec must not be blank"))
			})
		})
	})
})
