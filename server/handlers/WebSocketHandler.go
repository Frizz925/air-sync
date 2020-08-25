package handlers

import (
	repos "air-sync/repositories"
	"air-sync/util"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var ErrSessionNotFound = errors.New("Session not found")

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

func NewWebSocketHandler(repo *repos.SessionRepository) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
		repo: repo,
	}
}

func (h *WebSocketHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/ws/sessions/{id}", h.SetupWebSocket)
}

func (h *WebSocketHandler) SetupWebSocket(w http.ResponseWriter, req *http.Request) {
	req = util.DecorateRequest(req)
	logger := util.RequestLogger(req)
	conn, err := h.upgrader.Upgrade(w, req, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	defer conn.Close()
	id := mux.Vars(req)["id"]
	session := h.repo.Get(id)
	if session == nil {
		util.WriteHttpError(w, req, ErrSessionNotFound)
		return
	}
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
	for item := range ws.Observe() {
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
