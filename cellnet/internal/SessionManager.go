package internal

import (
	"sync"
	"sync/atomic"

	"Chat/cellnet"
)

func NewSessionManager() SessionManager {
	return new(sessionManager)
}

type SessionManager interface {
	cellnet.SessionAccessor
	Add(session cellnet.Session)
	Remove(session cellnet.Session)
}

type sessionManager struct {
	sessionWithID sync.Map
	sessionIDGen  int64
	count         int64
}

func (s *sessionManager) GetSession(id int64) cellnet.Session {
	if v, ok := s.sessionWithID.Load(id); ok {
		return v.(cellnet.Session)
	}
	return nil
}

func (s *sessionManager) RangeSession(callback func(cellnet.Session) bool) {
	s.sessionWithID.Range(func(key, value interface{}) bool {
		return callback(value.(cellnet.Session))
	})
}

func (s *sessionManager) SessionCount() int64 {
	return atomic.LoadInt64(&s.count)
}

func (s *sessionManager) CloseAllSessions() {
	s.RangeSession(func(session cellnet.Session) bool {
		session.Close()
		return true
	})
}

func (s *sessionManager) Add(session cellnet.Session) {
	id := atomic.AddInt64(&s.sessionIDGen, 1)
	atomic.AddInt64(&s.count, 1)
	session.(interface {
		SetID(id int64)
	}).SetID(id)
	s.sessionWithID.Store(id, session)
}

func (s *sessionManager) Remove(session cellnet.Session) {
	s.sessionWithID.Delete(session.ID())
	atomic.AddInt64(&s.count, -1)
}
