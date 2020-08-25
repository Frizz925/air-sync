package stores

import (
	"air-sync/models"
	"air-sync/util"
	"sync"
)

type StreamSession struct {
	*models.Session
	*util.Stream
}

type SessionRepository struct {
	sync.RWMutex
	sessions map[string]*StreamSession
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions: make(map[string]*StreamSession),
	}
}

func (r *SessionRepository) Create() *StreamSession {
	defer r.Unlock()
	r.Lock()
	session := &StreamSession{
		Session: models.NewSession(),
		Stream:  util.NewStream(),
	}
	r.sessions[session.Id] = session
	return session
}

func (r *SessionRepository) All() []string {
	defer r.RUnlock()
	r.RLock()
	sessionIds := make([]string, len(r.sessions))
	idx := 0
	for _, session := range r.sessions {
		sessionIds[idx] = session.Id
		idx++
	}
	return sessionIds
}

func (r *SessionRepository) Get(id string) *StreamSession {
	defer r.RUnlock()
	r.RLock()
	return r.sessions[id]
}

func (r *SessionRepository) Update(id string, content *models.Content) bool {
	defer r.Unlock()
	r.Lock()
	if s, ok := r.sessions[id]; ok {
		s.Content = content
		s.Fire(content)
		return true
	}
	return false
}

func (r *SessionRepository) Delete(id string) bool {
	defer r.Unlock()
	r.Lock()
	if session, ok := r.sessions[id]; ok {
		session.Shutdown()
		delete(r.sessions, id)
		return true
	}
	return false
}
