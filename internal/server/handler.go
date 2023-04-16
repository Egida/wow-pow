package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"wow-pow/pkg/protocol"
)

func (s *Server) handleRequest(log *log.Logger, conn net.Conn) error {
	defer func() {
		if err := conn.Close(); err != nil {
			s.log.Printf("ERROR: failure to close the connection: %e", err)
		}
	}()

	log.Printf("incomming connection from %s", conn.RemoteAddr().String())

	// DDoS protection
	pass, err := s.ddosProtection(log, conn)
	if err != nil {
		return fmt.Errorf("failure to do DDoS protection check: %e", err)
	}

	if !pass {
		return errors.New("DDoS protection check did not pass")
	}

	// Main login
	quote := s.quoter.Quote()
	msg := protocol.NewMessageWithPayload(protocol.NewQuote(quote))

	if _, err := conn.Write(msg.Serialize()); err != nil {
		return errors.New("failure to send a quote")
	}

	log.Printf("send the quote '%s'", quote)

	return nil
}
