package handlers

import (
	repos "air-sync/repositories"
	"air-sync/util"
	"encoding/json"
	"net/http"
	"strings"

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

func NewStreamingHandler(repo *repos.SessionRepository) *StreamingHandler {
	return &StreamingHandler{
		SessionHandler: NewSessionHandler(repo),
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
	session := h.repo.Get(id)
	if session == nil {
		http.Error(w, ErrSessionNotFound.Error(), http.StatusNotFound)
		return
	}
	logger := h.CreateSessionLogger(req, session.Session)
	logger.Info("Started event streaming session")
	defer logger.Info("Event streaming session ended")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(200)

	if err := h.SendEvent(rwf, "ping", ""); err != nil {
		logger.Error(err)
		return
	}

	sub := session.Subscribe()
	defer sub.Unsubscribe()

	for item := range sub.Observe() {
		if err := item.E; err != nil {
			if err != util.ErrStreamClosed {
				logger.Error(err)
			}
			break
		}
		b, err := json.Marshal(item.V)
		if err != nil {
			logger.Error(err)
			break
		}
		if err := h.SendEvent(rwf, "content", string(b)); err != nil {
			logger.Error(err)
			break
		}
	}
}

func (h *StreamingHandler) SendEvent(w ResponseWriteFlusher, event string, data string) error {
	payload := strings.Join([]string{
		"id: " + uuid.NewV4().String(),
		"event: " + event,
		"data: " + data,
		"\n",
	}, "\n")
	if _, err := w.Write([]byte(payload)); err != nil {
		return err
	}
	w.Flush()
	return nil
}
