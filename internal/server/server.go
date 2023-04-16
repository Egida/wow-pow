package server

import (
	"context"
	"errors"
	"log"
	"net"
	"wow-pow/internal/utils"
)

type Quoter interface {
	Quote() string
}

type Server struct {
	log      *log.Logger
	config   Config
	quoter   Quoter
	socket   net.Listener
	shutdown chan struct{}
}

type Config struct {
	ListenAddr      string
	ProofTokenSize  int
	ProofDifficulty int
}

func New(log *log.Logger, conf Config, quoter Quoter) *Server {
	log = utils.LoggerExtendWithPrefix(log, "server ->")

	return &Server{
		log:      log,
		config:   conf,
		quoter:   quoter,
		shutdown: make(chan struct{}, 1),
	}
}

func (s *Server) Start(ctx context.Context) chan error {
	serverErrors := make(chan error, 1)

	go func() {
		socket, err := net.Listen("tcp", s.config.ListenAddr)
		if err != nil {
			serverErrors <- err

			return
		}

		s.socket = socket

		s.log.Printf("start TCP Sever Listening %s", s.config.ListenAddr)

		serverErrors <- s.serve(ctx, socket)
	}()

	return serverErrors

}

func (s *Server) Shutdown() error {
	s.shutdown <- struct{}{}
	err := s.socket.Close()

	return err
}

func (s *Server) serve(ctx context.Context, socket net.Listener) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.shutdown:
			return nil
		default:
		}

		conn, err := socket.Accept()

		if err != nil && errors.Is(err, net.ErrClosed) {
			return err
		}

		if err != nil {
			s.log.Printf("ERROR: failed to listen socket %s", err)
			continue
		}

		go func(conn net.Conn) {
			log := utils.LoggerExtendWithPrefix(s.log, "rID:"+utils.RequestID())

			if err := s.handleRequest(log, conn); err != nil {
				log.Printf("ERROR: failure to handle incomming request: %s", err)
			}
		}(conn)
	}
}
