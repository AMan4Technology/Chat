package socket

import (
	"fmt"
	"net"
	"sync"

	"Chat/cellnet"
	"Chat/cellnet/internal"
)

func NewAcceptor(f cellnet.EventFunc, queue cellnet.EventQueue) cellnet.Peer {
	return &acceptor{
		peer: peer{
			eventFunc:      f,
			queue:          queue,
			SessionManager: internal.NewSessionManager(),
		},
	}
}

type acceptor struct {
	peer
	listener net.Listener
	wg       sync.WaitGroup
}

func (a *acceptor) Start(address string) bool {
	if a.peer.Start(address) {
		go a.listen(address)
		return true
	}
	return false
}

func (a *acceptor) Stop() {
	if a.listener != nil {
		a.listener.Close()
		a.listener = nil
		a.wg.Wait()
		a.peer.Stop()
	}
}

func (a acceptor) Send(msg interface{}) {
	a.peer.Send(msg)
}

func (a *acceptor) listen(address string) {
	a.wg.Add(1)
	defer a.wg.Done()
	var err error
	a.listener, err = net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := a.listener.Accept()
		if err != nil {
			break
		}
		go a.handleSession(conn)
	}
}

func (a *acceptor) handleSession(conn net.Conn) {
	session := newSession(conn, a)
	session.start()
}
