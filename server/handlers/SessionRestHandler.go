package handlers

import (
	"air-sync/models"
	repos "air-sync/repositories"
	"air-sync/subscribers/events"
	"air-sync/util"
	"air-sync/util/pubsub"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type SessionRestHandler struct {
	*SessionHandler
}

var _ RouteHandler = (*SessionRestHandler)(nil)

func NewSessionRestHandler(repo repos.SessionRepository, stream *pubsub.Stream) *SessionRestHandler {
	return &SessionRestHandler{
		SessionHandler: NewSessionHandler(repo, stream),
	}
}

func (h *SessionRestHandler) RegisterRoutes(r *mux.Router) {
	s := r.PathPrefix("/sessions").Subrouter()
	s.HandleFunc("", util.WrapRestHandlerFunc(h.CreateSession)).Methods("POST")
	s.HandleFunc("/{id}", h.WrapSessionHandlerFunc(h.GetSession)).Methods("GET")
	s.HandleFunc("/{id}", util.WrapRestHandlerFunc(h.UpdateSession)).Methods("PUT")
	s.HandleFunc("/{id}", util.WrapRestHandlerFunc(h.DeleteSession)).Methods("DELETE")
}

func (h *SessionRestHandler) CreateSession(req *http.Request) (*util.RestResponse, error) {
	session, err := h.repo.Create()
	if err != nil {
		return h.HandleSessionRestError(err)
	}
	util.RequestLogger(req).WithField("session_id", session.Id).Info("Created new session")
	return &util.RestResponse{
		Message: "Session created",
		Data:    session.Id,
	}, nil
}

func (h *SessionRestHandler) GetSession(req *http.Request, session *models.Session) (interface{}, error) {
	return session, nil
}

func (h *SessionRestHandler) UpdateSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	message := models.NewMessage()
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&message); err != nil {
		return nil, err
	}
	if message.Content == "" {
		return &util.RestResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Malformed request",
			Error:      "Message content is empty",
		}, nil
	}
	h.stream.Topic("update").Fire(events.SessionUpdate{
		Id:      id,
		Message: message,
	})
	util.RequestLogger(req).WithFields(log.Fields{
		"session_id": id,
		"message":    message,
	}).Info("Updated session")
	return &util.RestResponse{
		Message: "Session updated",
	}, nil
}

func (h *SessionRestHandler) DeleteSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	h.stream.Topic("delete").Fire(events.SessionDelete(id))
	util.RequestLogger(req).WithField("session_id", id).Info("Deleted session")
	return &util.RestResponse{
		Message: "Session deleted",
	}, nil
}
