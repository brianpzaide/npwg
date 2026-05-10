package main

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = 2 * time.Second

func pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	interval := defaultPingInterval

	select {
	case <-ctx.Done():
		return
	case d := <-reset:
		if d > 0 {
			interval = d
		}
	default:
		if interval <= 0 {
			interval = defaultPingInterval
		}
	}

	timer := time.NewTimer(interval)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case interval = <-reset:
			if interval <= 0 {
				interval = defaultPingInterval
			}
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			if _, err := w.Write([]byte("ping")); err != nil {
				return
			}
		}
		timer.Reset(interval)
	}
}
