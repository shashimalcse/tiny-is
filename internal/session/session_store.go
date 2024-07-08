package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type SessionInfo struct {
	SessionId      string
	UserID         string
	OrganizationId string
	ClientID       string
	ExpiresAt      time.Time
}

type SessionStore interface {
	CreateSession(userID, OrganizationId, clientID string, expireTime time.Duration) string
	GetSession(sessionID string) (SessionInfo, bool)
	DeleteSession(sessionID string)
}

type inMemorySessionStore struct {
	c *cache.Cache
}

func NewInMemorySessionStore() SessionStore {
	return &inMemorySessionStore{
		c: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (s inMemorySessionStore) CreateSession(userID, OrganizationId, clientID string, expireTime time.Duration) string {
	sessionID := uuid.New().String()
	s.c.Set(sessionID, SessionInfo{
		UserID:         userID,
		OrganizationId: OrganizationId,
		ClientID:       clientID,
		ExpiresAt:      time.Now().Add(expireTime),
		SessionId:      sessionID,
	}, expireTime)
	return sessionID
}

func (s inMemorySessionStore) GetSession(sessionID string) (SessionInfo, bool) {
	if data, found := s.c.Get(sessionID); found {
		sessionInfo := data.(SessionInfo)
		if time.Now().Before(sessionInfo.ExpiresAt) {
			return sessionInfo, true
		}
		s.c.Delete(sessionID)
	}
	return SessionInfo{}, false
}

func (s inMemorySessionStore) DeleteSession(sessionID string) {
	s.c.Delete(sessionID)
}
