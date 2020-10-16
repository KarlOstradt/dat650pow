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
	block := Block{PrevBlockHash: prevHash, Transactions: txs}
	block.Timestamp = time.Now().Unix()
	block.Hash = []byte{}
	block.Nonce = -1
	mBlock := MarshalBlock(block)
	sendChallenge(mBlock)
	block = awaitResponse()
	fmt.Print("   ", block.Nonce)
	fmt.Println("\n ")
	return block
}

func sendChallenge(mBlock []byte) {
	_, err := conn.WriteToUDP(mBlock, slave1)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = conn.WriteToUDP(mBlock, slave2)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func awaitResponse() Block {

	for {
		buffer := make([]byte, 1024)
		n, fromAddr, err := conn.ReadFromUDP(buffer)

		var block Block
		err = json.Unmarshal(buffer[3:n], &block)
		if err != nil {
			fmt.Println(err.Error())
		}
		if chain.ValidateBlock(&block) {
			fmt.Printf("%v\n", fromAddr)
			// fmt.Println(addrEquals(fromAddr, slave2))
			return block
		}

	}

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

func addrEquals(a, b *net.UDPAddr) bool {
	return a.IP.Equal(b.IP) && a.Port == b.Port
}
