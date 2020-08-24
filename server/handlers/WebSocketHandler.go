package handlers

import (
	"air-sync/models"
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
	*models.Session
	*websocket.Conn
	request *http.Request
	logger  *log.Logger
}

func NewWebSocketHandler(repo *repos.SessionRepository) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  2048,
			WriteBufferSize: 2048,
		},
		repo: repo,
	}
}

func (h *WebSocketHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/ws/sessions/{id}", h.SetupWebSocket)
}

func (h *WebSocketHandler) SetupWebSocket(w http.ResponseWriter, req *http.Request) {
	req = util.DecorateRequest(req)
	conn, err := h.upgrader.Upgrade(w, req, nil)
	if err != nil {
		util.WriteHttpError(w, req, err)
		return
	}
	defer conn.Close()
	logger := util.RequestLogger(req)
	id := mux.Vars(req)["id"]
	session := h.repo.Get(id)
	ws := &WebSocketSession{
		Session: session,
		Conn:    conn,
		request: req,
	}
	if err := ws.Start(); err != nil {
		logger.Error(err)
	}
}

func (ws *WebSocketSession) Start() error {
	logger := ws.logger
	for item := range ws.Observe() {
		if err := item.E; err != nil {
			if err != util.ErrStreamClosed {
				return err
			} else {
				break
			}
		}
		if err := ws.WriteJSON(item.V); err != nil {
			logger.Errorf("Error writing JSON: %+s", err)
		}
	}
	return nil
}
