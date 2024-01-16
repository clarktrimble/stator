package stator

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stator/collector/runtime"
	"stator/entity"
	"stator/formatter/prometheus"
)

func TestStat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stat Suite")
}

var _ = Describe("Stat", func() {
	var (
		lgr *LoggerMock
		svc *Svc
	)

	BeforeEach(func() {

		lgr = &LoggerMock{
			ErrorFunc: func(ctx context.Context, msg string, err error, kv ...any) {},
		}
	})

	Describe("exposing runtime stats", func() {
		var (
			appId string
			runId string
			rtr   *RouterMock
		)

		BeforeEach(func() {
			appId = "bargla"
			runId = "456"
			rtr = &RouterMock{
				HandleFuncFunc: func(pattern string, handler http.HandlerFunc) {},
			}

			svc = ExposeRuntime(appId, runId, rtr, lgr)
		})

		When("all goes well", func() {
			It("creates the service and registers route", func() {
				Expect(svc).To(Equal(&Svc{
					Collectors: []Collector{
						&runtime.Runtime{AppId: "bargla", RunId: "456"},
					},
					Formatter: prometheus.Prometheus{},
					Logger:    lgr,
				}))

				Expect(rtr.HandleFuncCalls()).To(HaveLen(1))
				Expect(rtr.HandleFuncCalls()[0].Pattern).To(Equal("GET /metrics"))
			})
		})
	})

	Describe("adding a collector", func() {
		BeforeEach(func() {
			svc = &Svc{}
			svc.AddCollector(&CollectorMock{})
		})

		When("all goes well", func() {
			It("is appended to Collectors", func() {
				Expect(svc.Collectors).To(HaveLen(1))
			})
		})
	})

	Describe("handling a request for stats", func() {
		var (
			collOne *CollectorMock
			collTwo *CollectorMock
			fmtr    *FormatterMock

			writer  http.ResponseWriter
			request *http.Request
		)

		BeforeEach(func() {
			collOne = &CollectorMock{
				CollectFunc: func(timeMoqParam time.Time) (entity.PointsAt, error) {
					return entity.PointsAt{}, fmt.Errorf("oops")
				},
			}
			collTwo = &CollectorMock{
				CollectFunc: func(timeMoqParam time.Time) (entity.PointsAt, error) {
					return entity.PointsAt{}, nil
				},
			}

			fmtr = &FormatterMock{
				FormatFunc: func(stats entity.PointsAt) []byte {
					return []byte("stuff")
				},
			}

			svc = &Svc{
				Collectors: []Collector{
					collOne,
					collTwo,
				},
				Formatter: fmtr,
				Logger:    lgr,
			}

			request = &http.Request{}
		})

		JustBeforeEach(func() {
			svc.GetStats(writer, request)
		})

		When("all goes well, mostly", func() {
			BeforeEach(func() {
				writer = httptest.NewRecorder()
			})

			It("collects, logging any errors, formats, and writes to response", func() {

				Expect(collOne.CollectCalls()).To(HaveLen(1))
				Expect(collTwo.CollectCalls()).To(HaveLen(1))

				Expect(lgr.ErrorCalls()).To(HaveLen(1))
				Expect(lgr.ErrorCalls()[0].Msg).To(Equal("failed to collect stats"))

				Expect(fmtr.FormatCalls()).To(HaveLen(1))

				recorder, ok := writer.(*httptest.ResponseRecorder)
				Expect(ok).To(BeTrue())
				Expect(recorder.Body.String()).To(Equal("stuff"))
			})
		})

		When("write to response fails", func() {
			BeforeEach(func() {
				writer = &errorResponder{}
				request = &http.Request{}
			})

			It("logs another error", func() {
				Expect(lgr.ErrorCalls()).To(HaveLen(2))
				Expect(lgr.ErrorCalls()[1].Msg).To(Equal("failed to write stats to response"))
			})
		})

	})
})

type errorResponder struct{}

func (er *errorResponder) Header() (hdr http.Header) {
	return http.Header{}
}
func (er *errorResponder) Write(body []byte) (count int, err error) {
	return 0, fmt.Errorf("oops")
}
func (er *errorResponder) WriteHeader(status int) {}
