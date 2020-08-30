package handlers

import (
	repos "air-sync/repositories"
	"air-sync/subscribers/events"
	"air-sync/util/pubsub"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type ResponseWriteFlusher interface {
	http.ResponseWriter
	http.Flusher
}

type StreamingHandler struct {
	*SessionHandler
}

var _ RouteHandler = (*StreamingHandler)(nil)

func NewStreamingHandler(repo repos.SessionRepository, stream *pubsub.Stream) *StreamingHandler {
	return &StreamingHandler{
		SessionHandler: NewSessionHandler(repo, stream),
	}
}

func (h *StreamingHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/sse/sessions/{id}", h.SendSessionEvent).Methods("GET")
}

func (h *StreamingHandler) SendSessionEvent(w http.ResponseWriter, req *http.Request) {
	rwf, ok := w.(ResponseWriteFlusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	id := mux.Vars(req)["id"]
	session, err := h.repo.Get(id)
	if err != nil {
		h.HandleSessionError(w, err)
		return
	}

	logger := h.CreateSessionLogger(req, session)
	logger.Info("Started event streaming session")
	defer logger.Info("Event streaming session ended")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(200)

	if err := h.SendEvent(rwf, "heartbeat", ""); err != nil {
		logger.Error(err)
		return
	}

	topic := h.stream.Topic(events.SessionEventName(id))
	if err := h.HandleStream(rwf, req, topic); err != nil {
		if err != io.EOF {
			logger.Error(err)
		}
	} else {
		if err := h.SendEvent(rwf, "close", ""); err != nil {
			logger.Error(err)
		}
	}
}

func (h *StreamingHandler) HandleStream(rwf ResponseWriteFlusher, req *http.Request, topic *pubsub.Topic) error {
	sub := topic.Subscribe()
	defer sub.Unsubscribe()

	ctx := req.Context()
	ch := sub.Observe()
	for {
		timeout := time.After(30 * time.Second)
		select {
		case item := <-ch:
			if item.E != nil {
				if item.E != pubsub.ErrStreamClosed {
					return item.E
				} else {
					return nil
				}
			}
			b, err := json.Marshal(item.V)
			if err != nil {
				return err
			}
			if err := h.SendEvent(rwf, "message", string(b)); err != nil {
				return err
			}
		case <-timeout:
			if err := h.SendEvent(rwf, "heartbeat", ""); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (h *StreamingHandler) SendEvent(rwf ResponseWriteFlusher, event string, data string) error {
	payload := strings.Join([]string{
		"id: " + uuid.NewV4().String(),
		"event: " + event,
		"data: " + data,
		"\n",
	}, "\n")
	if _, err := rwf.Write([]byte(payload)); err != nil {
		return err
	}
	rwf.Flush()
	return nil
}
