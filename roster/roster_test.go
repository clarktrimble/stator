package roster

import (
	"context"
	"fmt"
	"net"
	"stator/roster/entity"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRoster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Roster Suite")
}

var _ = Describe("Roster", func() {
	var (
		svc       entity.Service
		cfg       *Config
		port      int
		registrar *RegistrarMock
		lgr       *LoggerMock
		roster    *Roster
	)

	BeforeEach(func() {
		svc = entity.Service{
			Id:          "123",
			Name:        "bargla",
			Tags:        []string{"one", "two"},
			IpAddress:   "1.2.3.4",
			Port:        8082,
			MonitorSpec: "http://%s:%d/monitor",
		}
		cfg = &Config{
			Interval: 15 * time.Minute,
			Service: &ServiceConfig{
				Id:          "123",
				Name:        "bargla",
				Tags:        []string{"one", "two"},
				IpAddress:   "1.2.3.4",
				MonitorSpec: "http://%s:%d/monitor",
			},
		}

		port = 8082

		registrar = &RegistrarMock{
			RegisterFunc: func(ctx context.Context, svc entity.Service) error {
				return nil
			},
			UnregisterFunc: func(ctx context.Context, svc entity.Service) error {
				return nil
			},
		}

		lgr = &LoggerMock{
			InfoFunc:  func(ctx context.Context, msg string, kv ...any) {},
			ErrorFunc: func(ctx context.Context, msg string, err error, kv ...any) {},
			WithFieldsFunc: func(ctx context.Context, kv ...interface{}) context.Context {
				return ctx
			},
		}

		roster = cfg.New(port, registrar, lgr)
	})

	Describe("creating a roster", func() {

		When("all goes well", func() {
			It("creates one", func() {
				Expect(roster).To(Equal(&Roster{
					Registrar: registrar,
					Logger:    lgr,
					Service:   svc,
					Interval:  15 * time.Minute,
				}))
			})
		})

		When("looking up ip address", func() {
			var (
				rosterToo *Roster
			)

			BeforeEach(func() {
				cfg.Service.IpAddress = "lookup"
				rosterToo = cfg.New(port, registrar, lgr)
			})

			It("finds one", func() {
				Expect(net.ParseIP(rosterToo.Service.IpAddress)).ToNot(BeNil())
			})
		})
	})

	Describe("starting a roster", func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
			wg     sync.WaitGroup
		)

		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())
			roster.Interval = 100 * time.Millisecond
		})

		JustBeforeEach(func() {
			roster.Start(ctx, &wg)
		})

		When("all goes well", func() {
			It("registers, re-registers periodically, and shuts down when cancelled", func() {

				ic := lgr.InfoCalls
				rc := registrar.RegisterCalls
				uc := registrar.UnregisterCalls

				Expect(lgr.WithFieldsCalls()).To(HaveLen(1))
				Expect(lgr.WithFieldsCalls()[0].Kv[0]).To(Equal("worker_id"))

				Expect(ic()).To(HaveLen(1))
				Expect(ic()[0].Msg).To(Equal("worker starting"))

				Expect(rc()).To(HaveLen(1))
				Expect(rc()[0].Ctx).To(Equal(ctx))
				Expect(rc()[0].Svc).To(Equal(svc))

				Eventually(rc).Should(HaveLen(2))

				cancel()
				wg.Wait()

				Eventually(ic).Should(HaveLen(3))
				Expect(ic()[1].Msg).To(Equal("worker shutting down"))
				Expect(ic()[2].Msg).To(Equal("worker stopped"))

				Expect(uc()).Should(HaveLen(1))
				Expect(uc()[0].Ctx).To(Equal(context.WithoutCancel(ctx)))
				Expect(uc()[0].Svc).To(Equal(svc))
			})
		})

		When("registrar errors", func() {
			BeforeEach(func() {
				registrar.RegisterFunc = func(ctx context.Context, svc entity.Service) error {
					return fmt.Errorf("error from reg")
				}
				registrar.UnregisterFunc = func(ctx context.Context, svc entity.Service) error {
					return fmt.Errorf("error from unreg")
				}
			})

			It("logs errors and keeps trying", func() {

				ic := lgr.InfoCalls
				ec := lgr.ErrorCalls
				rc := registrar.RegisterCalls
				uc := registrar.UnregisterCalls

				Expect(lgr.WithFieldsCalls()).To(HaveLen(1))
				Expect(lgr.WithFieldsCalls()[0].Kv[0]).To(Equal("worker_id"))

				Expect(ic()).To(HaveLen(1))
				Expect(ic()[0].Msg).To(Equal("worker starting"))

				Expect(ec()).To(HaveLen(1))
				Expect(ec()[0].Msg).To(Equal("failed to register"))

				Expect(rc()).To(HaveLen(1))
				Expect(rc()[0].Ctx).To(Equal(ctx))
				Expect(rc()[0].Svc).To(Equal(svc))

				Eventually(rc).Should(HaveLen(2))

				Expect(ec()).To(HaveLen(2))
				Expect(ec()[1].Msg).To(Equal("failed to register"))

				cancel()
				wg.Wait()

				Eventually(ic).Should(HaveLen(3))
				Expect(ic()[1].Msg).To(Equal("worker shutting down"))
				Expect(ic()[2].Msg).To(Equal("worker stopped"))

				Expect(uc()).Should(HaveLen(1))
				Expect(uc()[0].Ctx).To(Equal(context.WithoutCancel(ctx)))
				Expect(uc()[0].Svc).To(Equal(svc))

				Expect(ec()).To(HaveLen(3))
				Expect(ec()[2].Msg).To(Equal("failed to unregister"))
			})
		})

		When("svc is invalid", func() {
			BeforeEach(func() {
				roster.Service.IpAddress = ""
			})

			It("logs an error and does not reg", func() {
				Expect(lgr.ErrorCalls()).To(HaveLen(1))
				Expect(lgr.ErrorCalls()[0].Msg).To(Equal("worker abort"))

				Expect(registrar.RegisterCalls()).To(HaveLen(0))
			})
		})

	})

})
