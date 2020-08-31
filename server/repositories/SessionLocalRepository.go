package stores

import (
	"air-sync/models"
	"sync"
)

type SessionMutex struct {
	sync.Mutex
	models.Session
}

type SessionLocalRepository struct {
	sync.RWMutex
	sessions map[string]*SessionMutex
}

var _ SessionRepository = (*SessionLocalRepository)(nil)

func NewSessionLocalRepository() *SessionLocalRepository {
	return &SessionLocalRepository{
		sessions: make(map[string]*SessionMutex),
	}
}

func (r *SessionLocalRepository) Create() (*models.Session, error) {
	defer r.Unlock()
	r.Lock()
	session := models.NewSession()
	r.sessions[session.Id] = &SessionMutex{
		Session: session,
	}
	return &session, nil
}

func (r *SessionLocalRepository) Get(id string) (*models.Session, error) {
	defer r.RUnlock()
	r.RLock()
	if session, ok := r.sessions[id]; ok {
		return &session.Session, nil
	}
	return nil, ErrSessionNotFound
}

func (r *SessionLocalRepository) InsertMessage(id string, message models.Message) error {
	defer r.Unlock()
	r.Lock()

	s, ok := r.sessions[id]
	if !ok {
		return ErrSessionNotFound
	}
	defer s.Unlock()
	s.Lock()

	s.Messages = append([]models.Message{message}, s.Messages...)
	r.sessions[id] = s
	return nil
}

func (r *SessionLocalRepository) DeleteMessage(id string, messageId string) error {
	defer r.Unlock()
	r.Lock()

	s, ok := r.sessions[id]
	if !ok {
		return ErrSessionNotFound
	}
	defer s.Unlock()
	s.Lock()

	for index, message := range s.Messages {
		if message.Id != messageId {
			continue
		}
		messages := s.Messages
		temp := make([]models.Message, 0)
		if index > 0 {
			temp = append(temp, messages[0:index]...)
		}
		temp = append(temp, messages[index+1:]...)
		s.Messages = temp
		return nil
	}

	return ErrMessageNotFound
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
