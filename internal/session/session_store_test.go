package session

import (
	"testing"
	"time"
)

func TestInMemorySessionStore(t *testing.T) {
	s := NewInMemorySessionStore()
	sessionID := s.CreateSession("test-user-id", "test-organization-id", "test-client-id", time.Minute)
	if sessionID == "" {
		t.Error("Expected a session ID to be returned")
	}
	_, found := s.GetSession(sessionID)
	if !found {
		t.Error("Expected to find the session")
	}
	s.DeleteSession(sessionID)
	_, found = s.GetSession(sessionID)
	if found {
		t.Error("Expected not to find the session")
	}
}
