package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newMockBlockchain(address string) *Blockchain {
	genesis := testBlockchainData["block0"]
	return &Blockchain{[]*Block{genesis}}
}

func addMockBlock(bc *Blockchain, newBlock *Block) {
	bc.blocks = append(bc.blocks, newBlock)
}

var testBlockchainData = map[string]*Block{
	"block0": &Block{
		Timestamp:     BlockTime,
		Transactions:  []*Transaction{testTransactions["tx0"]},
		PrevBlockHash: []byte{},
		Hash:          Hex2Bytes("00d4eeaee903dce5468d4c6975376dfbc4c45ea1bc6c5bbbfd8e13b26aaf6e3b"),
		Nonce:         59,
	},
	"block1": &Block{
		Timestamp:     BlockTime,
		Transactions:  []*Transaction{testTransactions["tx1"]},
		PrevBlockHash: Hex2Bytes("00d4eeaee903dce5468d4c6975376dfbc4c45ea1bc6c5bbbfd8e13b26aaf6e3b"),
		Hash:          Hex2Bytes("00361766580e4eb8d0b199b2b24b5caa48cd1b4682938fbffc0f584d0f2ae5aa"),
		Nonce:         1953,
	},
	"block2": &Block{
		Timestamp: BlockTime,
		Transactions: []*Transaction{
			testTransactions["tx2"],
			testTransactions["tx3"],
		},
		PrevBlockHash: Hex2Bytes("00361766580e4eb8d0b199b2b24b5caa48cd1b4682938fbffc0f584d0f2ae5aa"),
		Hash:          Hex2Bytes("0000dd6985d09b2accd326e386ac4fa6d46bd542a6316bcf072a89c5ef63a871"),
		Nonce:         809,
	},
	"block3": &Block{
		Timestamp:     BlockTime,
		Transactions:  []*Transaction{testTransactions["tx4"]},
		PrevBlockHash: Hex2Bytes("0000dd6985d09b2accd326e386ac4fa6d46bd542a6316bcf072a89c5ef63a871"),
		Hash:          Hex2Bytes("00ae7ee91b5b6547e7ae2ce0afb194dcf96552ce48960e63ec7c4c9d7ca5ee35"),
		Nonce:         0,
	},
	"block4": &Block{
		Timestamp:     BlockTime,
		Transactions:  []*Transaction{testTransactions["tx5"]},
		PrevBlockHash: Hex2Bytes("00ae7ee91b5b6547e7ae2ce0afb194dcf96552ce48960e63ec7c4c9d7ca5ee35"),
		Hash:          Hex2Bytes("00b2683ffe7f8502103ee822390c724e226fc22b992029a69aac1c745470e711"),
		Nonce:         114,
	},
}

func TestBlockchain(t *testing.T) {
	// NewBlockchain
	bc := NewBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	assert.NotNil(t, bc)
	assert.Equal(t, 1, len(bc.blocks))

	// GetGenesisBlock
	gb := bc.GetGenesisBlock()
	assert.NotNil(t, gb)
	assert.Equal(t, []byte{}, gb.PrevBlockHash, "Genesis block shouldn't has PrevBlockHash")

	// Genesis block should contains a genesis transaction
	coinbaseTx := gb.Transactions[0]
	assert.Equal(t, 1, len(gb.Transactions))
	assert.Equal(t, testTransactions["tx0"], coinbaseTx)
}

func TestAddBlock(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	assert.Equal(t, 1, len(bc.blocks))

	// AddBlock
	b1 := bc.AddBlock([]*Transaction{testTransactions["tx1"]})
	assert.NotNil(t, b1)
	assert.Equal(t, 2, len(bc.blocks))

	gb := bc.blocks[0]
	assert.Equalf(t, gb.Hash, b1.PrevBlockHash, "Genesis block Hash: %x isn't equal to current PrevBlockHash: %x", gb.Hash, b1.PrevBlockHash)

	b2 := bc.AddBlock([]*Transaction{testTransactions["tx3"]})
	assert.NotNil(t, b2)
	assert.Equal(t, 3, len(bc.blocks))
	assert.Equalf(t, b1.Hash, b2.PrevBlockHash, "Previous block Hash: %x isn't equal to the expected: %x", b2.PrevBlockHash, b1.Hash)
}

func TestCurrentBlock(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")

	// CurrentBlock
	b := bc.CurrentBlock()
	assert.NotNil(t, b)
	expectedBlock := bc.blocks[0]
	assert.Equalf(t, expectedBlock.Hash, b.Hash, "Current block Hash: %x isn't the expected: %x", b.Hash, expectedBlock.Hash)

	addMockBlock(bc, testBlockchainData["block1"])

	b = bc.CurrentBlock()
	assert.NotNil(t, b)
	expectedBlock = bc.blocks[1]
	assert.Equalf(t, expectedBlock.Hash, b.Hash, "Current block Hash: %x isn't the expected: %x", b.Hash, expectedBlock.Hash)
}

func TestGetBlock(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	// GetBlock
	b, err := bc.GetBlock(bc.blocks[0].Hash)
	assert.Nil(t, err)
	assert.NotNil(t, b)

	assert.Equalf(t, bc.blocks[0].Hash, b.Hash, "Block Hash: %x isn't the expected: %x", b.Hash, bc.blocks[0].Hash)
}

func TestMineBlockWithInvalidTxInput(t *testing.T) {
	fromAddress := "14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"
	toAddress := "1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"
	bc := newMockBlockchain(fromAddress)

	fromPublicKey := decodePublicKey(testEncPubKeyUser1)
	// Ignore transaction that refer to non-existent transaction input
	invalidTx := &Transaction{
		ID: Hex2Bytes("bce268225bc12a0015bcc39e91d59f47fd176e64ca42e4f8aecf107fe38f3bfa"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("non-existentID"),
				OutIdx:    0,
				Signature: nil,
				PubKey:    pubKeyToByte(*fromPublicKey),
			},
		},
		Vout: []TXOutput{
			TXOutput{Value: 5, PubKeyHash: pubKeyHashFromAddress(toAddress)},
			TXOutput{Value: 5, PubKeyHash: pubKeyHashFromAddress(fromAddress)},
		},
	}

	b, err := bc.MineBlock([]*Transaction{invalidTx})
	assert.Error(t, err, "there are no valid transactions to be mined")
	assert.Nil(t, b)
}

func TestMineBlock(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	tx := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"),
				OutIdx:    0,
				Signature: Hex2Bytes("17b6db89942bb02b485332c9a3b37638e02a3dfafdf4c3a4fad7fc4c7b062cc8156b75957050e049cd307853522f5ef49339b1b1230359f59571af12c612bde2"),
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}

	b, err := bc.MineBlock([]*Transaction{tx})
	assert.Nil(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, 2, len(bc.blocks))
	gb := bc.blocks[0]
	assert.Equalf(t, gb.Hash, b.PrevBlockHash, "Genesis block Hash: %x isn't equal to current PrevBlockHash: %x", gb.Hash, b.PrevBlockHash)

	minedBlock, err := bc.GetBlock(b.Hash)
	assert.Equal(t, b, minedBlock)

	txMinedBlock := bc.blocks[1].Transactions[0]
	assert.NotNil(t, txMinedBlock)
	assert.Equal(t, tx.ID, txMinedBlock.ID)
}

func TestValidateBlockWithoutCoinbase(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	assert.False(t, bc.ValidateBlock(testBlockchainData["block1"]))
}

func TestValidateBlockCoinbaseOrder(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	assert.True(t, bc.ValidateBlock(testBlockchainData["block0"]))

	coinbaseTX := newMockCoinbaseTX("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh", "test", "f51bb83bf43031b5c1bf1b01bb0a57b837c52f560412924cf3786e3fe5e26337")

	validTX := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"),
				OutIdx:    0,
				Signature: Hex2Bytes("17b6db89942bb02b485332c9a3b37638e02a3dfafdf4c3a4fad7fc4c7b062cc8156b75957050e049cd307853522f5ef49339b1b1230359f59571af12c612bde2"),
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}
	block := &Block{
		Timestamp:     BlockTime,
		Transactions:  []*Transaction{validTX, coinbaseTX}, // wrong coinbase order, should be the first transaction in the list!
		PrevBlockHash: testBlockchainData["block0"].Hash,
		Hash:          Hex2Bytes("00506151f37ae36c1d82006f8b7f49c96be9b8e835cbc326c73d22ede3c0914b"),
		Nonce:         89,
	}
	assert.False(t, bc.ValidateBlock(block))
}

func TestValidateBlockWithInvalidHash(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	block := &Block{
		Timestamp:     BlockTime,
		Transactions:  []*Transaction{testTransactions["tx0"]},
		PrevBlockHash: []byte{},
		Hash:          Hex2Bytes("73d40a0510b6327d0fbcd4a2baf6e7a70f2de174ad2c84538a7b09320e9db3f2"),
		Nonce:         23,
	}
	assert.False(t, bc.ValidateBlock(block))
}

func TestValidateBlockWithInvalidNonce(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	block := &Block{
		Timestamp:     BlockTime,
		Transactions:  []*Transaction{testTransactions["tx0"]},
		PrevBlockHash: []byte{},
		Hash:          Hex2Bytes("00d4eeaee903dce5468d4c6975376dfbc4c45ea1bc6c5bbbfd8e13b26aaf6e3b"),
		Nonce:         1,
	}
	assert.False(t, bc.ValidateBlock(block))
}

func TestValidateBlockSuccess(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	assert.True(t, bc.ValidateBlock(testBlockchainData["block0"]))

	coinbaseTX := newMockCoinbaseTX("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh", "test", "f51bb83bf43031b5c1bf1b01bb0a57b837c52f560412924cf3786e3fe5e26337")

	validTX := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"),
				OutIdx:    0,
				Signature: Hex2Bytes("17b6db89942bb02b485332c9a3b37638e02a3dfafdf4c3a4fad7fc4c7b062cc8156b75957050e049cd307853522f5ef49339b1b1230359f59571af12c612bde2"),
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}
	block := &Block{
		Timestamp:     BlockTime,
		Transactions:  []*Transaction{coinbaseTX, validTX},
		PrevBlockHash: testBlockchainData["block0"].Hash,
		Hash:          Hex2Bytes("00fd22275934205b0f5cfa4adbc4648a5ae61fb148493e59c43cfbbe41187306"),
		Nonce:         249,
	}
	assert.True(t, bc.ValidateBlock(block))
}

func TestSignTransactionSuccess(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	tx := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"),
				OutIdx:    0,
				Signature: nil,
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}
	privKey, _ := decodeKeyPair(testEncPrivKeyUser1, testEncPubKeyUser1)
	bc.SignTransaction(tx, *privKey)

	assert.NotNil(t, tx.Vin[0].Signature)
}

func TestSignTransactionWithInvalidTxInput(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	tx := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("non-existentID"),
				OutIdx:    0,
				Signature: nil,
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}
	privKey, _ := decodeKeyPair(testEncPrivKeyUser1, testEncPubKeyUser1)

	assert.PanicsWithValue(t, "Transaction not found in any block", func() { bc.SignTransaction(tx, *privKey) })
	assert.Nil(t, tx.Vin[0].Signature)
}

func TestVerifyTransactionSuccess(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")

	coinbaseTx := bc.GetGenesisBlock().Transactions[0]
	assert.True(t, bc.VerifyTransaction(coinbaseTx))

	tx := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"),
				OutIdx:    0,
				Signature: Hex2Bytes("17b6db89942bb02b485332c9a3b37638e02a3dfafdf4c3a4fad7fc4c7b062cc8156b75957050e049cd307853522f5ef49339b1b1230359f59571af12c612bde2"),
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}
	assert.True(t, bc.VerifyTransaction(tx))
}
func TestVerifyTransactionInvalidTxInput(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	tx := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("non-existentID"),
				OutIdx:    0,
				Signature: nil,
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      5,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}
	assert.False(t, bc.VerifyTransaction(tx))
}

func TestFindTransactionSuccess(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")

	// Find coinbase genesis transaction
	tx, err := bc.FindTransaction(Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"))
	assert.Nil(t, err)
	assert.NotNil(t, tx)
}

func TestFindTransactionFailure(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")

	notFoundTx, err := bc.FindTransaction(Hex2Bytes("non-existentID"))
	assert.Error(t, err, "Transaction not found")
	assert.Nil(t, notFoundTx)
}

func TestFindUTXOSet(t *testing.T) {
	bc := newMockBlockchain("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")

	utxos := bc.FindUTXOSet()
	expectedUTXOs := getTestExpectedUTXOSet("genesis")[0]
	assert.Equal(t, expectedUTXOs, utxos)

	addMockBlock(bc, testBlockchainData["block1"])
	expectedUTXOs = getTestExpectedUTXOSet("block1")[0]

	utxos = bc.FindUTXOSet()
	assert.Equal(t, expectedUTXOs, utxos)
}
