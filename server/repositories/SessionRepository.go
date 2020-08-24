package stores

import (
	"air-sync/models"

	log "github.com/sirupsen/logrus"
)

type SessionRepository struct {
	sessions map[string]*models.Session
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[string]*models.Session),
	}
}

func (r *SessionRepository) Create() *models.Session {
	session := models.NewSession()
	r.sessions[session.Id] = session
	log.Infof("Created new session ID %s", session.Id)
	return session
}

func (r *SessionRepository) All() []*models.Session {
	sessions := make([]*models.Session, len(r.sessions))
	idx := 0
	for _, session := range r.sessions {
		sessions[idx] = session
		idx++
	}
	return sessions
}

func (r *SessionRepository) Get(id string) *models.Session {
	return r.sessions[id]
}

func (r *SessionRepository) Delete(id string) bool {
	if _, ok := r.sessions[id]; ok {
		delete(r.sessions, id)
		log.Infof("Deleted session ID %s", id)
		return true
	}
	return false
}
