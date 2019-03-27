package socket

import (
	"Chat/cellnet"
	"Chat/cellnet/internal"
)

type eventWorker interface {
	fireEvent(event interface{}) interface{}
}

type peer struct {
	eventFunc cellnet.EventFunc
	queue     cellnet.EventQueue
	name      string
	address   string
	internal.SessionManager
}

func (s *peer) Start(address string) bool {
	if s.address != "" {
		return false
	}
	s.address = address
	return true
}

func (s *peer) Stop() {
	s.CloseAllSessions()
	s.address = ""
}

func (s peer) Name() string {
	return s.name
}

func (s *peer) SetName(name string) {
	s.name = name
}

func (s peer) Queue() cellnet.EventQueue {
	return s.queue
}

func (s *peer) SetEventFunc(f cellnet.EventFunc) {
	s.eventFunc = f
}

func (s peer) fireEvent(event interface{}) interface{} {
	if s.eventFunc == nil {
		return nil
	}
	return s.eventFunc(event)
}

func (s peer) Send(msg interface{}) {
	s.RangeSession(func(session cellnet.Session) bool {
		session.Send(msg)
		return true
	})
}
