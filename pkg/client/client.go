package client

import (
	"fmt"
	"log"
	"net"
	"wow-pow/internal/pow"
	"wow-pow/internal/utils"
	"wow-pow/pkg/protocol"
)

type Client struct {
	log    *log.Logger
	config Config
	socket net.Listener
}

type Config struct {
	Address        string
	PoWConcurrency int
}

func New(log *log.Logger, conf Config) *Client {
	log = utils.LoggerExtendWithPrefix(log, "client ->")

	return &Client{
		log:    log,
		config: conf,
	}
}

func (c *Client) GetWorkdOfWisdomQuote() (string, error) {
	conn, err := net.Dial("tcp", c.config.Address)
	if err != nil {
		return "", err
	}

	c.log.Printf("establish connection %s -> %s",
		conn.LocalAddr().String(), conn.RemoteAddr().String())

	err = c.ddosProtection(conn)
	if err != nil {
		return "", fmt.Errorf("did not pass DDoS protection: %e", err)
	}

	quote, err := c.readQuote(conn)
	if err != nil {
		return "", fmt.Errorf("failure to receive a quote: %e", err)
	}

	return quote, nil
}

func (c *Client) ddosProtection(conn net.Conn) error {
	payload, err := readMessage(
		conn, protocol.MessagaTypeProofOfWorkChallengeRequest)
	if err != nil {
		return err
	}

	req := &protocol.ProofOfWorkChallengeRequest{}
	req.Deserialize(payload)

	c.log.Printf("got PoW challenge with difficulty: %d, token: %v",
		req.Difficulty, req.Token)

	result := pow.FindNonce(req.Difficulty, req.Token, c.config.PoWConcurrency)
	if result.Error != nil {
		return result.Error

	}

	c.log.Printf("find the solution for %v nonce: %d leading zeros: %d hash: %v ",
		result.Duration, result.Nonce, result.LeadingZeros, result.Hash)

	resp := protocol.NewMessageWithPayload(
		protocol.NewProofOfWorkChallengeResponce(result.Nonce),
	)

	if _, err := conn.Write(resp.Serialize()); err != nil {
		return err
	}

	return nil
}

func (c *Client) readQuote(conn net.Conn) (string, error) {
	payload, err := readMessage(conn, protocol.MessageTypeQuote)
	if err != nil {
		return "", err
	}

	msg := &protocol.Quote{}
	msg.Deserialize(payload)

	c.log.Printf("got the quote form Word Of Wisdom: %s", msg.Quote)

	return msg.Quote, nil
}
