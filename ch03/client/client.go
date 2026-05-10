package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	begin := time.Now()
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}
	_, err = conn.Write([]byte("Pong"))
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
				return
			}
			break
		}
		log.Printf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}
}
