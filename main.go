/**
 * see https://gist.github.com/zubietaroberto/0beed814f8df28036af5
 */
package main

import (
	"encoding/hex"
	"context"
	"fmt"
	"net"
  "github.com/orangecms/wsd-cuse/pkg/client"
)

const (
  // wsdAddress = "localhost"
  wsdAddress = "192.160.100.1"
  wsdPort = ":28000"
	clientPort = ":28001"
	protocol = "udp"

  bufSize int = 2048
)

func main() {
	udpAddr, err := net.ResolveUDPAddr(protocol, clientPort)
	if err != nil {
		fmt.Println("Wrong Address")
		return
	}

	udpConn, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		fmt.Println(err)
	}

  ctx := context.Background()
  go client.Send(ctx, wsdAddress + wsdPort)

	display(udpConn)
}

func display(conn *net.UDPConn) {
  buf := make([]byte, bufSize)
	n, err := conn.Read(buf[0:])
	if err != nil {
		fmt.Println("Error Reading")
		return
	} else {
		fmt.Println(hex.EncodeToString(buf[0:n]))
		fmt.Println("Package Done")
	}
}
