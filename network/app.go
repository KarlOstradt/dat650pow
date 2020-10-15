package network

import (
	"dat650/base"
	"encoding/json"
	"fmt"
	"net"
)

const (
	port = ":1234"
)

var ourID int = -1
var networkNodes = []int{}
var addresses = []string{}

func main() {

	// Setup
	ourID, networkNodes, addresses = setup()

	// Resolve UDP address
	s, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		fmt.Println(err)
	}

	// Connection to listen for requests / commands
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
	}

	// Start go rutines
	go handleRequest(connection)
	// go handleResponse(connection)

	stop := make(chan struct{})
	select {
	case <-stop:
		fmt.Println("Stopping application")
	}
}

func handleRequest(connection *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err.Error())
		}
		tag := string(buffer[0:3])
		switch tag {
		case "STP":
			//TODO: Notify node to stop solving pow
		case "POW":
			//TODO: Update blockchain
			var block base.Block
			err = json.Unmarshal(buffer[3:n], &block)
			if err != nil {
				fmt.Println(err.Error())
			}

			block.Mine()
			sendResponse(connection, block)
		}
	}
}

func sendResponse(connection *net.UDPConn, block base.Block) {
	// buffer
	// buffer := make([]byte, 2000)
	addr, _ := net.ResolveUDPAddr("udp4", "192.168.39.135:1234")
	connection.WriteToUDP(base.MarshalBlock(block), addr)
}
