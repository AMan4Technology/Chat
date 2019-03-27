package socket

import (
	"net"

	"Chat/cellnet"
	"Chat/cellnet/internal"
)

func NewConnector(f cellnet.EventFunc, queue cellnet.EventQueue) cellnet.Peer {
	return &connector{
		peer: peer{
			eventFunc:      f,
			queue:          queue,
			SessionManager: internal.NewSessionManager(),
		},
	}
}

type connector struct {
	peer
	session cellnet.Session
}

func (c *connector) Start(address string) bool {
	if c.peer.Start(address) {
		go c.connect(address)
		return true
	}
	return false
}

func (c *connector) Stop() {
	if c.session != nil {
		c.session = nil
		c.peer.Stop()
	}
}

func (c *connector) connect(address string) {
	conn, err := net.Dial("tcp", address)
	session := newSession(conn, c)
	c.session = session
	if err != nil {
		c.fireEvent(ConnectErrorEvent{session, err})
		c.Stop()
		return
	}
	// todo log
	session.start()
}
