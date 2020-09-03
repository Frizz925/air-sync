package handlers

import (
	"air-sync/models"
	"air-sync/models/events"
	"air-sync/models/formatters"
	repos "air-sync/repositories"
	"air-sync/util"
	"air-sync/util/pubsub"
	"context"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketOptions struct {
	Repository repos.SessionRepository
	Publisher  *pubsub.Publisher
	EnableCORS bool
}

type WebSocketHandler struct {
	*SessionHandler
	upgrader *websocket.Upgrader
}

type WebSocketSession struct {
	models.Session
	*pubsub.Subscriber
	conn    *websocket.Conn
	request *http.Request
	logger  *log.Logger
}

type OriginCheck func(req *http.Request) bool

var _ RouteHandler = (*WebSocketHandler)(nil)

func NewWebSocketHandler(opts WebSocketOptions) *WebSocketHandler {
	var checkOrigin OriginCheck = nil
	if opts.EnableCORS {
		checkOrigin = acceptAllOrigin
	}
	return &WebSocketHandler{
		SessionHandler: NewSessionHandler(opts.Repository, opts.Publisher),
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
	session, err := h.repo.Find(id)
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

	sub := h.pub.Topic(events.EventSessionID(id)).Subscribe()
	defer sub.Unsubscribe()

	ws := &WebSocketSession{
		Session:    session,
		Subscriber: sub,
		conn:       conn,
		request:    req,
		logger:     logger,
	}

	ws.Setup()
	if err := ws.Start(); err != nil {
		if err != io.EOF {
			logger.Error(err)
		}
	}
}

func (ws *WebSocketSession) Setup() {
	ws.conn.SetPingHandler(nil)
	ws.conn.SetPongHandler(nil)
	ws.conn.SetCloseHandler(nil)
}

func (ws *WebSocketSession) Start() error {
	ws.logger.WithField("session_id", ws.ID).Info("New WebSocket client connected")
	defer ws.logger.WithField("session_id", ws.ID).Info("WebSocket client disconnected")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := ws.HandleReads(); err != nil {
			if _, ok := err.(*websocket.CloseError); !ok {
				ws.logger.Error(err)
			}
		}
		cancel()
	}()

	ch := ws.Channel()
	for {
		timeout := time.After(30 * time.Second)
		select {
		case v := <-ch:
			event, ok := v.(events.SessionEvent)
			if !ok {
				continue
			}
			err := ws.conn.WriteJSON(formatters.FromSessionEvent(event))
			if err != nil {
				return err
			}
		case <-timeout:
			err := ws.conn.WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (ws *WebSocketSession) HandleReads() error {
	for {
		messageType, _, err := ws.conn.ReadMessage()
		if err != nil {
			return err
		}
		switch messageType {
		case websocket.PingMessage:
			err := ws.conn.WriteMessage(websocket.PongMessage, []byte(""))
			if err != nil {
				return err
			}
		case websocket.CloseMessage:
			return nil
		}
	}
}

func acceptAllOrigin(_ *http.Request) bool {
	return true
}
