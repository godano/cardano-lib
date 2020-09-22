package chain

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/fxamacker/cbor/v2"
	cardano "github.com/godano/cardano-lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

// This method shall test, if the transaction below can be encoded correctly. It
// should match the output of the cardano-cli tool.
//
// cardano-cli shelley transaction build-raw
//        --tx-in 4e3a6e7fdcb0d0efa17bf79c13aed2b4cb9baf37fb1aa2e39553d5bd720c5c99#4 \
//        --tx-out addr1v9j7tj3qh877jcuk2qst4fu8vpaxjm9pzm8pyuhkt0evmuqdmzrua+100000000 \
//        --tx-out addr1vy96pgqgxw0hsef3p79usvxzqwqagsg48qen3wnq0u9tlugdvptzr+999899832035 \
//        --ttl 796546 \
//        --fee 167965 \
//        --out-file tx_01.raw
func TestTransaction_AssembleValidRawPaymentTransaction01_EncodeCorrectly(t *testing.T) {
	expectedCBORHex := "82a400818258204e3a6e7fdcb0d0efa17bf79c13aed2b4cb9baf37fb1aa2e39553d5bd720c5c9904018282581d6165e5ca20b9fde963965020baa787607a696ca116ce1272f65bf2cdf01a05f5e10082581d610ba0a008339f7865310f8bc830c20381d44115383338ba607f0abff11b000000e8ceac9ee3021a0002901d031a000c2782f6"

	transaction := Transaction{
		Body: TransactionBody{
			Inputs: []TransactionInput{
				{
					ID:    []byte{},
					Index: 4,
				},
			},
			Outputs: []TransactionOutput{
				{
					Address: []byte{},
					Amount:  100000000,
				},
				{
					Address: []byte{},
					Amount:  999899832035,
				},
			},
			Fee: 167965,
			TTL: 796546,
		},
	}

	var b bytes.Buffer
	dataWriter := bufio.NewWriter(&b)

	encoder := cbor.NewEncoder(dataWriter)
	err := encoder.Encode(&transaction)
	_ = dataWriter.Flush()
	assert.Nil(t, err)

	cborHex := hex.EncodeToString(b.Bytes())
	assert.Equal(t, expectedCBORHex, cborHex)
}

// This method shall test, if the transaction below can be decoded correctly. It
// should contain all the details of the transaction.
//
// cardano-cli shelley transaction build-raw
//        --tx-in 4e3a6e7fdcb0d0efa17bf79c13aed2b4cb9baf37fb1aa2e39553d5bd720c5c99#4 \
//        --tx-out addr1v9j7tj3qh877jcuk2qst4fu8vpaxjm9pzm8pyuhkt0evmuqdmzrua+100000000 \
//        --tx-out addr1vy96pgqgxw0hsef3p79usvxzqwqagsg48qen3wnq0u9tlugdvptzr+999899832035 \
//        --ttl 796546 \
//        --fee 167965 \
//        --out-file tx_01.raw
func TestTransaction_AssembleValidRawPaymentTransaction01_DecodeCorrectly(t *testing.T) {
	expectedCBORHex := "82a400818258204e3a6e7fdcb0d0efa17bf79c13aed2b4cb9baf37fb1aa2e39553d5bd720c5c9904018282581d6165e5ca20b9fde963965020baa787607a696ca116ce1272f65bf2cdf01a05f5e10082581d610ba0a008339f7865310f8bc830c20381d44115383338ba607f0abff11b000000e8ceac9ee3021a0002901d031a000c2782f6"
	data, _ := hex.DecodeString(expectedCBORHex)

	var transaction Transaction
	dec := cbor.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&transaction)

	if assert.Nil(t, err) {
		assert.NotNil(t, transaction.Body)

		// Fees and TTL
		assert.Equal(t, uint64(167965), transaction.Body.Fee, "the fee must be passed correctly to 167965")
		assert.Equal(t, uint64(796546), transaction.Body.TTL, "the TTL must be passed correctly to 796546")

		// transaction inputs
		assert.NotNil(t, transaction.Body.Inputs)
		assert.NotEmpty(t, transaction.Body.Inputs)
		assert.Equal(t, 1, len(transaction.Body.Inputs), "the transaction must have one input")
		assert.Equal(t, uint64(4), transaction.Body.Inputs[0].Index)

		// transaction outputs
		assert.NotNil(t, transaction.Body.Outputs)
		assert.NotEmpty(t, transaction.Body.Outputs)
		if assert.Equal(t, 2, len(transaction.Body.Outputs), "the transaction must have two outputs") {
			assert.Equal(t, cardano.Coin(100000000), transaction.Body.Outputs[0].Amount)
			addrOut1Data := transaction.Body.Outputs[0].Address
			addressOut1, err := bech32.Encode("addr", addrOut1Data[1:])
			if assert.Nil(t, err) {
				assert.Equal(t, "addr1v9j7tj3qh877jcuk2qst4fu8vpaxjm9pzm8pyuhkt0evmuqdmzrua", addressOut1)
			}

			assert.Equal(t, cardano.Coin(999899832035), transaction.Body.Outputs[1].Amount)
			addrOut2Data := transaction.Body.Outputs[1].Address
			addressOut2, err := bech32.Encode("addr", addrOut2Data[1:])
			if assert.Nil(t, err) {
				assert.Equal(t, "addr1vy96pgqgxw0hsef3p79usvxzqwqagsg48qen3wnq0u9tlugdvptzr", addressOut2)
			}
		}
	}
}
