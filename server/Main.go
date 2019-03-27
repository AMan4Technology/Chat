package main

import (
	"Chat/cellnet"
	"Chat/cellnet/packet"
	"Chat/cellnet/socket"
	"Chat/proto"
	"fmt"
)

func main() {
	queue := cellnet.NewEventQueue()
	peer := socket.NewAcceptor(packet.NewMessageCallback(onMessage), queue)
	peer.Start("127.0.0.1:8801")
	peer.SetName("server")
	queue.StartLoop()
	queue.Wait()
	peer.Stop()
}

func onMessage(session cellnet.Session, event interface{}) {
	switch event := event.(type) {
	case socket.AcceptedEvent:
		// todo log
		fmt.Println("client accepted: ", session.ID())
	case packet.MsgEvent:
		msg := event.Msg.(*proto.ChatREQ)
		ack := proto.ChatACK{
			ID:      session.ID(),
			Content: msg.Content,
		}
		session.Peer().Send(ack)
	case socket.SessionClosedEvent:
		fmt.Println("client disconnected: ", session.ID())
	}
}
