package handlers

import (
	repos "air-sync/repositories"
	"air-sync/util"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiHandler struct {
	*SessionRestHandler
	QrRestHandler
}

var _ RouteHandler = (*ApiHandler)(nil)

func NewApiHandler(repo *repos.SessionRepository) *ApiHandler {
	return &ApiHandler{
		SessionRestHandler: NewSessionRestHandler(repo),
	}
}

func (h *ApiHandler) RegisterRoutes(r *mux.Router) {
	s := r.PathPrefix("/api").Subrouter()
	h.SessionRestHandler.RegisterRoutes(s)
	h.QrRestHandler.RegisterRoutes(s)
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
