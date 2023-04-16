package protocol

import (
	"encoding/binary"
	"errors"
	"strconv"
)

type ProtocolVersion uint8

const (
	ProtocolVersionOne ProtocolVersion = 1
)

const (
	MessageHeaderSize   = 10
	SizeOfUint64InBytes = 8
	SizeOfUint8InBytes  = 1
)

type MessageType uint8

const (
	MessagaTypeProofOfWorkChallengeRequest MessageType = iota
	MessagaTypeProofOfWorkChallengeResponce
	MessageTypeQuote
)

func (t MessageType) String() string {
	switch t {
	case MessagaTypeProofOfWorkChallengeRequest:
		return "ProofOfWorkChallengeRequest"
	case MessagaTypeProofOfWorkChallengeResponce:
		return "ProofOfWorkChallengeResponce"
	case MessageTypeQuote:
		return "MessageTypeQuote"
	default:
		return "UNSUPPORTED"
	}
}

var (
	ErrWrongHeaderSize = errors.New(
		"message header size should be " + strconv.Itoa(MessageHeaderSize))
	ErrUnsupportedProtocolVersion = errors.New("unsupported protocol version")
)

type Serializeable interface {
	MessageType() MessageType
	Serialize([]byte)
	Size() uint64
}

type Message struct {
	Type        MessageType
	Version     ProtocolVersion
	PayloadSize uint64
	Body        Serializeable
}

func NewMessageWithPayload(body Serializeable) *Message {
	return &Message{
		Type:    body.MessageType(),
		Version: ProtocolVersionOne,
		Body:    body,
	}
}

func MessageFromBuf(buf []byte) (*Message, error) {
	if len(buf) != MessageHeaderSize {
		return nil, ErrWrongHeaderSize
	}

	if ProtocolVersion(buf[0]) != ProtocolVersionOne {
		return nil, ErrUnsupportedProtocolVersion
	}

	return &Message{
		Version:     ProtocolVersion(buf[0]),
		Type:        MessageType(buf[1]),
		PayloadSize: binary.BigEndian.Uint64(buf[2:]),
	}, nil
}

func (m *Message) Serialize() []byte {
	buf := make([]byte, 2+SizeOfUint64InBytes+m.Body.Size())

	buf[0] = byte(m.Version)
	buf[1] = byte(m.Type)

	binary.BigEndian.PutUint64(buf[2:2+SizeOfUint64InBytes], uint64(m.Body.Size()))

	m.Body.Serialize(buf[MessageHeaderSize:])

	return buf
}
