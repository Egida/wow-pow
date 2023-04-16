package client

import (
	"errors"
	"net"
	"wow-pow/pkg/protocol"
)

func readMessage(
	conn net.Conn, messageType protocol.MessageType,
) ([]byte, error) {
	head := make([]byte, protocol.MessageHeaderSize)

	if _, err := conn.Read(head); err != nil {
		return nil, err
	}

	msg, err := protocol.MessageFromBuf(head)
	if err != nil {
		return nil, err
	}

	if msg.Type != messageType {
		return nil, errors.New("expected a " + messageType.String() + " message")
	}

	payload := make([]byte, msg.PayloadSize)
	if _, err = conn.Read(payload); err != nil {
		return nil, err
	}

	return payload, nil
}
