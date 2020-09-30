package main

import (
	"bytes"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newMockUTXOSet(tx *Transaction) UTXOSet {
	txid := hex.EncodeToString(tx.ID)
	return UTXOSet{txid: tx.Vout}
}

func getTestExpectedUTXOSet(block string) []UTXOSet {
	return testUTXOs[block].expectedUTXOs
}

func getTestSpendableOutputs(utxos UTXOSet, pubKeyHash []byte) map[string][]int {
	unspentOutputs := make(map[string][]int)

	for txID, outputs := range utxos {
		for outIdx, out := range outputs {
			if bytes.Equal(out.PubKeyHash, pubKeyHash) {
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}
	return unspentOutputs
}

var testUTXOs = map[string]struct {
	utxos         UTXOSet
	expectedUTXOs []UTXOSet
}{
	"genesis": { // (0 input -> 1 output, generating "coins")
		utxos: UTXOSet{},
		expectedUTXOs: []UTXOSet{
			UTXOSet{
				"9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2": testTransactions["tx0"].Vout,
				//tx0: Address 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh create coinbase transaction and received 10 "coins"
			},
		},
	},
	"block1": { // (1 input -> 2 outputs, splitting one input)
		utxos: UTXOSet{
			"9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2": testTransactions["tx0"].Vout,
		},
		expectedUTXOs: []UTXOSet{
			UTXOSet{
				"397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13": testTransactions["tx1"].Vout,
				//tx1: 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh sent 5 "coins" to 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX and get 5 as remainder
			},
		},
	},
	"block2": { // (1 input -> 2 output, with multiple txs)
		utxos: UTXOSet{
			"397b990007845099b4fe50ba23490f277b3bf6f5316b4082c343b14c5504ab13": testTransactions["tx1"].Vout,
		},
		expectedUTXOs: []UTXOSet{
			UTXOSet{
				"e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b": testTransactions["tx2"].Vout,
				//tx2: 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh sent 1 "coin" to 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX and get 4 as remainder
				"dcd76d254f7a41888e6bda9958c4ceadf510e1bd5fd251f617c91b704fbf9492": testTransactions["tx3"].Vout,
				//tx3: 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX sent 3 "coins" to 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh and get 2 as remainder
			},
		},
	},
	"block3": { // (1 input -> 2 outputs)
		utxos: UTXOSet{
			"e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b": testTransactions["tx2"].Vout[:1], // ignore output to reduce combination of possible input sources
			"dcd76d254f7a41888e6bda9958c4ceadf510e1bd5fd251f617c91b704fbf9492": testTransactions["tx3"].Vout,
		},
		expectedUTXOs: []UTXOSet{
			UTXOSet{
				"e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b": testTransactions["tx2"].Vout[:1],
				"dcd76d254f7a41888e6bda9958c4ceadf510e1bd5fd251f617c91b704fbf9492": testTransactions["tx3"].Vout[1:], // spent first output
				// 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh spent his previous output that was here (i.e., TXOutput{3, rodrigoPubKeyHash}) using it as the input of the new tx below:
				"91d6fe8fe351e50fa6e16bb391ff74f5dc650646ce6ad02442e647742566b31b": testTransactions["tx4"].Vout,
				//tx4: 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh sent 2 "coins" to 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX and get 1 as remainder
			},
		},
	},
	"block4": { // (2 inputs -> 1 output)
		utxos: UTXOSet{
			"e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b": testTransactions["tx2"].Vout,
			"91d6fe8fe351e50fa6e16bb391ff74f5dc650646ce6ad02442e647742566b31b": testTransactions["tx4"].Vout,
		},
		expectedUTXOs: []UTXOSet{ // A map is not ordered different positions of the inputs generates a different hash of the transaction, this is why here we are duplicating the utxo set only changing the tx's hash and testing if either exists.
			UTXOSet{
				"e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b": testTransactions["tx2"].Vout[1:], // spent first output
				"91d6fe8fe351e50fa6e16bb391ff74f5dc650646ce6ad02442e647742566b31b": testTransactions["tx4"].Vout[1:], // spent first output
				"b63d956b234d27c3494d9935ac9764634db0232f32ef7f576979d8ba5ec93fbc": testTransactions["tx5"].Vout,
				//tx5: 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX sent 3 "coins" to 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh
			},
			UTXOSet{ //FIXME: isn't the best approach, but will work for now.
				"91d6fe8fe351e50fa6e16bb391ff74f5dc650646ce6ad02442e647742566b31b": testTransactions["tx4"].Vout[1:], // spent first output
				"e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b": testTransactions["tx2"].Vout[1:], // spent first output
				"7e062c43cdc7f7b04f9c4962b88c38fd562b667100a2e108421a9ccc853eeab8": testTransactions["tx5"].Vout,
				//tx5: 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX sent 3 "coins" to 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh
			},
		},
	},
}

func TestFindSpendableOutputsFromOneOutput(t *testing.T) {
	utxos := getTestExpectedUTXOSet("genesis")[0]
	expectedOut := utxos["9402c56f49de02d2b9c4633837d82e3881227a3ea90c4073c02815fdcf5afaa2"][0]
	expectedValue := expectedOut.Value

	pubKeyHash := expectedOut.PubKeyHash
	expectedUnspentOutputs := getTestSpendableOutputs(utxos, pubKeyHash)

	// Find the 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh unspent TXOutputs
	accumulatedAmount, unspentOutputs := utxos.FindSpendableOutputs(pubKeyHash, 5)

	assert.Equal(t, expectedValue, accumulatedAmount)
	assert.Equal(t, expectedUnspentOutputs, unspentOutputs)
}

func TestFindSpendableOutputsFromMultipleOutputs(t *testing.T) {
	utxos := getTestExpectedUTXOSet("block2")[0]
	out1 := utxos["e9e5fc159f24b2b33310f77aef4e425e77ed71be87dbf9a0c7764b5417bd3e4b"][1]
	out2 := utxos["dcd76d254f7a41888e6bda9958c4ceadf510e1bd5fd251f617c91b704fbf9492"][0]
	expectedValue := out1.Value + out2.Value

	expectedUnspentOutputs := getTestSpendableOutputs(utxos, out1.PubKeyHash)

	// Find the 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh unspent TXOutputs
	accumulatedAmount, unspentOutputs := utxos.FindSpendableOutputs(out1.PubKeyHash, 5)

	assert.Equal(t, expectedValue, accumulatedAmount)
	assert.Equal(t, expectedUnspentOutputs, unspentOutputs)
}

func TestFindUTXO(t *testing.T) {
	// 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh create a coinbase transaction, receiving 10 "coins"
	utxos := getTestExpectedUTXOSet("genesis")[0]

	rodrigoPubKeyHash := pubKeyHashFromAddress("14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh")
	leanderPubKeyHash := pubKeyHashFromAddress("1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX")

	utxoRodrigo := utxos.FindUTXO(rodrigoPubKeyHash)
	assert.Equal(t, []TXOutput{TXOutput{BlockReward, rodrigoPubKeyHash}}, utxoRodrigo)

	utxoLeander := utxos.FindUTXO(leanderPubKeyHash)
	assert.Equal(t, []TXOutput(nil), utxoLeander)

	// 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh sent 5 "coins" to 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX
	// update utxo
	utxos = getTestExpectedUTXOSet("block1")[0]
	utxoRodrigo = utxos.FindUTXO(rodrigoPubKeyHash)
	assert.Equal(t, []TXOutput{TXOutput{5, rodrigoPubKeyHash}}, utxoRodrigo)

	utxoLeander = utxos.FindUTXO(leanderPubKeyHash)
	assert.Equal(t, []TXOutput{TXOutput{5, leanderPubKeyHash}}, utxoLeander)

	// 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh sent 1 "coin" to 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX and
	// 1HrwWkjdwQuhaHSco9H7u7SVsmo4aeDZBX sent 3 "coins" to 14vRYoWsjqC61tNmaLPPzjKnxirSxFoehh
	// update utxo
	utxos = getTestExpectedUTXOSet("block2")[0]

	utxoRodrigo = utxos.FindUTXO(rodrigoPubKeyHash)
	assert.ElementsMatch(t, []TXOutput{
		TXOutput{4, rodrigoPubKeyHash},
		TXOutput{3, rodrigoPubKeyHash},
	}, utxoRodrigo)
	assert.Equal(t, 2, len(utxoRodrigo))

	utxoLeander = utxos.FindUTXO(leanderPubKeyHash)
	assert.ElementsMatch(t, []TXOutput{
		TXOutput{2, leanderPubKeyHash},
		TXOutput{1, leanderPubKeyHash},
	}, utxoLeander)
	assert.Equal(t, 2, len(utxoLeander))
}

func TestCountUTXOs(t *testing.T) {
	utxos := getTestExpectedUTXOSet("genesis")[0]
	assert.Equal(t, 1, utxos.CountUTXOs())

	utxos = getTestExpectedUTXOSet("block1")[0]
	assert.Equal(t, 2, utxos.CountUTXOs())

	utxos = getTestExpectedUTXOSet("block2")[0]
	assert.Equal(t, 4, utxos.CountUTXOs())
}

func TestUpdate(t *testing.T) {
	for k, m := range testUTXOs {
		utxos := m.utxos
		switch k {
		case "genesis":
			utxos.Update([]*Transaction{testTransactions["tx0"]})
			assert.True(t, utxos.Equal(m.expectedUTXOs[0]))
		case "block1":
			utxos.Update([]*Transaction{testTransactions["tx1"]})
			assert.True(t, utxos.Equal(m.expectedUTXOs[0]))
		case "block2":
			utxos.Update([]*Transaction{testTransactions["tx2"], testTransactions["tx3"]})
			assert.True(t, utxos.Equal(m.expectedUTXOs[0]))
		case "block3":
			utxos.Update([]*Transaction{testTransactions["tx4"]})
			assert.True(t, utxos.Equal(m.expectedUTXOs[0]))
		case "block4":
			utxos.Update([]*Transaction{testTransactions["tx5"]})
			ok1 := utxos.Equal(m.expectedUTXOs[0])
			ok2 := utxos.Equal(m.expectedUTXOs[1])
			assert.True(t, ok1 || ok2)
		}
	}
}
