package handlers

import (
	"air-sync/services"
	"air-sync/util"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type CronEnvironment string

const (
	CronEnvLocal     CronEnvironment = "local"
	CronEnvAppEngine CronEnvironment = "app_engine"
)

type CronHandler struct {
	env  CronEnvironment
	cron *services.CronJobService
}

var _ RouteHandler = (*CronHandler)(nil)

func NewCronHandler(env CronEnvironment, cron *services.CronJobService) *CronHandler {
	return &CronHandler{
		env:  env,
		cron: cron,
	}
}

func (h *CronHandler) RegisterRoutes(r *mux.Router) {
	s := r.PathPrefix("/cron").Subrouter()
	s.Use(h.Middleware)
	s.HandleFunc("/cleanup", util.WrapHandlerFunc(h.CleanupJob))
}

func (h *CronHandler) CleanupJob(req *http.Request) (*util.Response, error) {
	if err := h.cron.RunCleanupJob(); err != nil {
		return nil, err
	}
	return util.SuccessResponse, nil
}

func (h *CronHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.WithFields(log.Fields{
			"cron_env":  h.env,
			"client_ip": util.GetClientIP(req),
		}).Info("Receiving cron job request")
		if !h.ValidateRequest(req) {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func (h *CronHandler) ValidateRequest(req *http.Request) bool {
	switch h.env {
	case CronEnvLocal:
		if req.Header.Get("X-Cron-Agent") != "Postman" {
			return false
		}
		return util.GetClientIP(req) == "127.0.0.1"
	case CronEnvAppEngine:
		if req.Header.Get("X-Appengine-Cron") != "true" {
			return false
		}
		return util.GetClientIP(req) == "10.0.0.1"
	}
	return false
}
