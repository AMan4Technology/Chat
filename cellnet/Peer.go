package cellnet

type Session interface {
	Send(msg interface{})
	Raw() interface{}
	Close()
	Peer() Peer
	ID() int64
}

type Peer interface {
	Start(address string) bool
	Stop()
	Queue() EventQueue
	SetEventFunc(f EventFunc)
	Name() string
	SetName(string)
	Send(msg interface{})
	SessionAccessor
}

type EventFunc func(event interface{}) interface{}

type SessionAccessor interface {
	GetSession(int64) Session
	RangeSession(func(Session) bool)
	SessionCount() int64
	CloseAllSessions()
}
