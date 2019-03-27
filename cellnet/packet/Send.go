package packet

import (
	"Chat/cellnet"
	"net"
)

func onSendLTVPacket(session cellnet.Session, msg interface{}) error {
	conn, ok := session.Raw().(net.Conn)
	if !ok || conn == nil {
		return nil
	}
	data, msgID, err := cellnet.EncodeMessage(msg)
	if err != nil {
		return err
	}
	var pktWriter Writer
	if err := pktWriter.WriteValue(uint16(msgID)); err != nil {
		return err
	}
	if err := pktWriter.WriteValue(data); err != nil {
		return err
	}
	return SendVariableLengthPacket(conn, pktWriter)
}
