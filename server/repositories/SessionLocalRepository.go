package stores

import (
	"air-sync/models"
	"sync"
)

type SessionLocalRepository struct {
	sync.RWMutex
	sessions map[string]models.Session
}

var _ SessionRepository = (*SessionLocalRepository)(nil)

func NewSessionLocalRepository() *SessionLocalRepository {
	return &SessionLocalRepository{
		sessions: make(map[string]models.Session),
	}
}

func (r *SessionLocalRepository) Create() (*models.Session, error) {
	defer r.Unlock()
	r.Lock()
	session := models.NewSession()
	r.sessions[session.Id] = session
	return &session, nil
}

func (r *SessionLocalRepository) Get(id string) (*models.Session, error) {
	defer r.RUnlock()
	r.RLock()
	if session, ok := r.sessions[id]; ok {
		return &session, nil
	}
	return nil, ErrSessionNotFound
}

func (r *SessionLocalRepository) Update(id string, message models.Message) error {
	defer r.Unlock()
	r.Lock()
	if s, ok := r.sessions[id]; ok {
		s.Messages = append([]models.Message{message}, s.Messages...)
		r.sessions[id] = s
		return nil
	}
	return ErrSessionNotFound
}

func (r *SessionLocalRepository) Delete(id string) error {
	defer r.Unlock()
	r.Lock()
	if _, ok := r.sessions[id]; ok {
		delete(r.sessions, id)
		return nil
	}
	return ErrSessionNotFound
}
