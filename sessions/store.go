package sessions

import "sync"

type SessionStore struct {
	sessions map[string]*Session
	mu sync.Mutex
}

func NewSessionsStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *SessionStore) CreateSession() {
	session := &Session{

	}
	s.mu.Lock()
}