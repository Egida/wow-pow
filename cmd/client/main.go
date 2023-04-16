package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wow-pow/pkg/client"
)

func main() {
	var (
		address = flag.String("listen", "0.0.0.0:9000",
			"TCP/IP address of the server")
		powConcurrency = flag.Int("pow-concurrency", 4,
			"how many cores should be used for solve a pow puzzle")
		fetchConcurrency = flag.Int("fetch-concurrency", 1,
			"number of simultaneously client requests")
		pauseBetweenCallsSecs = flag.Int("pause-between-calls", 2,
			"pause between client calls in seconds")
	)

	flag.Parse()

	log := log.New(
		os.Stdout,
		"CLIENT",
		log.Lshortfile|log.Lmicroseconds,
	)

	client := client.New(log, client.Config{
		Address:        *address,
		PoWConcurrency: *powConcurrency,
	})

	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < *fetchConcurrency; i++ {
		go fetch(ctx, client, time.Second*time.Duration(*pauseBetweenCallsSecs))
	}

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case <-osSignals:
		cancel()
	}
}

func fetch(ctx context.Context, client *client.Client, pause time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			quote, err := client.GetWorkdOfWisdomQuote()
			if err != nil {
				log.Printf("ERROR: %s", err)
			}

			fmt.Printf("Quote: %s\n", quote)

			time.Sleep(pause)
		}
	}
}
