package socket

import (
	"Chat/cellnet"
	"Chat/cellnet/internal"
	"io"
	"net"
	"sync"
)

const SendQueueLen = 10

func newSession(conn net.Conn, peer cellnet.Peer) *session {
	return &session{
		conn:     conn,
		peer:     peer,
		sendChan: make(chan interface{}, SendQueueLen),
	}
}

type session struct {
	conn              net.Conn
	closeOnce         sync.Once
	exitWG            sync.WaitGroup
	peer              cellnet.Peer
	id                int64
	sendChan          chan interface{}
	closeSendChanOnce sync.Once
}

func (s *session) Send(msg interface{}) {
	defer func() { recover() }()
	s.sendChan <- msg
}

func (s session) Raw() interface{} {
	return s.conn
}

func (s *session) Close() {
	s.sessionManager().Remove(s)
	// defer func() { recover() // 通过recover()恢复重复关闭channel导致的panic }()
	s.closeSendChanOnce.Do(func() {
		close(s.sendChan)
	})
}

func (s session) Peer() cellnet.Peer {
	return s.peer
}

func (s session) ID() int64 {
	return s.id
}

func (s *session) SetID(id int64) {
	s.id = id
}

func (s *session) start() {
	eventWorker := s.eventWorker()
	s.sessionManager().Add(s)
	eventWorker.fireEvent(SessionStartedEvent{s})
	s.exitWG.Add(2)
	go func() {
		s.exitWG.Wait()
		eventWorker.fireEvent(SessionExitEvent{s})
	}()
	go s.sendLoop()
	go s.recvLoop()
}

func (s *session) sendLoop() {
	var err interface{}
	defer func() { s.cleanup(err) }()
	eventWorker := s.eventWorker()
	for msg := range s.sendChan {
		err = eventWorker.fireEvent(SendEvent{s, msg})
		if err != nil {
			eventWorker.fireEvent(SendErrorEvent{s, msg, err.(error)})
			break
		}
	}
}

func (s *session) recvLoop() {
	var err interface{}
	defer func() { s.cleanup(err) }()
	eventWorker := s.eventWorker()
	for {
		err = eventWorker.fireEvent(RecvEvent{s})
		if err == io.EOF {
			break
		} else if err != nil && s.conn != nil {
			eventWorker.fireEvent(RecvErrorEvent{s, err.(error)})
			break
		}
	}
}

func (s *session) cleanup(err interface{}) {
	defer func() {
		s.exitWG.Done()
	}()
	s.closeOnce.Do(func() {
		if s.conn != nil {
			s.Close()
			var conn net.Conn
			conn, s.conn = s.conn, nil
			conn.Close()
			s.eventWorker().fireEvent(SessionClosedEvent{s, err.(error)})
		}
	})
}

func (s session) sessionManager() internal.SessionManager {
	return s.peer.(internal.SessionManager)
}

func (s session) eventWorker() eventWorker {
	return s.peer.(eventWorker)
}
