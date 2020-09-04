package app

import (
	"air-sync/util"
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type WebApplication struct {
	Addr             string
	Router           *mux.Router
	EnableCORS       bool
	CloudEnvironment string
}

var _ Application = (*WebApplication)(nil)

func (s *WebApplication) Start(ctx context.Context) error {
	if s.CloudEnvironment == string(util.CloudEnvGoogleCloud) {
		s.Router.HandleFunc("/_ah/warmup", util.WrapHandlerFunc(s.handleWarmup))
	}

	err := s.Router.Walk(func(r *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := r.GetPathTemplate()
		if err != nil {
			return err
		}
		log.Infof("Route registered: %s", t)
		return nil
	})
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp4", s.Addr)
	if err != nil {
		return err
	}

	var handler http.Handler
	if s.EnableCORS {
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
		})
		handler = c.Handler(s.Router)
		log.Info("CORS enabled")
	} else {
		handler = s.Router
	}

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  45 * time.Second,
		WriteTimeout: 45 * time.Second,
	}

	go func() {
		if err := server.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Infof("Server listening at %s", s.Addr)
	<-ctx.Done()
	log.Info("Shutting down server")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		return err
	}

	log.Info("Server shutdown properly")
	return nil
}

func (s *WebApplication) handleWarmup(req *http.Request) (*util.Response, error) {
	util.RequestLogger(req).Info("Warmup request received")
	return util.SuccessResponse, nil
}
