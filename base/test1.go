package base

import (
	"encoding/json"
	"fmt"
	"net"
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
		// _, err := chain.MineBlock(txs)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	return nil
		// }
		block := mine(txs, chain.CurrentBlock().Hash)
		chain.blocks = append(chain.blocks, &block)

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

func mine(txs []*Transaction, prevHash []byte) Block {
	addr, err := net.ResolveUDPAddr("udp4", "192.168.39.140:1234")

	s, err := net.ResolveUDPAddr("udp4", ":1234")
	if err != nil {
		fmt.Println(err.Error())
	}
	connection, err := net.ListenUDP("udp4", s)
	defer connection.Close()

	block := Block{PrevBlockHash: prevHash, Transactions: txs}
	block.Timestamp = time.Now().Unix()
	block.Hash = []byte{}
	block.Nonce = -1
	mBlock := MarshalBlock(block)
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()
	// fmt.Printf("%v\n", connection)
	_, err = connection.WriteToUDP(mBlock, addr)
	time.Sleep(time.Second * 5)

	_, err = connection.WriteToUDP([]byte("STP"), addr)
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println("AAA")
	// chanReceive := make(chan []byte)
	buffer := make([]byte, 1024)
	n, _ := connection.Read(buffer)

	err = json.Unmarshal(buffer[3:n], &block)
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println(block.String())
	// fmt.Println(block.Nonce)
	return block
}

// MarshalBlock marshals the block
func MarshalBlock(block Block) []byte {
	mBlock, err := json.Marshal(block)
	if err != nil {
		fmt.Println(err.Error())
	}
	mBlock = append([]byte("POW"), mBlock...)
	return mBlock
}
