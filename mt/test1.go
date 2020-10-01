package main

import (
	"fmt"
	"time"
)

func runTest1(n int) []int64 {
	createBlockchain()
	t := []int64{}

	for i := 0; i < n; i++ {
		newTransaction(10) // Send 10
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
