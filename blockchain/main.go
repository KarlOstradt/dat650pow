package main

import (
	"fmt"
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
)

const fileName = "data20.csv"

func main() {
	verbose = false
	t := [][]int64{} // Time vector
	t = append(t, runTest1(2000))
	// t = append(t, runTest2(2000))
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
