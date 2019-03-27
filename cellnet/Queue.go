package cellnet

import "runtime/debug"

const DefaultQueueLength = 100

func NewEventQueue() EventQueue {
	return NewEventQueueWithLen(DefaultQueueLength)
}

func NewEventQueueWithLen(length int) EventQueue {
	e := eventQueue{
		queue: make(chan func(), length),
		stop:  make(chan int),
		exit:  make(chan int),
	}
	return &e
}

type EventQueue interface {
	StartLoop()
	StopLoop(result int)
	Wait() int
	Post(callback func())
	EnableCapturePanic(sure bool)
}

type eventQueue struct {
	queue        chan func()
	stop         chan int
	exit         chan int
	capturePanic bool
}

func (e *eventQueue) StartLoop() {
	go func() {
		for callback := range e.queue {
			e.protectedCall(callback)
		}
		e.exit <- <-e.stop
	}()
}

func (e *eventQueue) StopLoop(result int) {
	close(e.queue)
	e.queue = make(chan func(), cap(e.queue))
	e.stop <- result
}

func (e *eventQueue) Wait() int {
	return <-e.exit
}

func (e *eventQueue) Post(callback func()) {
	if callback == nil {
		return
	}
	e.queue <- callback
}

func (e *eventQueue) EnableCapturePanic(sure bool) {
	e.capturePanic = sure
}

func (e eventQueue) protectedCall(callback func()) {
	if callback == nil {
		return
	}
	if e.capturePanic {
		defer func() {
			if p := recover(); p != nil {
				debug.PrintStack()
			}
		}()
	}
	callback()
}
