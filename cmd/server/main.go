package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wow-pow/internal/quotes"
	"wow-pow/internal/server"
)

func main() {
	var (
		quotesFile = flag.String(
			"quotes", "data/quotes.txt", "file with quotes")
		listenAddr = flag.String("listen", "0.0.0.0:9000",
			"TCP/IP address on which the server is to listen for connections")
		proofDifficulty = flag.Int("proof-difficulty", 24, "proof difficulty (bits)")
		proofTokenSize  = flag.Int("proof-token-size", 64, "proof token size (bytes)")
	)

	flag.Parse()

	log := log.New(
		os.Stdout,
		"SERVER",
		log.Lshortfile|log.Lmicroseconds,
	)

	qs, err := loadQuotes(*quotesFile)
	if err != nil {
		log.Printf("failure to load quotes: %s", err)
		os.Exit(1)
	}

	quoteKeeper := quotes.New(qs)

	server := server.New(log, server.Config{
		ListenAddr:      *listenAddr,
		ProofTokenSize:  *proofTokenSize,
		ProofDifficulty: *proofDifficulty,
	}, quoteKeeper)

	ctx := context.Background()

	errTCPServer := server.Start(ctx)

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errTCPServer:
		log.Printf("problem with TCP Server %s", err)
		ctx.Done()
	case <-osSignals:
		log.Print("shutdown the server")
		ctx.Done()

		if err := server.Shutdown(); err != nil {
			log.Printf("ERROR: failure to shutdown TCP Server: %s", err)
		}
	}
}

func loadQuotes(filename string) ([]string, error) {
	result := []string{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		quote := scanner.Text()

		result = append(result, quote)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil

}
