package main

import (
	"context"
	"log"
	"net"
	"time"
)

func handleConnection(mainCtx context.Context, conn net.Conn) {
	begin := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		conn.Close()
	}()

	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second
	go pinger(ctx, conn, resetTimer)

	err := conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Println(err)
		return
	}
	buf := make([]byte, 1024)
	for {
		select {
		case <-mainCtx.Done():
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("connection was closed", err)
				return
			}
			log.Printf("[%s] %s",
				time.Since(begin).Truncate(time.Second), buf[:n])
			resetTimer <- 0
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
