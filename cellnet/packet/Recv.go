package packet

import (
	"Chat/cellnet"
	"net"
)

func onRecvLTVPacket(session cellnet.Session, f SessionMessageFunc) error {
	conn, ok := session.Raw().(net.Conn)
	if !ok || conn == nil {
		return nil
	}
	pktReader, err := RecvVariableLengthPacket(conn)
	if err != nil {
		return err
	}
	var msgID uint16
	if err := pktReader.ReadValue(&msgID); err != nil {
		return err
	}
	msg, err := cellnet.DecodeMessage(uint32(msgID), pktReader.RemainBytes())
	if err != nil {
		return err
	}
	invokeMsgFunc(session, f, MsgEvent{session, msg})
	return nil
}
