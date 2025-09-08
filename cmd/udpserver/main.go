package main

import (
	"log"
	"net"
)

func main() {
	UDPAddress, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Couldn't resolve UDP address: %v", err)
	}

	UDPConn, err := net.DialUDP("udp", UDPAddress, nil)
	if err != nil {
		log.Fatalf("Coudn't open UDP connection: %v", err)
	}

}
