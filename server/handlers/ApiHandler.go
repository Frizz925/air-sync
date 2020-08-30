package handlers

import (
	"air-sync/util"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiHandler struct {
	handlers []RouteHandler
}

var _ RouteHandler = (*ApiHandler)(nil)

func NewApiHandler(handlers ...RouteHandler) *ApiHandler {
	return &ApiHandler{handlers}
}

func (h *ApiHandler) RegisterRoutes(r *mux.Router) {
	s := r.PathPrefix("/api").Subrouter()
	for _, rh := range h.handlers {
		rh.RegisterRoutes(s)
	}
	s.PathPrefix("").HandlerFunc(util.WrapRestHandlerFunc(h.NotFound))
}

func (h *ApiHandler) NotFound(req *http.Request) (*util.RestResponse, error) {
	return &util.RestResponse{
		StatusCode: 404,
		Status:     "error",
		Message:    "Resource not found",
		Error:      "Path does not exist",
	}, nil
}
