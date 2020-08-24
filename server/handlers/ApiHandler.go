package handlers

import (
	repos "air-sync/repositories"
	"air-sync/util"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiHandler struct {
	repository *repos.SessionRepository
}

func NewApiHandler(repository *repos.SessionRepository) *ApiHandler {
	return &ApiHandler{repository: repository}
}

func (h *ApiHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/sessions", util.WrapRestHandlerFunc(h.GetSessions)).Methods("GET")
	r.HandleFunc("/api/sessions/create", util.WrapRestHandlerFunc(h.CreateSession)).Methods("POST")
	r.HandleFunc("/api/sessions/{id}", util.WrapRestHandlerFunc(h.DeleteSession)).Methods("DELETE")
}

func (h *ApiHandler) GetSessions(req *http.Request) (*util.RestResponse, error) {
	sessions := h.repository.All()
	ids := make([]string, len(sessions))
	for idx, session := range sessions {
		ids[idx] = session.Id
	}
	return &util.RestResponse{
		Data: ids,
	}, nil
}

func (h *ApiHandler) CreateSession(req *http.Request) (*util.RestResponse, error) {
	session := h.repository.Create()
	return &util.RestResponse{
		Message: "Session created",
		Data:    session.Id,
	}, nil
}

func (h *ApiHandler) DeleteSession(req *http.Request) (*util.RestResponse, error) {
	vars := mux.Vars(req)
	if ok := h.repository.Delete(vars["id"]); !ok {
		return &util.RestResponse{
			StatusCode: 404,
			Error:      "Session not found",
		}, nil
	}
	return &util.RestResponse{
		Message: "Session deleted",
	}, nil
}
