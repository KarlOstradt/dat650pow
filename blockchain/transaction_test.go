package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newMockCoinbaseTX(address, data, txid string) *Transaction {
	return &Transaction{
		ID: Hex2Bytes(txid),
		Vin: []TXInput{
			TXInput{Txid: []byte{}, OutIdx: -1, Signature: nil, PubKey: []byte(data)},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      BlockReward,
				PubKeyHash: pubKeyHashFromAddress(address),
			},
		},
	}
}

func removeTXInputSignature(tx *Transaction) {
	var inputs []TXInput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.OutIdx, nil, vin.PubKey})
	}
	*tx = Transaction{tx.ID, inputs, tx.Vout}
}

var testTransactions = map[string]*Transaction{
	"tx0": &Transaction{
		ID: Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"),
		Vin: []TXInput{
			TXInput{Txid: []byte{}, OutIdx: -1, Signature: nil, PubKey: []byte(GenesisCoinbaseData)},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      BlockReward,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	},
	"tx1": &Transaction{
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
	},
	"tx2": &Transaction{
		ID: Hex2Bytes("e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
				OutIdx:    1,
				Signature: nil,
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      1,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      4,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	},
	"tx3": &Transaction{
		ID: Hex2Bytes("dcd76d254f7a41888e6bda9958c4ceadf510e1bd5fd251f617c91b704fbf9492"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
				OutIdx:    0,
				Signature: nil,
				PubKey:    Hex2Bytes("c36d68bc641029e53a38252b436c596ef3d03a4a754743da50fb9a321020e882dd401732381783c7444112abc729b3bee04643015d80fe67e0c28a5b28a20910"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      3,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
			TXOutput{
				Value:      2,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
		},
	},
	"tx4": &Transaction{
		ID: Hex2Bytes("91d6fe8fe351e50fa6e16bb391ff74f5dc650646ce6ad02442e647742566b31b"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("dcd76d254f7a41888e6bda9958c4ceadf510e1bd5fd251f617c91b704fbf9492"),
				OutIdx:    0,
				Signature: nil,
				PubKey:    Hex2Bytes("f86aa0caf08359ee4227d2901ab490172c69a801910f4140cdde2f5dc8f8bb3dc19da2c9fb0ed041db106a8fea0382de25edbc83df6893574e40fc2e1e493748"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      2,
				PubKeyHash: pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"),
			},
			TXOutput{
				Value:      1,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	},
	"tx5": &Transaction{
		ID: Hex2Bytes("b63d956b234d27c3494d9935ac9764634db0232f32ef7f576979d8ba5ec93fbc"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b"),
				OutIdx:    0,
				Signature: nil,
				PubKey:    Hex2Bytes("c36d68bc641029e53a38252b436c596ef3d03a4a754743da50fb9a321020e882dd401732381783c7444112abc729b3bee04643015d80fe67e0c28a5b28a20910"),
			},
			TXInput{
				Txid:      Hex2Bytes("91d6fe8fe351e50fa6e16bb391ff74f5dc650646ce6ad02442e647742566b31b"),
				OutIdx:    0,
				Signature: nil,
				PubKey:    Hex2Bytes("c36d68bc641029e53a38252b436c596ef3d03a4a754743da50fb9a321020e882dd401732381783c7444112abc729b3bee04643015d80fe67e0c28a5b28a20910"),
			},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      3,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	},
}

func TestHash(t *testing.T) {

	tx := &Transaction{
		ID: []byte{},
		Vin: []TXInput{
			TXInput{Txid: []byte{}, OutIdx: -1, Signature: nil, PubKey: []byte(GenesisCoinbaseData)},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      BlockReward,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}
	assert.Equal(t, Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"), tx.Hash())

	tx = &Transaction{
		ID: []byte{},
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
	assert.Equal(t, Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"), tx.Hash())
}

func TestIsCoinbase(t *testing.T) {

	tx := testTransactions["tx0"]
	assert.True(t, tx.IsCoinbase())

	tx = testTransactions["tx1"]
	assert.False(t, tx.IsCoinbase())
}

func TestNewCoinbaseTX(t *testing.T) {

	// Passing data to the coinbase transaction
	tx := NewCoinbaseTX("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh", GenesisCoinbaseData)
	assert.Equal(t, Hex2Bytes("5468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73"), tx.Vin[0].PubKey)
	assert.Equal(t, -1, tx.Vin[0].OutIdx)
	assert.Equal(t, []byte{}, tx.Vin[0].Txid)
	assert.Nil(t, tx.Vin[0].Signature)
	assert.Equal(t, BlockReward, tx.Vout[0].Value)
	assert.Equal(t, pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"), tx.Vout[0].PubKeyHash)

	// Using default data
	tx = NewCoinbaseTX("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh", "")
	assert.Equal(t, []byte{}, tx.Vin[0].Txid)
	assert.Equal(t, -1, tx.Vin[0].OutIdx)
	assert.Equal(t, BlockReward, tx.Vout[0].Value)
	assert.Equal(t, pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"), tx.Vout[0].PubKeyHash)

	tx = NewCoinbaseTX("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX", "")
	assert.Equal(t, []byte{}, tx.Vin[0].Txid)
	assert.Equal(t, -1, tx.Vin[0].OutIdx)
	assert.Equal(t, BlockReward, tx.Vout[0].Value)
	assert.Equal(t, pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"), tx.Vout[0].PubKeyHash)
}

func TestNewUTXOTransaction(t *testing.T) {
	priv1, pub1 := decodeKeyPair(testEncPrivKeyUser1, testEncPubKeyUser1)
	fromWallet := CreateWallet(priv1, pub1)
	fromAddress := "14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"

	priv2, pub2 := decodeKeyPair(testEncPrivKeyUser2, testEncPubKeyUser2)
	toWallet := CreateWallet(priv2, pub2)
	toAddress := "1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX"

	// "from" address have 10 (i.e., genesis coinbase) and "to" address have 0
	bc := newMockBlockchain(fromAddress)
	utxos := UTXOSet{
		"9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2": testTransactions["tx0"].Vout,
	}

	// Reject if there is not sufficient funds
	tx1, err := NewUTXOTransaction(toWallet, fromAddress, 5, utxos, bc)
	assert.Errorf(t, err, "Not enough funds")
	assert.Nil(t, tx1)

	// Accept otherwise
	tx1, err = NewUTXOTransaction(fromWallet, toAddress, 5, utxos, bc)
	// NOTE: We are ignoring the signatures checks for now
	// so blocks without signatures or with invalid signatures can be accepted
	removeTXInputSignature(tx1)
	assert.Equal(t, testTransactions["tx1"], tx1)
	assert.Nil(t, err)

	// update utxo and blockchain with tx1
	bc.AddBlock([]*Transaction{testTransactions["tx1"]})
	utxos = UTXOSet{
		"397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13": testTransactions["tx1"].Vout,
	}

	tx2, err := NewUTXOTransaction(fromWallet, toAddress, 1, utxos, bc)
	removeTXInputSignature(tx2)
	assert.Equal(t, testTransactions["tx2"], tx2)
	assert.Nil(t, err)

	tx3, err := NewUTXOTransaction(toWallet, fromAddress, 3, utxos, bc)
	removeTXInputSignature(tx3)
	assert.Equal(t, testTransactions["tx3"], tx3)
	assert.Nil(t, err)
}

func TestSign(t *testing.T) {
	privKey, _ := decodeKeyPair(testEncPrivKeyUser1, testEncPubKeyUser1)

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

	prevTXs := make(map[string]*Transaction)
	prevTXs["9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"] = testTransactions["tx0"]

	tx.Sign(*privKey, prevTXs)
	assert.NotNil(t, tx.Vin[0].Signature)
}

func TestSignIgnoreCoinbaseTX(t *testing.T) {
	privKey, _ := decodeKeyPair(testEncPrivKeyUser1, testEncPubKeyUser1)

	tx := &Transaction{
		ID: Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"),
		Vin: []TXInput{
			TXInput{Txid: []byte{}, OutIdx: -1, Signature: nil, PubKey: []byte(GenesisCoinbaseData)},
		},
		Vout: []TXOutput{
			TXOutput{
				Value:      BlockReward,
				PubKeyHash: pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh"),
			},
		},
	}
	prevTXs := make(map[string]*Transaction)

	tx.Sign(*privKey, prevTXs)
	assert.Nil(t, tx.Vin[0].Signature)
}

func TestSignInvalidInputTX(t *testing.T) {
	privKey, _ := decodeKeyPair(testEncPrivKeyUser1, testEncPubKeyUser1)

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

	prevTXs := make(map[string]*Transaction)
	prevTXs["9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"] = testTransactions["tx0"]

	assert.PanicsWithValue(t, "Current input transaction isn't listed in previous transactions", func() { tx.Sign(*privKey, prevTXs) })
	assert.Nil(t, tx.Vin[0].Signature)
}

func TestVerifyIgnoreCoinbaseTX(t *testing.T) {
	tx := testTransactions["tx0"]
	prevTXs := make(map[string]*Transaction)
	assert.True(t, tx.Verify(prevTXs))
}

func TestVerify(t *testing.T) {
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

	prevTXs := make(map[string]*Transaction)
	prevTXs["9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"] = testTransactions["tx0"]

	assert.True(t, tx.Verify(prevTXs))
}

func TestVerifyInvalidInputTX(t *testing.T) {
	tx := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("non-existentID"),
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

	prevTXs := make(map[string]*Transaction)
	prevTXs["9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"] = testTransactions["tx0"]

	assert.PanicsWithValue(t, "Current input transaction isn't listed in previous transactions", func() { tx.Verify(prevTXs) })
}

func TestVerifyInvalidSignature(t *testing.T) {
	tx := &Transaction{
		ID: Hex2Bytes("397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13"),
		Vin: []TXInput{
			TXInput{
				Txid:      Hex2Bytes("9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"),
				OutIdx:    0,
				Signature: Hex2Bytes("invalid"),
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

	prevTXs := make(map[string]*Transaction)
	prevTXs["9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"] = testTransactions["tx0"]

	assert.False(t, tx.Verify(prevTXs))
}
func TestTrimmedCopy(t *testing.T) {
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

	txCopy := tx.TrimmedCopy()

	assert.Nil(t, txCopy.Vin[0].Signature)
	assert.Nil(t, txCopy.Vin[0].PubKey)
	assert.Equal(t, tx.Vin[0].Txid, txCopy.Vin[0].Txid)
	assert.Equal(t, tx.Vin[0].OutIdx, txCopy.Vin[0].OutIdx)
	assert.Equal(t, tx.Vout, txCopy.Vout)
	assert.Equal(t, tx.ID, txCopy.ID)
}
