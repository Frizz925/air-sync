package handlers

import (
	"air-sync/models"
	"air-sync/models/events"
	repos "air-sync/repositories"
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

func NewSessionRestHandler(repo repos.SessionRepository, pub *pubsub.Publisher) *SessionRestHandler {
	return &SessionRestHandler{
		SessionHandler: NewSessionHandler(repo, pub),
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
	h.topic.Publish(events.SessionEvent{
		Event:     events.EventSessionCreated,
		SessionId: session.Id,
		Value:     events.SessionCreate(session),
	})
	util.RequestLogger(req).WithField("session_id", session.Id).Info("Created new session")
	return &util.RestResponse{
		Message: "Session created",
		Data:    session.Id,
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
	h.topic.Publish(events.SessionEvent{
		Event:     events.EventSessionDeleted,
		SessionId: id,
		Value:     events.SessionDelete(id),
	})
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
	if insert.Body == "" && insert.AttachmentId == "" {
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
	h.topic.Publish(events.SessionEvent{
		Event:     events.EventMessageInserted,
		SessionId: id,
		Value: events.MessageInsert{
			SessionId: id,
			Message:   message,
		},
	})
	util.RequestLogger(req).WithFields(log.Fields{
		"session_id": id,
		"message_id": message.Id,
	}).Info("Inserted message")
	return &util.RestResponse{
		Message: "Message inserted",
		Data:    message.Id,
	}, nil
}

func (h *SessionRestHandler) DeleteMessage(req *http.Request) (*util.RestResponse, error) {
	vars := mux.Vars(req)
	sessionId := vars["id"]
	messageId := vars["messageId"]
	if err := h.repo.DeleteMessage(sessionId, messageId); err != nil {
		return h.HandleSessionRestError(err)
	}
	h.pub.Topic(events.EventSession).Publish(events.SessionEvent{
		Event:     events.EventMessageDeleted,
		SessionId: sessionId,
		Value: events.MessageDelete{
			SessionId: sessionId,
			MessageId: messageId,
		},
	})
	util.RequestLogger(req).WithFields(log.Fields{
		"session_id": sessionId,
		"message_id": messageId,
	}).Info("Deleted message")
	return &util.RestResponse{
		Message: "Message deleted",
	}, nil
}
