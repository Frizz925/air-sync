package handlers

import (
	"air-sync/models"
	repos "air-sync/repositories"
	"air-sync/util"
	"encoding/json"
	"io/ioutil"
	"net/http"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var ResSessionNotFound = &util.RestResponse{
	StatusCode: http.StatusNotFound,
	Error:      "Session not found",
}

type ApiHandler struct {
	repository *repos.SessionRepository
}

type SessionHandlerFunc func(req *http.Request, session *repos.StreamSession) (interface{}, error)

func NewApiHandler(repository *repos.SessionRepository) *ApiHandler {
	return &ApiHandler{repository: repository}
}

func (h *ApiHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/sessions", util.WrapRestHandlerFunc(h.GetSessions)).Methods("GET")
	r.HandleFunc("/api/sessions", util.WrapRestHandlerFunc(h.CreateSession)).Methods("POST")
	r.HandleFunc("/api/sessions/{id}", h.WrapSessionHandlerFunc(h.GetSession)).Methods("GET")
	r.HandleFunc("/api/sessions/{id}", util.WrapRestHandlerFunc(h.UpdateSession)).Methods("PUT")
	r.HandleFunc("/api/sessions/{id}", util.WrapRestHandlerFunc(h.DeleteSession)).Methods("DELETE")

	r.HandleFunc("/api/qr/generate", util.WrapHandlerFunc(h.GenerateQR)).Methods("POST")
}

func (h *ApiHandler) GetSessions(_ *http.Request) (*util.RestResponse, error) {
	return &util.RestResponse{
		Data: h.repository.All(),
	}, nil
}

func (h *ApiHandler) CreateSession(req *http.Request) (*util.RestResponse, error) {
	session := h.repository.Create()
	util.RequestLogger(req).WithField("session_id", session.Id).Info("Created new session")
	return &util.RestResponse{
		Message: "Session created",
		Data:    session.Id,
	}, nil
}

func (h *ApiHandler) GetSession(req *http.Request, session *repos.StreamSession) (interface{}, error) {
	return session.Session, nil
}

func (h *ApiHandler) UpdateSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	content := &models.Content{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(content); err != nil {
		return nil, err
	}
	if !h.repository.Update(id, content) {
		return ResSessionNotFound, nil
	}
	util.RequestLogger(req).WithFields(log.Fields{
		"session_id": id,
		"content":    content,
	}).Info("Updated session")
	return nil, nil
}

func (h *ApiHandler) DeleteSession(req *http.Request) (*util.RestResponse, error) {
	id := mux.Vars(req)["id"]
	if ok := h.repository.Delete(id); !ok {
		return ResSessionNotFound, nil
	}
	util.RequestLogger(req).WithField("session_id", id).Info("Deleted session")
	return &util.RestResponse{
		Message: "Session deleted",
	}, nil
}

func (h *ApiHandler) GenerateQR(req *http.Request) (*util.Response, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	q, err := qrcode.Encode(string(b), qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return &util.Response{
		ContentType: "image/png",
		Body:        q,
	}, nil
}

func (h *ApiHandler) WrapSessionHandlerFunc(handler SessionHandlerFunc) http.HandlerFunc {
	return util.WrapRestHandlerFunc(func(req *http.Request) (*util.RestResponse, error) {
		id := mux.Vars(req)["id"]
		session := h.repository.Get(id)
		if session == nil {
			return ResSessionNotFound, nil
		}
		data, err := handler(req, session)
		if err != nil {
			return nil, err
		}
		return &util.RestResponse{
			Data: data,
		}, nil
	})
}
