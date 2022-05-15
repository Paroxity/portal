package session

import (
	"github.com/google/uuid"
	"sync"
)

// Store represents a store which holds all the open sessions on the proxy.
type Store struct {
	mu           sync.Mutex
	sessions     map[uuid.UUID]*Session
	sessionNames map[string]*Session
}

// NewDefaultStore creates a new Store and returns it.
func NewDefaultStore() *Store {
	return &Store{
		sessions:     make(map[uuid.UUID]*Session),
		sessionNames: make(map[string]*Session),
	}
}

// All returns all the sessions stored on the proxy.
func (s *Store) All() (all []*Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.sessions {
		all = append(all, v)
	}
	return
}

// Load attempts to load a session from the UUID of a player.
func (s *Store) Load(x uuid.UUID) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.sessions[x]
	return v, ok
}

// LoadFromName attempts to load a session from the username of a player.
func (s *Store) LoadFromName(x string) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.sessionNames[x]
	return v, ok
}

// Store stores the session on the proxy.
func (s *Store) Store(x *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[x.UUID()] = x
	s.sessionNames[x.Conn().IdentityData().DisplayName] = x
}

// Delete deletes a session from the store.
func (s *Store) Delete(x uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.sessions[x]
	if ok {
		delete(s.sessions, x)
		delete(s.sessionNames, v.Conn().IdentityData().DisplayName)
	}
}
