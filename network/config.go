package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Servers - List of all servers
type Servers struct {
	Servers []Server `json:"servers"`
}

// Server info
type Server struct {
	Hostname string `json:"hostname"`
	Address  string `json:"address"`
	ID       int    `json:"id"`
}

var test = []string{}

func returnAddresses() []string {
	return test
}

// Function for network setup.
func setup() (int, []int, []string) {
	hostname, _ := os.Hostname()
	hostname = strings.Split(hostname, ".")[0]

	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var servers Servers

	json.Unmarshal(byteValue, &servers)

	var nodes = []int{}
	var address = []string{}
	var id int

	for i := 0; i < len(servers.Servers); i++ {
		nodes = append(networkNodes, servers.Servers[i].ID)
		address = append(addresses, servers.Servers[i].Address)

		if hostname == servers.Servers[i].Hostname {
			id = servers.Servers[i].ID
		}
	}

	// Localhost scenario
	// if ourID == -5 {
	// 	ourID = 2
	// 	hostname = "localhost:1234"
	// 	addresses = []string{"localhost:1234", "localhost:1234", "localhost:1234"}
	// }
	return id, nodes, address
}
