package main

import (
	"fmt"
	"time"
)

// Wallet 1 always sends 1 coin to Wallet 2
// One input tx, two output tx
// Every 10 transactions will be 1:1 ?
func runTest2(n int) []int64 {
	createBlockchain()
	t := []int64{}

	for i := 0; i < n; i++ {
		newTransaction(1) // Send 10
		txs := prepareTXs()
		t0 := time.Now()
		_, err := chain.MineBlock(txs)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		t = append(t, time.Now().Sub(t0).Milliseconds())
		utxos.Update(txs)
		txBuffer = []*Transaction{}
	}

	return t
}
