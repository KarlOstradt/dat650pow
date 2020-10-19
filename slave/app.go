package main

import (
	"dat650/base"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const (
	// Eirik
	// port = ":1234"

	// Karl
	port = ":1235"
)

// Karl
var ourID int = 1

// Eirik
// var ourID int = 0

var networkNodes = []int{}
var addresses = []string{}

func main() {

	fmt.Println("Program started")

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

	fmt.Println("Setup and address resolved")

	// Start go rutines
	go handleRequest(connection)

	stop := make(chan struct{})
	select {
	case <-stop:
		fmt.Println("Stopping application")
	}
}

func handleRequest(connection *net.UDPConn) {
	fmt.Println("HandleRequest")
	defer func(connection *net.UDPConn) {
		if err := recover(); err != nil {
			fmt.Println("Crashed. Starting new handleRequest routine.")
			go handleRequest(connection)
		}
	}(connection)
	buffer := make([]byte, 1024)
	stopChan := make(chan bool)
	close(stopChan)
	for {
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err.Error())
		}
		tag := string(buffer[0:3])

		switch tag {
		case "STP":
			_, ok := <-stopChan
			fmt.Println(ok)
			if ok {
				stopChan <- true
			}
		case "POW":
			// Attempt to stop previous pow
			ticker := time.NewTicker(50 * time.Millisecond)
			select {
			case _, ok := <-stopChan:
				if ok {
					stopChan <- true
				}
			case <-ticker.C:
				stopChan <- true
			}
			ticker.Stop()

			stopChan = make(chan bool)
			var block base.Block
			err = json.Unmarshal(buffer[3:n], &block)
			if err != nil {
				fmt.Println(err.Error())
			}

			go func(block base.Block) {
				block.Mine(stopChan, ourID, 2)
				if block.Nonce != -1 {
					sendResponse(connection, block)
				}
				defer func() {
					if err := recover(); err != nil {

					}
				}()
				close(stopChan)
			}(block)

		}
	}
}

func sendResponse(connection *net.UDPConn, block base.Block) {
	// Karl
	addr, _ := net.ResolveUDPAddr("udp4", ":1234")

	// Eirik
	// addr, _ := net.ResolveUDPAddr("udp4", "192.168.39.135:1234")
	connection.WriteToUDP(base.MarshalBlock(block), addr)
}
