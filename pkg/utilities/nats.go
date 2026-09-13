package utilities

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

func NATSConnect(ctx context.Context, url string) (*nats.Conn, error) {
	timeout := 30 * time.Second
	retryDelay := 5 * time.Second

	connCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		conn, err := nats.Connect(url)
		if err == nil {
			return conn, nil
		}

		select {
		case <-connCtx.Done():
			return nil, err
		case <-time.After(retryDelay):
		}
	}
}
