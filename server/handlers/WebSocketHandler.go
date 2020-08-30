package handlers

import (
	"air-sync/models"
	repos "air-sync/repositories"
	"air-sync/subscribers/events"
	"air-sync/util"
	"air-sync/util/pubsub"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	*SessionHandler
	upgrader *websocket.Upgrader
}

type WebSocketHandlerOptions struct {
	Repository repos.SessionRepository
	Stream     *pubsub.Stream
	EnableCORS bool
}

type WebSocketSession struct {
	*models.Session
	*pubsub.Topic
	conn    *websocket.Conn
	request *http.Request
	logger  *log.Logger
}

type OriginCheck func(req *http.Request) bool

var _ RouteHandler = (*WebSocketHandler)(nil)

func NewWebSocketHandler(opts WebSocketHandlerOptions) *WebSocketHandler {
	var checkOrigin OriginCheck = nil
	if opts.EnableCORS {
		checkOrigin = acceptAllOrigin
	}

	return &WebSocketHandler{
		SessionHandler: NewSessionHandler(opts.Repository, opts.Stream),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
			CheckOrigin:     checkOrigin,
		},
	}
}

func (h *WebSocketHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/ws/sessions/{id}", h.SetupWS)
}

func (h *WebSocketHandler) SetupWS(w http.ResponseWriter, req *http.Request) {
	req = util.DecorateRequest(req)
	id := mux.Vars(req)["id"]
	session, err := h.repo.Get(id)
	if err != nil {
		h.HandleSessionError(w, err)
		return
	}
	logger := util.RequestLogger(req)
	conn, err := h.upgrader.Upgrade(w, req, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	defer conn.Close()
	name := events.SessionEventName(id)
	ws := &WebSocketSession{
		Session: session,
		Topic:   h.stream.Topic(name),
		conn:    conn,
		request: req,
		logger:  logger,
	}
	if err := ws.Start(); err != nil {
		logger.Error(err)
	}
}

func (ws *WebSocketSession) Start() error {
	ws.logger.WithField("session_id", ws.Id).Info("New WebSocket client connected")
	defer ws.logger.WithField("session_id", ws.Id).Info("WebSocket client disconnected")
	return ws.ForEach(func(v interface{}) error {
		return ws.conn.WriteJSON(v)
	})
}

func acceptAllOrigin(_ *http.Request) bool {
	return true
}
