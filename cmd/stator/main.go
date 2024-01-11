package main

import (
	"context"
	"os"
	"stator"
	"sync"

	"github.com/clarktrimble/delish"
	"github.com/clarktrimble/delish/graceful"
	"github.com/clarktrimble/giant"
	"github.com/clarktrimble/hondo"
	"github.com/clarktrimble/launch"
	"github.com/clarktrimble/sabot"

	"stator/collector/diskusage"
	"stator/collector/wave"
	"stator/minroute"
	"stator/roster"
	"stator/roster/registrar/consul"
)

const (
	appId     string = "stator"
	cfgPrefix string = "sttr"
	blerb     string = "'stator' demonstrates registration and reporting of stats"
)

var (
	version string
	wg      sync.WaitGroup
)

type Config struct {
	Version string         `json:"version" ignored:"true"`
	Logger  *sabot.Config  `json:"logger"`
	Client  *giant.Config  `json:"consul_http_client"`
	Consul  *consul.Config `json:"consul"`
	Roster  *roster.Config `json:"roster"`
	Server  *delish.Config `json:"http_server"`
}

func main() {

	// load config and setup logger

	cfg := &Config{Version: version}
	launch.Load(cfg, cfgPrefix, blerb)

	lgr := cfg.Logger.New(os.Stdout)
	runId := hondo.Rand(7)
	ctx := lgr.WithFields(context.Background(), "app_id", appId, "run_id", runId)
	lgr.Info(ctx, "starting up", "config", cfg)

	// init graceful and create router

	ctx = graceful.Initialize(ctx, &wg, lgr)

	rtr := minroute.New(ctx, lgr)
	rtr.HandleFunc("GET /config", delish.ObjHandler("config", cfg, lgr))
	rtr.HandleFunc("GET /monitor", delish.ObjHandler("status", "ok", lgr))

	// setup and start registration

	client := cfg.Client.NewWithTrippers(lgr)
	csl := cfg.Consul.New(client)
	rstr := cfg.Roster.New(cfg.Server.Port, csl, lgr)
	rstr.Start(ctx, &wg)

	// setup stats expositor

	svc := stator.ExposeRuntime(appId, runId, rtr, lgr)
	svc.AddCollector(&diskusage.DiskUsage{Paths: []string{"/", "/boot/efi"}})
	svc.AddCollector(wave.New())

	// start api server and wait for shutdown

	server := cfg.Server.NewWithLog(ctx, rtr, lgr)
	server.Start(ctx, &wg)
	graceful.Wait(ctx)
}
