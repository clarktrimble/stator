// Package roster registers with discovery, repeatedly.
package roster

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/clarktrimble/hondo"

	"stator/roster/entity"
)

// Registrar specifies a registration interface.
type Registrar interface {
	Register(ctx context.Context, svc entity.Service) (err error)
	Unregister(ctx context.Context, svc entity.Service) (err error)
}

// Logger specifies a logging interface.
type Logger interface {
	Info(ctx context.Context, msg string, kv ...any)
	Error(ctx context.Context, msg string, err error, kv ...any)
	WithFields(ctx context.Context, kv ...any) context.Context
}

// ServiceConfig is configuration for service to be registered.
type ServiceConfig struct {
	Id          string   `json:"id" desc:"unique to service id" required:"true"`
	Name        string   `json:"name" desc:"name" required:"true"`
	Tags        []string `json:"tags" desc:"tags" required:"true"`
	MonitorSpec string   `json:"monitor_spec" desc:"specifier for monitor endpoint uri" default:"http://%s:%d/monitor"`
}

// Config is Roster configuration.
type Config struct {
	Interval time.Duration  `json:"reregister_interval" desc:"reregister period" default:"15m"`
	Service  *ServiceConfig `json:"service"`
}

// Roster repeatedly registers a service and unregisters when stopped.
type Roster struct {
	Registrar Registrar
	Logger    Logger
	Service   entity.Service
	Interval  time.Duration
}

// New creates a Roster from Config.
func (cfg *Config) New(ip string, port int, registrar Registrar, lgr Logger) *Roster {

	if ip == "lookup" {
		ip = getIp()
	}

	svc := entity.Service{
		Id:          cfg.Service.Id,
		Name:        cfg.Service.Name,
		Tags:        cfg.Service.Tags,
		IpAddress:   ip,
		Port:        port,
		MonitorSpec: cfg.Service.MonitorSpec,
	}

	return &Roster{
		Registrar: registrar,
		Logger:    lgr,
		Service:   svc,
		Interval:  cfg.Interval,
	}
}

// Start starts a Roster service.
func (roster *Roster) Start(ctx context.Context, wg *sync.WaitGroup) {

	err := roster.Service.Valid()
	if err != nil {
		roster.Logger.Error(ctx, "invalid service description", err)
		return
	}

	ctx = roster.Logger.WithFields(ctx, "worker_id", hondo.Rand(7))
	roster.Logger.Info(ctx, "worker starting", "name", "roster")

	roster.register(ctx)

	go roster.work(ctx, wg)
}

// unexported

func (roster *Roster) work(ctx context.Context, wg *sync.WaitGroup) {

	// open-loop, just keep registering ftw, but is simple

	wg.Add(1)
	defer wg.Done()

	tick := time.NewTicker(roster.Interval)

	for {
		select {
		case <-tick.C:
			roster.register(ctx)

		case <-ctx.Done():
			roster.Logger.Info(ctx, "worker shutting down")
			roster.unregister(ctx)
			roster.Logger.Info(ctx, "worker stopped")
			return
		}
	}
}

func (roster *Roster) register(ctx context.Context) {

	err := roster.Registrar.Register(ctx, roster.Service)
	if err != nil {
		roster.Logger.Error(ctx, "failed to register", err)
	}
}

func (roster *Roster) unregister(ctx context.Context) {

	ctx = context.WithoutCancel(ctx)

	err := roster.Registrar.Unregister(ctx, roster.Service)
	if err != nil {
		roster.Logger.Error(ctx, "failed to unregister", err)
	}
}

func getIp() (ip string) {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unable to determine external ip"
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
