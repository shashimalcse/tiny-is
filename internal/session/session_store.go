package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type SessionInfo struct {
	SessionId string
	UserID    string
	ClientID  string
	ExpiresAt time.Time
}

type SessionStore struct {
	cache *cache.Cache
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (s *SessionStore) CreateSession(userID, clientID string, expireTime time.Duration) string {
	sessionID := uuid.New().String()
	s.cache.Set(sessionID, SessionInfo{
		UserID:    userID,
		ClientID:  clientID,
		ExpiresAt: time.Now().Add(expireTime),
		SessionId: sessionID,
	}, expireTime)
	return sessionID
}

func (s *SessionStore) GetSession(sessionID string) (*SessionInfo, bool) {
	if data, found := s.cache.Get(sessionID); found {
		sessionInfo := data.(SessionInfo)
		if time.Now().Before(sessionInfo.ExpiresAt) {
			return &sessionInfo, true
		}
		// Session has expired, remove it
		s.cache.Delete(sessionID)
	}
	return nil, false
}

func (s *SessionStore) DeleteSession(sessionID string) {
	s.cache.Delete(sessionID)
}
