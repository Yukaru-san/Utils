package goutils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// ReadStdinWithTimeout reads from stdin with a given timeout
func ReadStdinWithTimeout(bufferSize int, timeoutSeconds time.Duration) []byte {
	c := make(chan []byte, 1)

	// Read in background to allow using a select for a timeout
	go (func() {
		r := bufio.NewReader(os.Stdin)
		buf := make([]byte, bufferSize)

		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		c <- buf[:n]
	})()

	select {
	case b := <-c:
		return b
	// Timeout
	case <-time.After(timeoutSeconds * time.Second):
		fmt.Println("No input received")
		os.Exit(1)
		return nil
	}
}
