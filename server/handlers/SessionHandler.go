package handlers

import (
	"air-sync/models"
	"air-sync/models/events"
	repos "air-sync/repositories"
	"air-sync/util"
	"air-sync/util/logging"
	"air-sync/util/pubsub"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var (
	RestSessionNotFound = util.RestResponse{
		StatusCode: http.StatusNotFound,
		Message:    "Resource not found",
		Error:      "Session not found",
	}
	RestMessageNotFound = util.RestResponse{
		StatusCode: http.StatusNotFound,
		Message:    "Resource not found",
		Error:      "Message not found",
	}
	RestAttachmentNotFound = util.RestResponse{
		StatusCode: http.StatusNotFound,
		Message:    "Resource not found",
		Error:      "Attachment not found",
	}
)

type SessionHandler struct {
	repo  repos.SessionRepository
	pub   *pubsub.Publisher
	topic *pubsub.Topic
}

type SessionHandlerFunc func(req *http.Request, session models.Session) (interface{}, error)

func NewSessionHandler(repo repos.SessionRepository, pub *pubsub.Publisher) *SessionHandler {
	return &SessionHandler{
		repo:  repo,
		pub:   pub,
		topic: pub.Topic(events.EventSession),
	}
}

func (h *SessionHandler) CreateSessionLogger(req *http.Request, session models.Session) *log.Logger {
	logger := util.CreateRequestLogger(req)
	h.ApplySessionLogger(logger, session)
	return logger
}

func (h *SessionHandler) ApplySessionLogger(logger *log.Logger, session models.Session) {
	logger.Formatter = logging.NewSessionLogFormatter(logger.Formatter, session)
}

func (h *SessionHandler) WrapSessionHandlerFunc(handler SessionHandlerFunc) http.HandlerFunc {
	return util.WrapRestHandlerFunc(func(req *http.Request) (*util.RestResponse, error) {
		id := mux.Vars(req)["id"]
		session, err := h.repo.Find(id)
		if err != nil {
			return h.HandleSessionRestError(err)
		}
		h.ApplySessionLogger(util.RequestLogger(req), session)
		data, err := handler(req, session)
		if err != nil {
			return nil, err
		}
		return &util.RestResponse{
			Data: data,
		}, nil
	})
}

func (h *SessionHandler) HandleSessionRestError(err error) (*util.RestResponse, error) {
	switch err {
	case repos.ErrSessionNotFound:
		return &RestSessionNotFound, nil
	case repos.ErrMessageNotFound:
		return &RestMessageNotFound, nil
	case repos.ErrAttachmentNotFound:
		return &RestAttachmentNotFound, nil
	}
	return nil, err
}

func (h *SessionHandler) HandleSessionError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	switch err {
	case repos.ErrSessionNotFound:
		code = http.StatusNotFound
	case repos.ErrMessageNotFound:
		code = http.StatusNotFound
	case repos.ErrAttachmentNotFound:
		code = http.StatusNotFound
	}
	http.Error(w, err.Error(), code)
}
