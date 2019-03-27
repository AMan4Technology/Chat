package packet

import (
	"Chat/cellnet"
	"Chat/cellnet/socket"
)

func NewMessageCallback(f SessionMessageFunc) cellnet.EventFunc {
	return func(event interface{}) interface{} {
		switch event := event.(type) {
		case socket.RecvEvent:
			return onRecvLTVPacket(event.Session, f)
		case socket.SendEvent:
			return onSendLTVPacket(event.Session, event.Msg)
		case socket.ConnectErrorEvent:
			invokeMsgFunc(event.Session, f, event)
		case socket.SessionStartedEvent:
			invokeMsgFunc(event.Session, f, event)
		case socket.SessionClosedEvent:
			invokeMsgFunc(event.Session, f, event)
		case socket.SessionExitEvent:
			invokeMsgFunc(event.Session, f, event)
		case socket.RecvErrorEvent:
			// todo log
			invokeMsgFunc(event.Session, f, event)
		case socket.SendErrorEvent:
			// todo log
			invokeMsgFunc(event.Session, f, event)
		}
		return nil
	}
}

func invokeMsgFunc(session cellnet.Session, f SessionMessageFunc, event interface{}) {
	queue := session.Peer().Queue()
	if queue != nil {
		queue.Post(func() {
			f(session, event)
		})
		return
	}
	f(session, event)
}

type SessionMessageFunc func(session cellnet.Session, event interface{})

type MsgEvent struct {
	Session cellnet.Session
	Msg     interface{}
}
