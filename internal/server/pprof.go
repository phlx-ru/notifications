package server

import (
	"net/http"
	"net/http/pprof"
	runtimePProf "runtime/pprof"

	kratosHTTP "github.com/go-kratos/kratos/v2/transport/http"
)

func registerPProfHandlers(s *kratosHTTP.Server) {
	r := s.Route("/")

	wrap := func(handle func(w http.ResponseWriter, r *http.Request)) kratosHTTP.HandlerFunc {
		return func(context kratosHTTP.Context) error {
			handle(context.Response(), context.Request())
			return nil
		}
	}

	r.GET("/debug/pprof/index", wrap(pprof.Index)) // `/index` 'cause of trailing slash redirection
	r.GET("/debug/pprof/cmdline", wrap(pprof.Cmdline))
	r.GET("/debug/pprof/profile", wrap(pprof.Profile))
	r.GET("/debug/pprof/symbol", wrap(pprof.Symbol))
	r.GET("/debug/pprof/trace", wrap(pprof.Trace))

	// goroutine, threadcreate, heap, allocs, block, mutex
	for _, p := range runtimePProf.Profiles() {
		r.GET("/debug/pprof/"+p.Name(), wrap(pprof.Handler(p.Name()).ServeHTTP))
	}
}
