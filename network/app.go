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
var blockchain base.Blockchain

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
	go handleResponse(connection)

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
		case "MSG":
			//TODO: Notify node to stop solving pow
		case "UPD":
			//TODO: Update blockchain
			var newBc base.Blockchain
			err = json.Unmarshal(buffer[3:n], &newBc)
			if err != nil {
				fmt.Println(err.Error())
			}
			blockchain = newBc
		case "POW":
			//TODO: Solve pow
		}
	}
}

func handleResponse(connection *net.UDPConn) {
	// buffer
	// buffer := make([]byte, 2000)
}
