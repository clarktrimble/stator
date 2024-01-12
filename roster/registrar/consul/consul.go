// Package consul provides for registration with Consul agent.
package consul

import (
	"context"
	"fmt"
	"time"

	"stator/roster/entity"
)

//go:generate moq -out mock_test.go . Client

const (
	registerPath   string = "/v1/agent/service/register"
	unregisterPath string = "/v1/agent/service/deregister/%s"
)

// Client specifies an http client.
type Client interface {
	SendObject(ctx context.Context, method, path string, snd, rcv any) (err error)
}

// Config is Consul configuration.
type Config struct {
	CheckInterval   time.Duration `json:"check_interval" desc:"health check period" default:"1m"`
	CheckTimeout    time.Duration `json:"check_timeout" desc:"health check timeout" default:"10s"`
	DeregisterAfter time.Duration `json:"deregister_after" desc:"deregister after failed check period" default:"30m"`
}

// Consul is a Consul client.
type Consul struct {
	Client          Client
	CheckInterval   time.Duration
	CheckTimeout    time.Duration
	DeregisterAfter time.Duration
}

// New creates a Consul from Config.
func (cfg *Config) New(client Client) *Consul {

	return &Consul{
		Client:          client,
		CheckInterval:   cfg.CheckInterval,
		CheckTimeout:    cfg.CheckTimeout,
		DeregisterAfter: cfg.DeregisterAfter,
	}
}

// Register registers the service.
func (csl *Consul) Register(ctx context.Context, svc entity.Service) (err error) {

	// Todo: note about catalog reg
	// Todo: support reregister, checking first?
	//       or could long-poll

	reg := register{
		ID:      svc.NameId(),
		Name:    svc.Name,
		Tags:    svc.Tags,
		Address: svc.IpAddress,
		Port:    svc.Port,
		Check: check{
			HTTP:                           fmt.Sprintf(svc.MonitorSpec, svc.IpAddress, svc.Port),
			Interval:                       csl.CheckInterval.String(),
			Timeout:                        csl.CheckTimeout.String(),
			DeregisterCriticalServiceAfter: csl.DeregisterAfter.String(),
			Status:                         "passing",
		},
	}

	err = csl.Client.SendObject(ctx, "PUT", registerPath, reg, nil)
	return
}

// Unregister deregisters the service.
func (csl *Consul) Unregister(ctx context.Context, svc entity.Service) (err error) {

	err = csl.Client.SendObject(ctx, "PUT", fmt.Sprintf(unregisterPath, svc.NameId()), nil, nil)
	return
}

// unexported

type check struct {
	Status                         string
	HTTP                           string
	Interval                       string
	Timeout                        string
	DeregisterCriticalServiceAfter string
}

type register struct {
	ID      string
	Name    string
	Tags    []string
	Address string
	Port    int
	Check   check
}
