package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	chain          Blockchain
	txBuffer       []*Transaction
	wallet1        *Wallet
	wallet1Address string
	wallet2        *Wallet
	wallet2Address string
	utxos          UTXOSet
	t              []int64 // Time vector
)

const fileName = "data.csv"

func main() {
	createWallets()
	createBlockchain()

	mine(10)
	writeToFile()
	// fmt.Println(chain.String())
	// fmt.Printf("%v\n", t)
}

func createBlockchain() {
	chain = *CreateBlockchain(wallet1.GetStringAddress())
	utxos = chain.FindUTXOSet()
}

func mine(n int) {
	for i := 0; i < n; i++ {
		newTransaction(10)
		txs := prepareTXs()
		t0 := time.Now()
		_, err := chain.MineBlock(txs)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t = append(t, time.Now().Sub(t0).Milliseconds())
		utxos.Update(txs)
		txBuffer = []*Transaction{}
	}

}

func newTransaction(amount int) {
	tx, err := NewUTXOTransaction(wallet1, wallet2Address, amount, utxos, &chain)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	txBuffer = append(txBuffer, tx)
}

func prepareTXs() []*Transaction {
	coinbaseTX := NewCoinbaseTX(wallet1Address, GenesisCoinbaseData)
	return append([]*Transaction{coinbaseTX}, txBuffer...)
}

func createWallets() {
	wallet1 = NewWallet()
	wallet1Address = wallet1.GetStringAddress()
	wallet2 = NewWallet()
	wallet2Address = wallet2.GetStringAddress()
}

func writeToFile() {
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	s := strings.Trim(strings.Replace(fmt.Sprint(t), " ", ",", -1), "[]")
	f.WriteString(s)
}
