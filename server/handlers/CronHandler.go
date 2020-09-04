package handlers

import (
	"air-sync/services"
	"air-sync/util"
	"net"
	"net/http"

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
		return h.GetClientIP(req) == "127.0.0.1"
	case CronEnvAppEngine:
		if req.Header.Get("X-Appengine-Cron") != "true" {
			return false
		}
		return h.GetClientIP(req) == "10.0.0.1"
	}
	return false
}

func (h *CronHandler) GetClientIP(req *http.Request) string {
	if host := req.Header.Get("X-Real-IP"); host != "" {
		return host
	}
	if host := req.Header.Get("X-Forwarded-For"); host != "" {
		return host
	}
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return ""
	}
	return host
}
