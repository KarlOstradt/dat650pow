package main

import (
	"fmt"
	"time"
)

// Wallet 1 always sends 10 coin to Wallet 2
// One input tx, one output tx
func runTest1(n int) []int64 {
	createBlockchain()
	t := []int64{}
	if verbose {
		balance1, _ := utxos.FindSpendableOutputs(HashPubKey(wallet1.PublicKey), 9999999)
		balance2, _ := utxos.FindSpendableOutputs(HashPubKey(wallet2.PublicKey), 9999999)
		fmt.Printf("%d %d %d\n", len(chain.blocks), balance1, balance2)
	}

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
		if verbose {
			balance1, _ := utxos.FindSpendableOutputs(HashPubKey(wallet1.PublicKey), 9999999)
			balance2, _ := utxos.FindSpendableOutputs(HashPubKey(wallet2.PublicKey), 9999999)
			fmt.Printf("%d %d %d\n", len(chain.blocks), balance1, balance2)
		}
		txBuffer = []*Transaction{}
	}

	return t
}
