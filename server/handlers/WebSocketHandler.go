package handlers

import (
	repos "air-sync/repositories"
	"air-sync/util"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader *websocket.Upgrader
	repo     *repos.SessionRepository
}

type WebSocketSession struct {
	*repos.StreamSession
	conn    *websocket.Conn
	request *http.Request
	logger  *log.Logger
}

type OriginCheck func(req *http.Request) bool

var _ RouteHandler = (*WebSocketHandler)(nil)

func NewWebSocketHandler(repo *repos.SessionRepository, cors bool) *WebSocketHandler {
	var checkOrigin OriginCheck = nil
	if cors {
		checkOrigin = acceptAllOrigin
	}

	return &WebSocketHandler{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
			CheckOrigin:     checkOrigin,
		},
		repo: repo,
	}
}

func (h *WebSocketHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/ws/sessions/{id}", h.SetupWS)
}

func (h *WebSocketHandler) SetupWS(w http.ResponseWriter, req *http.Request) {
	req = util.DecorateRequest(req)
	id := mux.Vars(req)["id"]
	session := h.repo.Get(id)
	if session == nil {
		http.Error(w, ErrSessionNotFound.Error(), 404)
		return
	}
	logger := util.RequestLogger(req)
	conn, err := h.upgrader.Upgrade(w, req, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	defer conn.Close()
	ws := &WebSocketSession{
		StreamSession: session,
		conn:          conn,
		request:       req,
		logger:        logger,
	}
	if err := ws.Start(); err != nil {
		logger.Error(err)
	}
}

func (ws *WebSocketSession) Start() error {
	ws.logger.WithField("session_id", ws.Id).Info("New WebSocket client connected")
	defer ws.logger.WithField("session_id", ws.Id).Info("WebSocket client disconnected")

	sub := ws.Subscribe()
	defer sub.Unsubscribe()

	for item := range sub.Observe() {
		if err := item.E; err != nil {
			if err != util.ErrStreamClosed {
				return err
			} else {
				break
			}
		}
		if err := ws.conn.WriteJSON(item.V); err != nil {
			return err
		}
	}

	return nil
}

func acceptAllOrigin(_ *http.Request) bool {
	return true
}
