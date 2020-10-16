package base

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	chain          Blockchain
	txBuffer       []*Transaction
	wallet1        *Wallet
	wallet1Address string
	wallet2        *Wallet
	wallet2Address string
	utxos          UTXOSet
	verbose        bool
	slave1         *net.UDPAddr
	slave2         *net.UDPAddr
	conn           *net.UDPConn
)

const nRoutines = 8
const fileName = "data20.csv"

// MainMethod func
func MainMethod() {
	resolveAddresses()
	fmt.Println("MainMethod")
	verbose = false
	t := [][]int64{} // Time vector
	t = append(t, runTest1(10))
	writeToFile(t)
	// fmt.Println(chain.String())
}

// createBlockchain will create new wallets and blockchain
func createBlockchain() {
	wallet1 = NewWallet()
	wallet1Address = wallet1.GetStringAddress()
	wallet2 = NewWallet()
	wallet2Address = wallet2.GetStringAddress()

	chain = *CreateBlockchain(wallet1.GetStringAddress())
	utxos = chain.FindUTXOSet()
	txBuffer = []*Transaction{}
}

func newTransaction(amount int) {
	tx, err := NewUTXOTransaction(wallet1, wallet2Address, amount, utxos, &chain)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	txBuffer = append(txBuffer, tx)
}

// prepareTXs will create a coinbase tx and append the buffer
func prepareTXs() []*Transaction {
	coinbaseTX := NewCoinbaseTX(wallet1Address, GenesisCoinbaseData)
	return append([]*Transaction{coinbaseTX}, txBuffer...)
}

func writeToFile(result [][]int64) {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, t := range result {
		s := strings.Trim(strings.Replace(fmt.Sprint(t), " ", ",", -1), "[]")
		f.WriteString(s + "\n")
	}

}

func resolveAddresses() {
	addr1, err := net.ResolveUDPAddr("udp4", "192.168.39.140:1234")
	if err != nil {
		fmt.Println(err.Error())
		panic("could not resolve 192.168.39.140:1234")
	}
	addr2, err := net.ResolveUDPAddr("udp4", "127.0.0.1:1235")
	if err != nil {
		fmt.Println(err.Error())
		panic("could not resolve :1235")
	}

	s, err := net.ResolveUDPAddr("udp4", ":1234")
	if err != nil {
		fmt.Println(err.Error())
		panic("could not resolve localhost :1234")
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err.Error())
		panic("could not establish listen connection on local port :1234")
	}
	slave1 = addr1
	slave2 = addr2
	conn = connection
}
