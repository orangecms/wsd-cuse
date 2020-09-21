/**
 * see https://ops.tips/blog/udp-client-and-server-in-go/
 */
package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

const (
	maxBufferSize = 2048
)

// client wraps the whole functionality of a UDP client that sends
// a message and waits for a response coming back from the server
// that it initially targetted.
func Send(ctx context.Context, address string) (err error) {
	// Resolve the UDP address so that we can make use of DialUDP
	// with an actual IP and port instead of a name (in case a
	// hostname is specified).
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}

	// Although we're not in a connection-oriented transport,
	// the act of `dialing` is analogous to the act of performing
	// a `connect(2)` syscall for a socket of type SOCK_DGRAM:
	// - it forces the underlying socket to only read and write
	//   to and from a specific remote address.
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}

	// Closes the underlying file descriptor associated with the,
	// socket so that it no longer refers to any file.
	defer conn.Close()

	doneChan := make(chan error, 1)

	go func() {
		// It is possible that this action blocks, although this
		// should only occur in very resource-intensive situations:
		// - when you've filled up the socket buffer and the OS
		//   can't dequeue the queue fast enough.
		//n, err := io.Copy(conn, reader)
		//n, err := fmt.Fprintf(conn, "WSD")
		n, err := fmt.Fprintf(conn,
			"FC1307"+string(0x1)+string(4)+string(0)+string(0)+string(0)+string(0)+string(0)+string(1)+string(5)+string(5)+"adminxxxxxxxxxxxadminxxxxxxxxxxx"+string(0)+string(0)+string(0)+string(1),
		)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("%d bytes sent\n", n)

		buffer := make([]byte, maxBufferSize)

		// Set a deadline for the ReadOperation so that we don't
		// wait forever for a server that might not respond on
		// a resonable amount of time.
		deadline := time.Now().Add(15 * time.Second)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			doneChan <- err
			return
		}

		nRead, _, err := conn.ReadFrom(buffer)
		if err != nil {
			doneChan <- err
			return
		}

		fmt.Printf("bytes=%d received\n", nRead)
		fmt.Println(hex.EncodeToString(buffer[0:nRead]))

		doneChan <- nil
	}()
	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}

/*

// Call the `Write()` method of the implementor
// of the `io.Writer` interface.
n, err = fmt.Fprintf(conn, "something")
*/
