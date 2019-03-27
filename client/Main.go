package main

import (
	"Chat/cellnet"
	"Chat/cellnet/packet"
	"Chat/cellnet/socket"
	"Chat/proto"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	queue := cellnet.NewEventQueue()
	peer := socket.NewConnector(packet.NewMessageCallback(onMessage), queue)
	peer.Start("127.0.0.1:8801")
	peer.SetName("client")
	queue.StartLoop()
	ReadConsole(func(line string) bool {
		peer.Send(proto.ChatREQ{Content: line})
		return true
	})
}

func ReadConsole(callback func(line string) bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if !callback(line) {
			break
		}
	}
}

func onMessage(session cellnet.Session, event interface{}) {
	switch event := event.(type) {
	case socket.ConnectedEvent:
		fmt.Println("connected")
	case packet.MsgEvent:
		msg := event.Msg.(*proto.ChatACK)
		// todo log
		fmt.Printf("sessionID%d say: %s\n", msg.ID, msg.Content)
	case socket.SessionClosedEvent:
		fmt.Println("disconnected")
	}
}
