package minroute

import (
	"context"
	"net/http"
	"strings"

	"github.com/clarktrimble/delish/respond"
	"github.com/pkg/errors"
)

// Logger specifies a logging interface.
type Logger interface {
	Info(ctx context.Context, msg string, kv ...any)
	Error(ctx context.Context, msg string, err error, kv ...any)
	WithFields(ctx context.Context, kv ...any) context.Context
}

// MinRoute maps http methods and paths to handlers
type MinRoute struct {
	Logger Logger
	routes map[string]map[string]http.HandlerFunc
}

// New creates a router with an empty route table
func New(lgr Logger) (rtr *MinRoute) {

	rtr = &MinRoute{
		Logger: lgr,
		routes: map[string]map[string]http.HandlerFunc{
			"GET":    map[string]http.HandlerFunc{},
			"PUT":    map[string]http.HandlerFunc{},
			"POST":   map[string]http.HandlerFunc{},
			"DELETE": map[string]http.HandlerFunc{},
		},
	}

	return
}

// ServeHTTP looks up the handler and invokes
func (rtr *MinRoute) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	routes, ok := rtr.routes[request.Method]
	if !ok {
		notFound(request.Context(), writer, rtr.Logger)
		return
	}

	handler, ok := routes[request.RequestURI]
	if !ok {
		notFound(request.Context(), writer, rtr.Logger)
		return
	}

	handler.ServeHTTP(writer, request)
}

// Set associates a method and path with a handler
func (rtr *MinRoute) HandleFunc(pattern string, handler http.HandlerFunc) {

	// Todo: unit!
	split := strings.Split(pattern, " ")
	if len(split) != 2 {
		panic(errors.Errorf("failed to split pattern: '%s' into method and path", pattern))
	}
	method := split[0]
	path := split[1]

	_, ok := rtr.routes[method]
	if !ok {
		panic(errors.Errorf("unsupported method from pattern: '%s'", pattern))
	}

	rtr.routes[method][path] = handler
}

// unexported

func notFound(ctx context.Context, writer http.ResponseWriter, lgr Logger) {

	rp := &respond.Respond{
		Writer: writer,
		Logger: lgr,
	}

	rp.NotFound(ctx)
}
