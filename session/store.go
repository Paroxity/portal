package session

import (
	"github.com/google/uuid"
	"strings"
	"sync"
)

// Store represents a store that holds all of the open sessions on the proxy.
type Store interface {
	// All returns all of the sessions stored on the proxy.
	All() []*Session
	// Load attempts to load a session from the UUID of a player.
	Load(x uuid.UUID) (*Session, bool)
	// LoadFromName attempts to load a session from the username of a player.
	LoadFromName(x string) (*Session, bool)
	// Store stores the session on the proxy.
	Store(x *Session)
	// Delete deletes a session from the store.
	Delete(x uuid.UUID)
}

// DefaultStore represents a session store with basic behaviour.
type DefaultStore struct {
	mu           sync.Mutex
	sessions     map[uuid.UUID]*Session
	sessionNames map[string]*Session
}

// NewDefaultStore creates a new DefaultStore and returns it.
func NewDefaultStore() *DefaultStore {
	return &DefaultStore{
		sessions:     make(map[uuid.UUID]*Session),
		sessionNames: make(map[string]*Session),
	}
}

// All ...
func (s *DefaultStore) All() (all []*Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.sessions {
		all = append(all, v)
	}
	return
}

// Load ...
func (s *DefaultStore) Load(x uuid.UUID) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.sessions[x]
	return v, ok
}

// LoadFromName ...
func (s *DefaultStore) LoadFromName(x string) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.sessionNames[strings.ToLower(x)]
	return v, ok
}

// Store ...
func (s *DefaultStore) Store(x *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[x.UUID()] = x
	s.sessionNames[x.Conn().IdentityData().DisplayName] = x
}

// Delete ...
func (s *DefaultStore) Delete(x uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.sessions[x]
	if ok {
		delete(s.sessions, x)
		delete(s.sessionNames, v.Conn().IdentityData().DisplayName)
	}
}
