package handlers

import (
	"air-sync/models/events"
	"air-sync/models/formatters"
	repos "air-sync/repositories"
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

func NewLongPollHandler(repo repos.SessionRepository, pub *pubsub.Publisher) *LongPollHandler {
	return &LongPollHandler{
		SessionHandler: NewSessionHandler(repo, pub),
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
	sub := h.pub.Topic(events.EventSessionId(id)).Subscribe()
	defer sub.Unsubscribe()
	ch := sub.Channel()

	timeout := time.After(30 * time.Second) // Poll for 30 seconds
	select {
	case v := <-ch:
		if event, ok := v.(events.SessionEvent); ok {
			return &util.RestResponse{
				Message: "New session event",
				Data:    formatters.FromSessionEvent(event),
			}, nil
		}
	case <-timeout:
	case <-ctx.Done():
	}

	return &util.RestResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}
