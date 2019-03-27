package socket

import "Chat/cellnet"

type RecvEvent struct {
	Session cellnet.Session
}

type SendEvent struct {
	Session cellnet.Session
	Msg     interface{}
}

type RecvErrorEvent struct {
	Session cellnet.Session
	Error   error
}

type SendErrorEvent struct {
	Session cellnet.Session
	Msg     interface{}
	Error   error
}

type ConnectErrorEvent struct {
	Session cellnet.Session
	Error   error
}

type SessionStartedEvent struct {
	Session cellnet.Session
}

type SessionClosedEvent struct {
	Session cellnet.Session
	Error   error
}

type ConnectedEvent = SessionStartedEvent

type AcceptedEvent = SessionStartedEvent

type SessionExitEvent struct {
	Session cellnet.Session
}
