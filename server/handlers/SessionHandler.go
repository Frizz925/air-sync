package handlers

import (
	"air-sync/models"
	repos "air-sync/repositories"
	"air-sync/util"
	"air-sync/util/logging"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var (
	ErrSessionNotFound = errors.New("Session not found")
	ResSessionNotFound = &util.RestResponse{
		StatusCode: http.StatusNotFound,
		Message:    "Resource not found",
		Error:      "Session not found",
	}
)

type SessionHandler struct {
	repo *repos.SessionRepository
}

type SessionHandlerFunc func(req *http.Request, session *repos.StreamSession) (interface{}, error)

func NewSessionHandler(repo *repos.SessionRepository) *SessionHandler {
	return &SessionHandler{
		repo: repo,
	}
}

func (h *SessionHandler) CreateSessionLogger(req *http.Request, session *models.Session) *log.Logger {
	logger := util.CreateRequestLogger(req)
	h.ApplySessionLogger(logger, session)
	return logger
}

func (h *SessionHandler) ApplySessionLogger(logger *log.Logger, session *models.Session) {
	logger.Formatter = logging.NewSessionLogFormatter(logger.Formatter, session)
}

func (h *SessionHandler) WrapSessionHandlerFunc(handler SessionHandlerFunc) http.HandlerFunc {
	return util.WrapRestHandlerFunc(func(req *http.Request) (*util.RestResponse, error) {
		id := mux.Vars(req)["id"]
		session := h.repo.Get(id)
		if session == nil {
			return ResSessionNotFound, nil
		}
		h.ApplySessionLogger(util.RequestLogger(req), session.Session)
		data, err := handler(req, session)
		if err != nil {
			return nil, err
		}
		return &util.RestResponse{
			Data: data,
		}, nil
	})
}
