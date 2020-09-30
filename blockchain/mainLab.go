package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

// var (
// 	chain          Blockchain
// 	txBuffer       []*Transaction
// 	wallet1        *Wallet
// 	wallet1Address string
// 	wallet2        *Wallet
// 	wallet2Address string
// 	utxos          UTXOSet
// )

func mainLab() {
	prompt := promptui.Select{
		Label: "Select Command",
		Items: []string{"createBlockchain", "addTransaction", "mineBlock", "printChain", "getBalance", "printWallets", "exit"},
	}

	createWallets()
	createBlockchain()

	shouldTerminate := false
	for !shouldTerminate {
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "createBlockchain":
			createBlockchain()
		case "addTransaction":
			addTransaction()
		case "mineBlock":
			mineBlock()
		case "printChain":
			printChain()
		case "getBalance":
			getBalance()
		case "printWallets":
			printWallets()
		case "exit":
			shouldTerminate = true
		}

	}

}

// func createBlockchain() {
// 	chain = *CreateBlockchain(wallet1.GetStringAddress())
// 	utxos = chain.FindUTXOSet()
// }

func addTransaction() {
	fromWallet, toWallet := selectWallet("Select sender")
	if fromWallet == nil {
		return
	}
	validate := func(input string) error {
		val, err := strconv.ParseFloat(input, 64)
		if val < 0 {
			return errors.New("Invalid number")
		}
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter the amount",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	amount, err := strconv.ParseInt(result, 10, 64)

	tx, err := NewUTXOTransaction(fromWallet, toWallet.GetStringAddress(), int(amount), utxos, &chain)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Println("Added transaction to the buffer:")
	fmt.Printf("%s\n", tx.String())

	txBuffer = append(txBuffer, tx)
	// utxos.Update([]*Transaction{tx})
	// utxos.Update(txBuffer)
}

func selectWallet(label string) (*Wallet, *Wallet) {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"wallet1: " + wallet1Address, "wallet2: " + wallet2Address},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, nil
	}
	switch result {
	case "wallet1: " + wallet1Address:
		return wallet1, wallet2
	case "wallet2: " + wallet2Address:
		return wallet2, wallet1
	}
	return nil, nil
}

func mineBlock() {
	if len(txBuffer) == 0 {
		fmt.Println("No transactions to add")
		return
	}

	minerWallet, _ := selectWallet("Select miner")
	coinbaseTX := NewCoinbaseTX(minerWallet.GetStringAddress(), GenesisCoinbaseData)
	// coinbaseTX := &Transaction{
	// 	Vin: []TXInput{
	// 		{Txid: []byte{}, OutIdx: -1, Signature: nil, PubKey: minerWallet.PublicKey},
	// 	},
	// 	Vout: []TXOutput{
	// 		*NewTXOutput(BlockReward, minerWallet.GetStringAddress()),
	// 	},
	// }
	// coinbaseTX.ID = coinbaseTX.Hash()

	txs := append([]*Transaction{coinbaseTX}, txBuffer...)
	block, err := chain.MineBlock(txs)
	fmt.Println(len(block.Transactions))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	utxos.Update(txs)
	fmt.Println("Added block to the chain:")
	fmt.Printf("%s\n", block.String())
	txBuffer = []*Transaction{}
}

// func createWallets() {
// 	wallet1 = NewWallet()
// 	wallet1Address = wallet1.GetStringAddress()
// 	wallet2 = NewWallet()
// 	wallet2Address = wallet2.GetStringAddress()
// }

func getBalance() {
	wallet, _ := selectWallet("Get balance of:")

	balance, _ := utxos.FindSpendableOutputs(HashPubKey(wallet.PublicKey), 1)
	fmt.Printf("Balance for %s: %d\n", wallet.GetStringAddress(), balance)
}

func printChain() {
	fmt.Println(chain.String())
}

func printWallets() {
	fmt.Println("\n--------Wallet1--------\t")
	fmt.Println("Address:\t" + wallet1.GetStringAddress())
	balance, _ := utxos.FindSpendableOutputs(HashPubKey(wallet1.PublicKey), 1)
	fmt.Printf("Balance:\t%d\n", balance)
	fmt.Printf("Public key:\t%x\n", wallet1.PublicKey)
	fmt.Printf("Pub key hash:\t%x\n", GetPubKeyHashFromAddress(wallet1.GetStringAddress()))

	fmt.Println("\n--------Wallet2--------\t")
	fmt.Println("Address:\t" + wallet2.GetStringAddress())
	balance, _ = utxos.FindSpendableOutputs(HashPubKey(wallet2.PublicKey), 1)
	fmt.Printf("Balance:\t%d\n", balance)
	fmt.Printf("Public key:\t%x\n", wallet2.PublicKey)
	fmt.Printf("Pub key hash:\t%x\n\n", GetPubKeyHashFromAddress(wallet2.GetStringAddress()))
}
