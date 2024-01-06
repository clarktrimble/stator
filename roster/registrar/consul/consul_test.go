package consul

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stator/roster/entity"
)

func TestConsul(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Consul Suite")
}

var _ = Describe("Consul", func() {
	var (
		cfg    *Config
		client *ClientMock
		csl    *Consul
	)

	BeforeEach(func() {

		client = &ClientMock{
			SendObjectFunc: func(ctx context.Context, method string, path string, snd any, rcv any) error {
				return nil
			},
		}
		cfg = &Config{
			CheckInterval:   time.Minute,
			CheckTimeout:    10 * time.Second,
			DeregisterAfter: 30 * time.Minute,
		}

		csl = cfg.New(client)
	})

	Describe("creating a registrar", func() {

		When("all goes well", func() {
			It("creates one with client and cfg durations", func() {
				Expect(csl).To(Equal(&Consul{
					Client:          client,
					CheckInterval:   time.Minute,
					CheckTimeout:    10 * time.Second,
					DeregisterAfter: 30 * time.Minute,
				}))
			})
		})
	})

	Describe("interacting with discovery", func() {
		var (
			ctx context.Context
			svc entity.Service
			err error
		)

		BeforeEach(func() {
			ctx = context.Background()
			svc = entity.Service{
				Id:          "123",
				Name:        "foobear",
				Tags:        []string{"one", "two"},
				IpAddress:   "1.2.3.4",
				Port:        8082,
				MonitorSpec: "http://%s:%d/monitor",
			}
		})

		Describe("registering", func() {

			JustBeforeEach(func() {
				err = csl.Register(ctx, svc)
			})

			When("all goes well", func() {
				It("calls the client", func() {
					Expect(err).ToNot(HaveOccurred())

					soc := client.SendObjectCalls()
					Expect(soc).To(HaveLen(1))
					Expect(soc[0].Ctx).To(Equal(ctx))
					Expect(soc[0].Method).To(Equal("PUT"))
					Expect(soc[0].Path).To(Equal("/v1/agent/service/register"))
					Expect(soc[0].Snd).To(Equal(register{
						ID:      "foobear-123",
						Name:    "foobear",
						Tags:    []string{"one", "two"},
						Address: "1.2.3.4",
						Port:    8082,
						Check: check{
							Status:                         "passing",
							HTTP:                           "http://1.2.3.4:8082/monitor",
							Interval:                       "1m0s",
							Timeout:                        "10s",
							DeregisterCriticalServiceAfter: "30m0s",
						},
					}))
					Expect(soc[0].Rcv).To(BeNil())
				})
			})

			When("client has trouble", func() {
				BeforeEach(func() {
					client.SendObjectFunc = func(ctx context.Context, method string, path string, snd any, rcv any) error {
						return fmt.Errorf("oops")
					}
				})

				It("relays the error", func() {
					Expect(err).To(MatchError("oops"))
				})
			})

		})

		Describe("unregistering", func() {

			JustBeforeEach(func() {
				err = csl.Unregister(ctx, svc)
			})

			When("all goes well", func() {
				It("calls the client", func() {
					Expect(err).ToNot(HaveOccurred())

					soc := client.SendObjectCalls()
					Expect(soc).To(HaveLen(1))
					Expect(soc[0].Ctx).To(Equal(ctx))
					Expect(soc[0].Method).To(Equal("PUT"))
					Expect(soc[0].Path).To(Equal("/v1/agent/service/deregister/foobear-123"))
					Expect(soc[0].Snd).To(BeNil())
					Expect(soc[0].Rcv).To(BeNil())
				})
			})

			When("client has trouble", func() {
				BeforeEach(func() {
					client.SendObjectFunc = func(ctx context.Context, method string, path string, snd any, rcv any) error {
						return fmt.Errorf("oops")
					}
				})

				It("relays the error", func() {
					Expect(err).To(MatchError("oops"))
				})
			})

		})

	})

})
