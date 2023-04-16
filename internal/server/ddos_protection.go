package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"wow-pow/internal/pow"
	"wow-pow/pkg/protocol"
)

func (s *Server) ddosProtection(log *log.Logger, conn net.Conn) (bool, error) {
	difficulty := uint8(s.config.ProofDifficulty)
	token := make([]byte, s.config.ProofTokenSize)

	if err := pow.GenerateToken(token); err != nil {
		return false, fmt.Errorf("failure to generate proof token: %s", err)
	}

	// Send puzzle
	msg := protocol.NewMessageWithPayload(
		protocol.NewProofOfWorkChallengeRequest(difficulty, token),
	)

	if _, err := conn.Write(msg.Serialize()); err != nil {
		return false, err
	}

	log.Printf("sent PoW challenge with difficulty: %d, token: %v",
		difficulty, token)

	// Receive the solution
	head := make([]byte, protocol.MessageHeaderSize)

	if _, err := conn.Read(head); err != nil {
		return false, err
	}

	msg, err := protocol.MessageFromBuf(head)
	if err != nil {
		return false, err
	}

	if msg.Type != protocol.MessagaTypeProofOfWorkChallengeResponce {
		return false, errors.New("expected a ProofOfWorkChallengeResponce message")
	}

	payload := make([]byte, msg.PayloadSize)
	if _, err = conn.Read(payload); err != nil {
		return false, err
	}

	resp := &protocol.ProofOfWorkChallengeResponce{}
	resp.Deserialize(payload)

	log.Printf("got the solution nonce: %d", resp.Nonce)

	if pow.CheckSolution(token, resp.Nonce, difficulty) {
		return true, nil
	}

	return false, nil
}
