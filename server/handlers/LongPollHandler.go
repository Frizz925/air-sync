package handlers

import (
	"air-sync/models/formatters"
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
	session, err := h.repo.Find(id)
	if err != nil {
		return h.HandleSessionRestError(err)
	}
	logger := util.RequestLogger(req)
	h.ApplySessionLogger(logger, session)

	logger.Info("Started long-polling session")
	defer logger.Info("Long-polling session ended")

	ctx := req.Context()
	sub := h.stream.Topic(events.SessionEventName(id)).Subscribe()
	defer sub.Unsubscribe()

	timeout := time.After(30 * time.Second) // Poll for 30 seconds
	select {
	case item := <-sub.Observe():
		if err := item.E; err != nil {
			if err != pubsub.ErrStreamClosed {
				return nil, err
			}
		}
		if v, ok := item.V.(events.SessionEvent); ok {
			return &util.RestResponse{
				Message: "New session event",
				Data:    formatters.FromSessionEvent(v),
			}, nil
		}
	case <-timeout:
	case <-ctx.Done():
	}

	return &util.RestResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}
