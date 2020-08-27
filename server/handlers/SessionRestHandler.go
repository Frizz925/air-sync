package handlers

import (
	"air-sync/models"
	repos "air-sync/repositories"
	"air-sync/util"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type SessionRestHandler struct {
	*SessionHandler
}

var _ RouteHandler = (*SessionRestHandler)(nil)

func NewSessionRestHandler(repo *repos.SessionRepository) *SessionRestHandler {
	return &SessionRestHandler{
		SessionHandler: NewSessionHandler(repo),
	}
}

func (h *SessionRestHandler) RegisterRoutes(r *mux.Router) {
	s := r.PathPrefix("/sessions").Subrouter()
	// s.HandleFunc("/", util.WrapRestHandlerFunc(h.GetSessions)).Methods("GET")
	s.HandleFunc("", util.WrapRestHandlerFunc(h.CreateSession)).Methods("POST")
	s.HandleFunc("/{id}", h.WrapSessionHandlerFunc(h.GetSession)).Methods("GET")
	s.HandleFunc("/{id}", util.WrapRestHandlerFunc(h.UpdateSession)).Methods("PUT")
	s.HandleFunc("/{id}", util.WrapRestHandlerFunc(h.DeleteSession)).Methods("DELETE")
}

func (h *SessionRestHandler) GetSessions(_ *http.Request) (*util.RestResponse, error) {
	return &util.RestResponse{
		Data: h.repo.All(),
	}, nil
}

func (h *SessionRestHandler) CreateSession(req *http.Request) (*util.RestResponse, error) {
	session := h.repo.Create()
	util.RequestLogger(req).WithField("session_id", session.Id).Info("Created new session")
	return &util.RestResponse{
		Message: "Session created",
		Data:    session.Id,
	}, nil
}

func (h *SessionRestHandler) GetSession(req *http.Request, session *repos.StreamSession) (interface{}, error) {
	return session.Session, nil
}

func (h *SessionRestHandler) UpdateSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	content := &models.Content{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(content); err != nil {
		return nil, err
	}
	if !h.repo.Update(id, content) {
		return ResSessionNotFound, nil
	}
	util.RequestLogger(req).WithFields(log.Fields{
		"session_id": id,
		"content":    content,
	}).Info("Updated session")
	return nil, nil
}

func (h *SessionRestHandler) DeleteSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	if ok := h.repo.Delete(id); !ok {
		return ResSessionNotFound, nil
	}
	util.RequestLogger(req).WithField("session_id", id).Info("Deleted session")
	return &util.RestResponse{
		Message: "Session deleted",
	}, nil
}
