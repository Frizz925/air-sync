package handlers

import (
	repos "air-sync/repositories"
	"air-sync/subscribers/events"
	"air-sync/util"
	"air-sync/util/pubsub"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type LongPollHandler struct {
	*SessionHandler
}

var _ RouteHandler = (*LongPollHandler)(nil)

func NewLongPollHandler(repo repos.SessionRepository, stream *pubsub.Stream) *LongPollHandler {
	return &LongPollHandler{
		SessionHandler: NewSessionHandler(repo, stream),
	}
}

func (h *LongPollHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/lp/sessions/{id}", util.WrapRestHandlerFunc(h.PollSession)).Methods("GET")
}

func (h *LongPollHandler) PollSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	session, err := h.repo.Get(id)
	if err != nil {
		return h.HandleSessionRestError(err)
	}
	logger := util.RequestLogger(req)
	h.ApplySessionLogger(logger, session)

	logger.Info("Started long-polling session")
	defer logger.Info("Long-polling session ended")

	sub := h.stream.Topic(events.SessionEventName(id)).Subscribe()
	defer sub.Unsubscribe()

	timeout := time.After(30 * time.Second) // Poll for 30 seconds
	select {
	case item := <-sub.Observe():
		if err := item.E; err != nil {
			if err != pubsub.ErrStreamClosed {
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
