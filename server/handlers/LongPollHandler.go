package handlers

import (
	repos "air-sync/repositories"
	"air-sync/util"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type LongPollHandler struct {
	*SessionHandler
}

var _ RouteHandler = (*LongPollHandler)(nil)

func NewLongPollHandler(repo *repos.SessionRepository) *LongPollHandler {
	return &LongPollHandler{
		SessionHandler: NewSessionHandler(repo),
	}
}

func (h *LongPollHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/lp/sessions/{id}", util.WrapRestHandlerFunc(h.PollSession)).Methods("GET")
}

func (h *LongPollHandler) PollSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	session := h.repo.Get(id)
	if session == nil {
		return ResSessionNotFound, nil
	}
	logger := util.RequestLogger(req)
	h.ApplySessionLogger(logger, session.Session)

	logger.Info("Started long-polling session")
	defer logger.Info("Long-polling session ended")

	sub := session.Subscribe()
	timeout := time.After(30 * time.Second) // Poll for 30 seconds
	defer sub.Unsubscribe()

	select {
	case item := <-sub.Observe():
		if err := item.E; err != nil {
			if err != util.ErrStreamClosed {
				return nil, err
			}
		} else {
			return &util.RestResponse{
				Data: item.V,
			}, nil
		}
	case <-timeout:
	}

	return &util.RestResponse{
		StatusCode: 204,
		Status:     "success",
		Message:    "No message received in time",
	}, nil
}
