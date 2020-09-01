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
	s.HandleFunc("/{id}", util.WrapRestHandlerFunc(h.DeleteSession)).Methods("DELETE")
	s.HandleFunc("/{id}", util.WrapRestHandlerFunc(h.InsertMessage)).Methods("PUT")
	s.HandleFunc("/{id}/{messageId}", util.WrapRestHandlerFunc(h.DeleteMessage)).Methods("DELETE")
}

func (h *SessionRestHandler) CreateSession(req *http.Request) (*util.RestResponse, error) {
	session, err := h.repo.Create()
	if err != nil {
		return h.HandleSessionRestError(err)
	}
	util.RequestLogger(req).WithField("session_id", session.ID).Info("Created new session")
	return &util.RestResponse{
		Message: "Session created",
		Data:    session.ID,
	}, nil
}

func (h *SessionRestHandler) GetSession(req *http.Request, session models.Session) (interface{}, error) {
	return session, nil
}

func (h *SessionRestHandler) DeleteSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	if err := h.repo.Delete(id); err != nil {
		return h.HandleSessionRestError(err)
	}
	h.stream.Topic(events.SessionDeleted).Fire(events.SessionDelete(id))
	util.RequestLogger(req).WithField("session_id", id).Info("Deleted session")
	return &util.RestResponse{
		Message: "Session deleted",
	}, nil
}

func (h *SessionRestHandler) InsertMessage(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	insert := models.InsertMessage{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&insert); err != nil {
		return nil, err
	}
	if insert.Body == "" && insert.AttachmentID == "" {
		return &util.RestResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Malformed request",
			Error:      "Message body and attachment are empty",
		}, nil
	}
	message, err := h.repo.InsertMessage(id, insert)
	if err != nil {
		return h.HandleSessionRestError(err)
	}
	h.stream.Topic(events.MessageInserted).Fire(events.MessageInsert{
		SessionId: id,
		Message:   message,
	})
	util.RequestLogger(req).WithFields(log.Fields{
		"session_id": id,
		"message_id": message.ID,
	}).Info("Inserted message")
	return &util.RestResponse{
		Message: "Message inserted",
		Data:    message.ID,
	}, nil
}

func (h *SessionRestHandler) DeleteMessage(req *http.Request) (*util.RestResponse, error) {
	vars := mux.Vars(req)
	sessionId := vars["id"]
	messageId := vars["messageId"]
	if err := h.repo.DeleteMessage(sessionId, messageId); err != nil {
		return h.HandleSessionRestError(err)
	}
	h.stream.Topic(events.MessageDeleted).Fire(events.MessageDelete{
		SessionId: sessionId,
		MessageId: messageId,
	})
	util.RequestLogger(req).WithFields(log.Fields{
		"session_id": sessionId,
		"message_id": messageId,
	}).Info("Deleted message")
	return &util.RestResponse{
		Message: "Message deleted",
	}, nil
}
